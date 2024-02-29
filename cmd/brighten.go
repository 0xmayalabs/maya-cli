package cmd

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/cmp"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"time"
)

const (
	MinPixelValue = 0
	MaxPixelValue = 255
)

var brighteningFactor int

// brightenConfig specifies the configuration for brightening an image by a brightening factor.
type brightenConfig struct {
	originalImg       string
	finalImg          string
	brighteningFactor int // TODO(xenowits): Convert it to floating-point
	proofDir          string
	markdownFile      string
	backend           string
}

// newBrightenCmd returns a new cobra.Command for brightening an image by a brightening factor.
func newBrightenCmd() *cobra.Command {
	var conf brightenConfig

	cmd := &cobra.Command{
		Use: "brighten",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveBrighten(conf)
		},
	}

	bindBrightenFlags(cmd, &conf)

	return cmd
}

// bindBrightenFlags binds the brightening configuration flags.
func bindBrightenFlags(cmd *cobra.Command, conf *brightenConfig) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().IntVar(&conf.brighteningFactor, "brightening-factor", 2, "The factor with which image is brightened.")
	cmd.Flags().StringVar(&conf.backend, "backend", "groth16", "The proving backend used for generating the proofs.")
}

// proveBrighten generates the zk proof of brightening an image by a brightening factor.
func proveBrighten(config brightenConfig) error {
	brighteningFactor = config.brighteningFactor

	fmt.Println("Brightening factor", brighteningFactor)

	// Open the original image file.
	originalImage, err := os.Open(config.originalImg)
	if err != nil {
		return err
	}
	defer originalImage.Close()

	// Open the final image file.
	finalImage, err := os.Open(config.finalImg)
	if err != nil {
		return err
	}
	defer finalImage.Close()

	// Get the pixel values for the original image.
	originalPixels, err := convertImgToPixels(originalImage)
	if err != nil {
		return err
	}

	// Get the pixel values for the final image.
	finalPixels, err := convertImgToPixels(finalImage)
	if err != nil {
		return err
	}

	proof, vk, circuitCompilationDuration, provingDuration, err := generateBrightenProof(config.backend, originalPixels, finalPixels)
	if err != nil {
		return err
	}

	brightenDir := path.Join(config.proofDir, "brighten")
	if err = os.MkdirAll(brightenDir, 0o777); err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(brightenDir, "proof.bin"))
	if err != nil {
		return err
	}
	defer proofFile.Close()

	n, err := proof.WriteTo(proofFile)
	if err != nil {
		return err
	}

	fmt.Println("Proof size: ", n)

	if config.markdownFile != "" {
		mdFile, err := os.OpenFile(config.markdownFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer mdFile.Close()

		if _, err = fmt.Fprintf(mdFile, "| %s | %f | %f | %d | %s |\n",
			fmt.Sprintf("%dx%d", len(finalPixels),
				len(finalPixels[0])),
			circuitCompilationDuration.Seconds(),
			provingDuration.Seconds(),
			n,
			config.backend,
		); err != nil {
			return err
		}
	}

	vkFile, err := os.Create(path.Join(brightenDir, "vkey.bin"))
	if err != nil {
		return err
	}
	defer vkFile.Close()

	_, err = vk.WriteTo(vkFile)
	if err != nil {
		return err
	}

	return nil
}

// generateBrightenProof returns the zk proof of brightening an image by a brightening factor.
func generateBrightenProof(backend string, original, brightened [][][]uint8) (io.WriterTo, io.WriterTo, time.Duration, time.Duration, error) {
	var circuit brightenCircuit
	circuit.Original = make([][][]frontend.Variable, len(original)) // First dimension
	for i := range original {
		circuit.Original[i] = make([][]frontend.Variable, len(original[i])) // Second dimension
		for j := range circuit.Original[i] {
			circuit.Original[i][j] = make([]frontend.Variable, len(original[i][j])) // Third dimension
		}
	}

	circuit.Brightened = make([][][]frontend.Variable, len(brightened)) // First dimension
	for i := range brightened {
		circuit.Brightened[i] = make([][]frontend.Variable, len(brightened[i])) // Second dimension
		for j := range circuit.Brightened[i] {
			circuit.Brightened[i][j] = make([]frontend.Variable, len(brightened[i][j])) // Third dimension
		}
	}

	t0 := time.Now()
	cs, err := compileCircuit(backend, &circuit)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	fmt.Println("Brighten compilation time:", time.Since(t0).Seconds())
	circuitCompilationDuration := time.Since(t0)

	t0 = time.Now()
	witness, err := frontend.NewWitness(&brightenCircuit{
		Original:   convertToFrontendVariable(original),
		Brightened: convertToFrontendVariable(brightened),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, 0, 0, err
	}

	proof, vk, err := generateProofByBackend(backend, cs, witness)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	fmt.Println("Time taken to prove: ", time.Since(t0).Seconds())
	proofDuration := time.Since(t0)

	return proof, vk, circuitCompilationDuration, proofDuration, nil
}

// brightenCircuit represents the arithmetic circuit to prove brighten transformations.
type brightenCircuit struct {
	Original   [][][]frontend.Variable `gnark:",secret"`
	Brightened [][][]frontend.Variable `gnark:",public"`
}

func (c *brightenCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(len(c.Original), len(c.Brightened))
	api.AssertIsEqual(len(c.Original[0]), len(c.Brightened[0]))
	api.AssertIsEqual(len(c.Original[0][0]), len(c.Brightened[0][0]))

	// The pixel values for the original and brightened images must match exactly.
	for i := 0; i < len(c.Original); i++ {
		for j := 0; j < len(c.Original[0]); j++ {
			r := api.Add(c.Original[i][j][0], brighteningFactor)
			r = api.Select(cmp.IsLess(api, r, MaxPixelValue), r, MaxPixelValue)
			r = api.Select(cmp.IsLess(api, r, MinPixelValue), MinPixelValue, r)

			g := api.Add(c.Original[i][j][1], brighteningFactor)
			g = api.Select(cmp.IsLess(api, g, MaxPixelValue), g, MaxPixelValue)
			g = api.Select(cmp.IsLess(api, g, MinPixelValue), MinPixelValue, g)

			b := api.Add(c.Original[i][j][2], brighteningFactor)
			b = api.Select(cmp.IsLess(api, b, MaxPixelValue), b, MaxPixelValue)
			b = api.Select(cmp.IsLess(api, b, MinPixelValue), MinPixelValue, b)

			api.AssertIsEqual(c.Brightened[i][j][0], r) // R
			api.AssertIsEqual(c.Brightened[i][j][1], g) // G
			api.AssertIsEqual(c.Brightened[i][j][2], b) // B
		}
	}

	return nil
}

// verifyBrightenConfig specifies the verification configuration for rotating an image by 90 degrees.
type verifyBrightenConfig struct {
	proofDir string
	finalImg string
}

// newVerifyBrightenCmd returns a new cobra.Command for rotating an image by 90 degrees.
func newVerifyBrightenCmd() *cobra.Command {
	var conf verifyBrightenConfig

	cmd := &cobra.Command{
		Use: "brighten",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyBrighten(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyBrighten verifies the zk proof of brightening an image by a brightening factor.
func verifyBrighten(config verifyBrightenConfig) error {
	// Open the final image file.
	finalImage, err := os.Open(config.finalImg)
	if err != nil {
		return err
	}
	defer finalImage.Close()

	// Get the pixel values for the final image.
	finalPixels, err := convertImgToPixels(finalImage)
	if err != nil {
		return err
	}

	witness, err := frontend.NewWitness(&brightenCircuit{
		Brightened: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	brightenDir := path.Join(config.proofDir, "brighten")
	if err = os.MkdirAll(brightenDir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(brightenDir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(brightenDir, "vkey.bin"))
	if err != nil {
		return err
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return err
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Println("Invalid proof 😞")
	} else {
		fmt.Println("Proof verified 🎉")
	}

	return nil
}
