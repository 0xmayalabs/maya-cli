package cmd

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/spf13/cobra"
	"os"
	"path"
	"time"
)

// flipHorizontalConfig specifies the configuration for flipping an image horizontally.
type flipHorizontalConfig struct {
	originalImg  string
	finalImg     string
	proofDir     string
	markdownFile string
}

// newFlipHorizontalCmd returns a new cobra.Command for flipping an image horizontally.
func newFlipHorizontalCmd() *cobra.Command {
	var conf flipHorizontalConfig

	cmd := &cobra.Command{
		Use: "flip-horizontal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveFlipHorizontal(conf)
		},
	}

	bindFlipHorizontalFlags(cmd, &conf)

	return cmd
}

// bindFlipHorizontalFlags binds the flip horizontal configuration flags.
func bindFlipHorizontalFlags(cmd *cobra.Command, conf *flipHorizontalConfig) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
}

// proveFlipHorizontal generates the zk proof of flip horizontal transformation.
func proveFlipHorizontal(config flipHorizontalConfig) error {
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

	proof, vk, circuitCompilationDuration, provingDuration, err := generateFlipHorizontalProof(originalPixels, finalPixels)
	if err != nil {
		return err
	}

	flipHorizontalDir := path.Join(config.proofDir, "flipHorizontal")
	if err = os.MkdirAll(flipHorizontalDir, 0o777); err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(flipHorizontalDir, "proof.bin"))
	if err != nil {
		return err
	}
	defer proofFile.Close()

	n, err := proof.WriteTo(proofFile)
	if err != nil {
		return err
	}

	fmt.Println("Proof Size: ", n)

	vkFile, err := os.Create(path.Join(flipHorizontalDir, "vkey.bin"))
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

		if _, err = fmt.Fprintf(mdFile, "| %s | %f | %f | %d |\n",
			fmt.Sprintf("%dx%d", len(finalPixels),
				len(finalPixels[0])),
			circuitCompilationDuration.Seconds(),
			provingDuration.Seconds(),
			n,
		); err != nil {
			return err
		}
	}

	return nil
}

// generateFlipHorizontalProof returns the proof of flipHorizontal transformation.
func generateFlipHorizontalProof(original, flipped [][][]uint8) (groth16.Proof, groth16.VerifyingKey, time.Duration, time.Duration, error) {
	var circuit FlipHorizontalCircuit
	circuit.Original = make([][][]frontend.Variable, len(original)) // First dimension
	for i := range original {
		circuit.Original[i] = make([][]frontend.Variable, len(original[i])) // Second dimension
		for j := range circuit.Original[i] {
			circuit.Original[i][j] = make([]frontend.Variable, len(original[i][j])) // Third dimension
		}
	}

	circuit.Flipped = make([][][]frontend.Variable, len(flipped)) // First dimension
	for i := range flipped {
		circuit.Flipped[i] = make([][]frontend.Variable, len(flipped[i])) // Second dimension
		for j := range circuit.Flipped[i] {
			circuit.Flipped[i][j] = make([]frontend.Variable, len(flipped[i][j])) // Third dimension
		}
	}

	t0 := time.Now()
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("Flip Horizontal compilation time:", time.Since(t0).Seconds())
	circuitCompilationDuration := time.Since(t0)

	t0 = time.Now()
	witness, err := frontend.NewWitness(&FlipHorizontalCircuit{
		Original: convertToFrontendVariable(original),
		Flipped:  convertToFrontendVariable(flipped),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, 0, 0, err
	}

	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	fmt.Println("Time taken to prove: ", time.Since(t0).Seconds())
	proofDuration := time.Since(t0)

	return proof, vk, circuitCompilationDuration, proofDuration, nil
}

// FlipHorizontalCircuit represents the arithmetic circuit to prove flip horizontal transformations.
type FlipHorizontalCircuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Flipped  [][][]frontend.Variable `gnark:",public"`
}

func (c *FlipHorizontalCircuit) Define(api frontend.API) error {
	api.AssertIsDifferent(len(c.Original), 0)
	api.AssertIsDifferent(len(c.Flipped), 0)
	api.AssertIsEqual(len(c.Original), len(c.Flipped))
	api.AssertIsEqual(len(c.Original[0]), len(c.Flipped[0]))

	// The pixel values for the original and flipped images must match exactly.
	for i := 0; i < len(c.Original); i++ {
		for j := 0; j < len(c.Original[i]); j++ {
			j2 := len(c.Original[i]) - j - 1
			api.AssertIsEqual(c.Original[i][j][0], c.Flipped[i][j2][0]) // R
			api.AssertIsEqual(c.Original[i][j][1], c.Flipped[i][j2][1]) // G
			api.AssertIsEqual(c.Original[i][j][2], c.Flipped[i][j2][2]) // B
		}
	}

	return nil
}

// verifyFlipHorizontalConfig specifies the verification configuration for rotating an image by 270 degrees.
type verifyFlipHorizontalConfig struct {
	proofDir string
	finalImg string
}

// newVerifyFlipHorizontalCmd returns a new cobra.Command for rotating an image by 270 degrees.
func newVerifyFlipHorizontalCmd() *cobra.Command {
	var conf verifyFlipHorizontalConfig

	cmd := &cobra.Command{
		Use: "flip-horizontal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyFlipHorizontal(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyFlipHorizontal verifies the zk proof of flip horizontal transformation.
func verifyFlipHorizontal(config verifyFlipHorizontalConfig) error {
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

	witness, err := frontend.NewWitness(&FlipHorizontalCircuit{
		Flipped: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	flipHorizontalDir := path.Join(config.proofDir, "flipHorizontal")
	if err = os.MkdirAll(flipHorizontalDir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(flipHorizontalDir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(flipHorizontalDir, "vkey.bin"))
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
