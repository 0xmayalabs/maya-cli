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

// rotate180Config specifies the configuration for rotating an image by 180 degrees.
type rotate180Config struct {
	originalImg  string
	finalImg     string
	proofDir     string
	markdownFile string
	backend      string
}

// newRotate180Cmd returns a new cobra.Command for rotating an image by 180 degrees.
func newRotate180Cmd() *cobra.Command {
	var conf rotate180Config

	cmd := &cobra.Command{
		Use: "rotate180",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveRotate180(conf)
		},
	}

	bindRotate180Flags(cmd, &conf)

	return cmd
}

// bindRotate180Flags binds the rotate180 configuration flags.
func bindRotate180Flags(cmd *cobra.Command, conf *rotate180Config) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.backend, "backend", "groth16", "The proving backend used for generating the proofs.")
}

// proveRotate180 generates the zk proof of rotated transformation 180.
func proveRotate180(config rotate180Config) error {
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

	proof, vk, circuitCompilationDuration, provingDuration, err := generateRotate180Proof(config.backend, originalPixels, finalPixels)
	if err != nil {
		return err
	}

	rotate90Dir := path.Join(config.proofDir, "rotate180")
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

	vkFile, err := os.Create(path.Join(rotate90Dir, "vkey.bin"))
	if err != nil {
		return err
	}
	defer vkFile.Close()

	_, err = vk.WriteTo(vkFile)
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

		if _, err = fmt.Fprintf(mdFile, "| %s | %f | %f | %d | %s|\n",
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

	return nil
}

// generateRotate180Proof returns the proof of rotate180 transformation.
func generateRotate180Proof(backend string, original, rotated [][][]uint8) (io.WriterTo, io.WriterTo, time.Duration, time.Duration, error) {
	var circuit Rotate180Circuit
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

	fmt.Println("Rotate180 compilation time:", time.Since(t0).Seconds())
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

// Rotate180Circuit represents the arithmetic circuit to prove rotate180 transformations.
type Rotate180Circuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Rotated  [][][]frontend.Variable `gnark:",public"`
}

func (c *Rotate180Circuit) Define(api frontend.API) error {
	// TODO(dhruv): Add AssertIsDifferent to compare len(Original) with 0.
	api.AssertIsEqual(len(c.Original), len(c.Rotated[0]))
	api.AssertIsEqual(len(c.Original[0]), len(c.Rotated))

	rows := len(c.Original)
	cols := len(c.Original[0])

	// The pixel values for the original and rotated180 images must match exactly.
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			api.AssertIsEqual(c.Original[i][j][0], c.Rotated[rows-1-i][cols-1-j][0]) // R
			api.AssertIsEqual(c.Original[i][j][1], c.Rotated[rows-1-i][cols-1-j][1]) // G
			api.AssertIsEqual(c.Original[i][j][2], c.Rotated[rows-1-i][cols-1-j][2]) // B
		}
	}

	return nil
}

// verifyRotate180Config specifies the verification configuration for rotating an image by 180 degrees.
type verifyRotate180Config struct {
	proofDir string
	finalImg string
}

// newVerifyRotate180Cmd returns a new cobra.Command for rotating an image by 180 degrees.
func newVerifyRotate180Cmd() *cobra.Command {
	var conf verifyRotate180Config

	cmd := &cobra.Command{
		Use: "rotate180",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyRotate180Crop(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyRotate180Crop verifies the zk proof of rotate180 transformation.
func verifyRotate180Crop(config verifyRotate180Config) error {
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

	witness, err := frontend.NewWitness(&Rotate180Circuit{
		Rotated: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	rotate180Dir := path.Join(config.proofDir, "rotate180")
	if err = os.MkdirAll(rotate180Dir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(rotate180Dir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(rotate180Dir, "vkey.bin"))
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
