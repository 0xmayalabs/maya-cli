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

// rotate90Config specifies the configuration for rotating an image by 90 degrees.
type rotate90Config struct {
	originalImg string
	finalImg    string
	proofDir    string
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
}

// proveRotate90 generates the zk proof of rotate 90 transformation.
func proveRotate90(config rotate90Config) error {
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

	proof, vk, err := generateRotate90Proof(originalPixels, finalPixels)
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
func generateRotate90Proof(original, rotated [][][]uint8) (groth16.Proof, groth16.VerifyingKey, error) {
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
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("Rotate90 compilation time:", time.Since(t0).Seconds())

	t0 = time.Now()
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

	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Time taken to prove: ", time.Since(t0).Seconds())

	return proof, vk, nil
}

// Rotate90Circuit represents the arithmetic circuit to prove rotate90 transformations.
type Rotate90Circuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Rotated  [][][]frontend.Variable `gnark:",public"`
}

func (c *Rotate90Circuit) Define(api frontend.API) error {
	api.AssertIsDifferent(len(c.Original), 0)
	api.AssertIsDifferent(len(c.Rotated), 0)
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
