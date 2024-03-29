package cmd

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"time"
)

// rotate90Config specifies the configuration for rotating an image by 90 degrees.
type rotate90Config struct {
	originalImg  string
	finalImg     string
	proofDir     string
	markdownFile string
	backend      string
}

// newRotate90Cmd returns a new cobra.Command for rotating an image by 90 degrees.
func newRotate90Cmd() *cobra.Command {
	var conf rotate90Config

	cmd := &cobra.Command{
		Use: "rotate90",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveRotate90(conf)
		},
	}

	bindRotate90Flags(cmd, &conf)

	return cmd
}

// bindRotate90Flags binds the rotate90 configuration flags.
func bindRotate90Flags(cmd *cobra.Command, conf *rotate90Config) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.backend, "backend", "groth16", "The proving backend used for generating the proofs.")
}

// proveRotate90 generates the zk proof of rotate 90 transformation.
func proveRotate90(config rotate90Config) error {
	// Open the original image file.
	originalImage, err := loadImage(config.originalImg)
	if err != nil {
		return err
	}

	// Open the final image file.
	finalImage, err := loadImage(config.finalImg)
	if err != nil {
		return err
	}

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

	proof, vk, circuitCompilationDuration, provingDuration, err := generateRotate90Proof(config.backend, originalPixels, finalPixels)
	if err != nil {
		return err
	}

	rotate90Dir := path.Join(config.proofDir, "rotate90")
	if err = os.MkdirAll(rotate90Dir, 0o777); err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(rotate90Dir, "proof.bin"))
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

	vkFile, err := os.Create(path.Join(rotate90Dir, "vkey.bin"))
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

// generateRotate90Proof returns the proof of rotate90 transformation.
func generateRotate90Proof(backend string, original, rotated [][][]uint8) (io.WriterTo, io.WriterTo, time.Duration, time.Duration, error) {
	var circuit Rotate90Circuit
	circuit.Original = make([][][]frontend.Variable, len(original)) // First dimension
	for i := range original {
		circuit.Original[i] = make([][]frontend.Variable, len(original[i])) // Second dimension
		for j := range circuit.Original[i] {
			circuit.Original[i][j] = make([]frontend.Variable, len(original[i][j])) // Third dimension
		}
	}

	circuit.Rotated = make([][][]frontend.Variable, len(rotated)) // First dimension
	for i := range rotated {
		circuit.Rotated[i] = make([][]frontend.Variable, len(rotated[i])) // Second dimension
		for j := range circuit.Rotated[i] {
			circuit.Rotated[i][j] = make([]frontend.Variable, len(rotated[i][j])) // Third dimension
		}
	}

	t0 := time.Now()
	cs, err := compileCircuit(backend, &circuit)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	fmt.Println("Rotate90 compilation time:", time.Since(t0).Seconds())
	circuitCompilationDuration := time.Since(t0)

	t0 = time.Now()
	witness, err := frontend.NewWitness(&Rotate90Circuit{
		Original: convertToFrontendVariable(original),
		Rotated:  convertToFrontendVariable(rotated),
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

// Rotate90Circuit represents the arithmetic circuit to prove rotate90 transformations.
type Rotate90Circuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Rotated  [][][]frontend.Variable `gnark:",public"`
}

func (c *Rotate90Circuit) Define(api frontend.API) error {
	// TODO(dhruv): Add AssertIsDifferent to compare len(Original) with 0.
	api.AssertIsEqual(len(c.Original), len(c.Rotated[0]))
	api.AssertIsEqual(len(c.Original[0]), len(c.Rotated))

	// The pixel values for the original and rotated90 images must match exactly.
	for i := 0; i < len(c.Original); i++ {
		for j := 0; j < len(c.Original[i]); j++ {
			api.AssertIsEqual(c.Original[i][j][0], c.Rotated[j][len(c.Original)-1-i][0]) // R
			api.AssertIsEqual(c.Original[i][j][1], c.Rotated[j][len(c.Original)-1-i][1]) // G
			api.AssertIsEqual(c.Original[i][j][2], c.Rotated[j][len(c.Original)-1-i][2]) // B
		}
	}

	return nil
}

// verifyRotate90Config specifies the verification configuration for rotating an image by 90 degrees.
type verifyRotate90Config struct {
	proofDir string
	finalImg string
}

// newVerifyRotate90Cmd returns a new cobra.Command for rotating an image by 90 degrees.
func newVerifyRotate90Cmd() *cobra.Command {
	var conf verifyRotate90Config

	cmd := &cobra.Command{
		Use: "rotate90",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyRotate90(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyRotate90 verifies the zk proof of rotate90 transformation.
func verifyRotate90(config verifyRotate90Config) error {
	// Open the final image file.
	finalImage, err := loadImage(config.finalImg)
	if err != nil {
		return err
	}

	// Get the pixel values for the final image.
	finalPixels, err := convertImgToPixels(finalImage)
	if err != nil {
		return err
	}

	witness, err := frontend.NewWitness(&Rotate90Circuit{
		Rotated: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	rotate90Dir := path.Join(config.proofDir, "rotate90")
	if err = os.MkdirAll(rotate90Dir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(rotate90Dir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(rotate90Dir, "vkey.bin"))
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
