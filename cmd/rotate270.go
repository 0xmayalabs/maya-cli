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

// rotate270Config specifies the configuration for rotating an image by 270 degrees.
type rotate270Config struct {
	originalImg  string
	finalImg     string
	proofDir     string
	markdownFile string
	backend      string
}

// newRotate270Cmd returns a new cobra.Command for rotating an image by 270 degrees.
func newRotate270Cmd() *cobra.Command {
	var conf rotate270Config

	cmd := &cobra.Command{
		Use: "rotate270",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveRotate270(conf)
		},
	}

	bindRotate270Flags(cmd, &conf)

	return cmd
}

// bindRotate270Flags binds the rotate270 configuration flags.
func bindRotate270Flags(cmd *cobra.Command, conf *rotate270Config) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
}

// proveRotate270 generates the zk proof of rotated transformation 270.
func proveRotate270(config rotate270Config) error {
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

	proof, vk, circuitCompilationDuration, provingDuration, err := generateRotate270Proof(config.backend, originalPixels, finalPixels)
	if err != nil {
		return err
	}

	rotate270Dir := path.Join(config.proofDir, "rotate270")
	if err = os.MkdirAll(rotate270Dir, 0o777); err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(rotate270Dir, "proof.bin"))
	if err != nil {
		return err
	}
	defer proofFile.Close()

	n, err := proof.WriteTo(proofFile)
	if err != nil {
		return err
	}

	fmt.Println("Proof size: ", n)

	vkFile, err := os.Create(path.Join(rotate270Dir, "vkey.bin"))
	if err != nil {
		return err
	}
	defer vkFile.Close()

	_, err = vk.WriteTo(vkFile)
	if err != nil {
		return err
	}

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

	return nil
}

// generateRotate270Proof returns the proof of rotate270 transformation.
func generateRotate270Proof(backend string, original, rotated [][][]uint8) (io.WriterTo, io.WriterTo, time.Duration, time.Duration, error) {
	var circuit Rotate270Circuit
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

	fmt.Println("Rotate270Circuit compilation time:", time.Since(t0).Seconds())
	circuitCompilationDuration := time.Since(t0)

	t0 = time.Now()
	witness, err := frontend.NewWitness(&Rotate270Circuit{
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

// Rotate270Circuit represents the arithmetic circuit to prove rotate270 transformations.
type Rotate270Circuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Rotated  [][][]frontend.Variable `gnark:",public"`
}

func (c *Rotate270Circuit) Define(api frontend.API) error {
	// TODO(dhruv): Add AssertIsDifferent to compare len(Original) with 0.
	api.AssertIsEqual(len(c.Original[0]), len(c.Rotated))
	api.AssertIsEqual(len(c.Original), len(c.Rotated[0]))

	// The pixel values for the original and rotated270 images must match exactly.
	for i := 0; i < len(c.Original); i++ {
		for j := 0; j < len(c.Original[i]); j++ {
			api.AssertIsEqual(c.Original[i][j][0], c.Rotated[len(c.Rotated)-j-1][i][0]) // R
			api.AssertIsEqual(c.Original[i][j][1], c.Rotated[len(c.Rotated)-j-1][i][1]) // G
			api.AssertIsEqual(c.Original[i][j][2], c.Rotated[len(c.Rotated)-j-1][i][2]) // B
		}
	}

	return nil
}

// verifyRotate270Config specifies the verification configuration for rotating an image by 270 degrees.
type verifyRotate270Config struct {
	proofDir string
	finalImg string
}

// newVerifyRotate270Cmd returns a new cobra.Command for rotating an image by 270 degrees.
func newVerifyRotate270Cmd() *cobra.Command {
	var conf verifyRotate270Config

	cmd := &cobra.Command{
		Use: "rotate270",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyRotate270(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyRotate270 verifies the zk proof of rotate270 transformation.
func verifyRotate270(config verifyRotate270Config) error {
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

	witness, err := frontend.NewWitness(&Rotate270Circuit{
		Rotated: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	rotate270Dir := path.Join(config.proofDir, "rotate270")
	if err = os.MkdirAll(rotate270Dir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(rotate270Dir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(rotate270Dir, "vkey.bin"))
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
