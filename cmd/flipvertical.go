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

// flipVerticalConfig specifies the configuration for flipping an image vertically.
type flipVerticalConfig struct {
	originalImg string
	finalImg    string
	proofDir    string
}

// newFlipVerticalCmd returns a new cobra.Command for flipping an image vertically.
func newFlipVerticalCmd() *cobra.Command {
	var conf flipVerticalConfig

	cmd := &cobra.Command{
		Use: "flip-vertical",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveFlipVertical(conf)
		},
	}

	bindFlipVerticalFlags(cmd, &conf)

	return cmd
}

// bindFlipVerticalFlags binds the flip vertical configuration flags.
func bindFlipVerticalFlags(cmd *cobra.Command, conf *flipVerticalConfig) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
}

// proveFlipVertical generates the zk proof of flip vertical transformation.
func proveFlipVertical(config flipVerticalConfig) error {
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

	proof, vk, err := generateFlipVerticalProof(originalPixels, finalPixels)
	if err != nil {
		return err
	}

	flipVerticalDir := path.Join(config.proofDir, "flipVertical")
	if err = os.MkdirAll(flipVerticalDir, 0o777); err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(flipVerticalDir, "proof.bin"))
	if err != nil {
		return err
	}
	defer proofFile.Close()

	n, err := proof.WriteTo(proofFile)
	if err != nil {
		return err
	}

	fmt.Println("Proof Size: ", n)

	vkFile, err := os.Create(path.Join(flipVerticalDir, "vkey.bin"))
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

// generateFlipVerticalProof returns the proof of flipVertical transformation.
func generateFlipVerticalProof(original, flipped [][][]uint8) (groth16.Proof, groth16.VerifyingKey, error) {
	var circuit FlipVerticalCircuit
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

	fmt.Println("flipVertical compilation time:", time.Since(t0).Seconds())

	t0 = time.Now()
	witness, err := frontend.NewWitness(&FlipVerticalCircuit{
		Original: convertToFrontendVariable(original),
		Flipped:  convertToFrontendVariable(flipped),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}

	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return nil, nil, err
	}

	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Time taken to prove: ", time.Since(t0).Seconds())

	return proof, vk, nil
}

// FlipVerticalCircuit represents the arithmetic circuit to prove FlipVertical transformations.
type FlipVerticalCircuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Flipped  [][][]frontend.Variable `gnark:",public"`
}

func (c *FlipVerticalCircuit) Define(api frontend.API) error {
	api.AssertIsDifferent(len(c.Original), 0)
	api.AssertIsDifferent(len(c.Flipped), 0)
	api.AssertIsEqual(len(c.Original), len(c.Flipped))
	api.AssertIsEqual(len(c.Original[0]), len(c.Flipped[0]))

	// The pixel values for the original and flip vertical images must match exactly.
	for i := 0; i < len(c.Original); i++ {
		for j := 0; j < len(c.Original[i]); j++ {
			api.AssertIsEqual(c.Original[i][j][0], c.Flipped[len(c.Flipped)-i-1][j][0]) // R
			api.AssertIsEqual(c.Original[i][j][1], c.Flipped[len(c.Flipped)-i-1][j][1]) // G
			api.AssertIsEqual(c.Original[i][j][2], c.Flipped[len(c.Flipped)-i-1][j][2]) // B
		}
	}

	return nil
}

// verifyFlipVerticalConfig specifies the verification configuration for rotating an image by 270 degrees.
type verifyFlipVerticalConfig struct {
	proofDir string
	finalImg string
}

// newVerifyFlipVerticalCmd returns a new cobra.Command for rotating an image by 270 degrees.
func newVerifyFlipVerticalCmd() *cobra.Command {
	var conf verifyFlipVerticalConfig

	cmd := &cobra.Command{
		Use: "flip-vertical",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyFlipVertical(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyFlipVertical verifies the zk proof of flip vertical transformation.
func verifyFlipVertical(config verifyFlipVerticalConfig) error {
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

	witness, err := frontend.NewWitness(&FlipVerticalCircuit{
		Flipped: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	flipVerticalDir := path.Join(config.proofDir, "flipVertical")
	if err = os.MkdirAll(flipVerticalDir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(flipVerticalDir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(flipVerticalDir, "vkey.bin"))
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
