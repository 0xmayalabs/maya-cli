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

// rotate180Config specifies the configuration for rotating an image by 180 degrees.
type rotate180Config struct {
	originalImg string
	finalImg    string
	proofDir    string
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
}

// proveRotate180 generates the zk proof of rotated transformation 180.
func proveRotate180(config rotate180Config) error {
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

	proof, vk, err := generateRotate180Proof(originalPixels, finalPixels)
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

	_, err = proof.WriteTo(proofFile)
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

	return nil
}

// generateRotate180Proof returns the proof of rotate180 transformation.
func generateRotate180Proof(original, rotated [][][]uint8) (groth16.Proof, groth16.VerifyingKey, error) {
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
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("Circuit compilation time:", time.Since(t0).Seconds())

	witness, err := frontend.NewWitness(&Rotate90Circuit{
		Original: convertToFrontendVariable(original),
		Rotated:  convertToFrontendVariable(rotated),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}

	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return nil, nil, err
	}

	t0 = time.Now()
	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Time taken to prove: ", time.Since(t0).Seconds())

	return proof, vk, nil
}

// Rotate180Circuit represents the arithmetic circuit to prove rotate180 transformations.
type Rotate180Circuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Rotated  [][][]frontend.Variable `gnark:",public"`
}

func (c *Rotate180Circuit) Define(api frontend.API) error {
	api.AssertIsDifferent(len(c.Original), 0)
	api.AssertIsDifferent(len(c.Rotated), 0)
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

func is180DegreeRotation(arr1, arr2 [][][]int) bool {
	if len(arr1) == 0 || len(arr2) == 0 || len(arr1) != len(arr2) || len(arr1[0]) != len(arr2[0]) {
		return false
	}

	rows := len(arr1)
	cols := len(arr1[0])

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if arr1[i][j][0] != arr2[rows-1-i][cols-1-j][0] { // Assuming 3rd dimension has at least 1 element
				return false
			}
		}
	}

	return true
}

// VERIFICATION CODE

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
