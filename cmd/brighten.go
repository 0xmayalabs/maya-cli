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

var brighteningFactor int

// brightenConfig specifies the configuration for brightening an image by a brightening factor.
type brightenConfig struct {
	originalImg       string
	finalImg          string
	brighteningFactor int // TODO(xenowits): Convert it to floating-point
	proofDir          string
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

	bindbrightenFlags(cmd, &conf)

	return cmd
}

// bindbrightenFlags binds the brighten configuration flags.
func bindbrightenFlags(cmd *cobra.Command, conf *brightenConfig) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
}

// proveBrighten generates the zk proof of brightening an image by a brightening factor.
func proveBrighten(config brightenConfig) error {
	brighteningFactor = config.brighteningFactor

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

	proof, vk, err := generateBrightenProof(originalPixels, finalPixels)
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
func generateBrightenProof(original, brightened [][][]uint8) (groth16.Proof, groth16.VerifyingKey, error) {
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
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("Brighten compilation time:", time.Since(t0).Seconds())

	t0 = time.Now()
	witness, err := frontend.NewWitness(&brightenCircuit{
		Original:   convertToFrontendVariable(original),
		Brightened: convertToFrontendVariable(brightened),
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

// brightenCircuit represents the arithmetic circuit to prove brighten transformations.
type brightenCircuit struct {
	Original   [][][]frontend.Variable `gnark:",secret"`
	Brightened [][][]frontend.Variable `gnark:",public"`
}

func (c *brightenCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(len(c.Original), len(c.Brightened))

	// The pixel values for the original and brightened images must match exactly.
	for i := 0; i < len(c.Original[0]); i++ { // Columns
		for j := 0; j < len(c.Original); j++ { // Rows
			api.AssertIsEqual(len(c.Original[j][i]), len(c.Brightened[j][i]))

			api.AssertIsEqual(c.Brightened[j][i][0], api.Mul(c.Original[j][i][0], brighteningFactor)) // R
			api.AssertIsEqual(c.Brightened[j][i][1], api.Mul(c.Original[j][i][1], brighteningFactor)) // G
			api.AssertIsEqual(c.Brightened[j][i][2], api.Mul(c.Original[j][i][2], brighteningFactor)) // B
		}
	}

	return nil
}

func isBrightenedVersion(arr1, arr2 [][][]float64, value float64) bool {
	// Check if the dimensions of both images are the same
	if len(arr1) != len(arr2) {
		return false
	}
	for x := range arr1[0] { // Swap: Iterate over x first
		for y := range arr1 { // Swap: Then iterate over y
			if len(arr1[y][x]) != len(arr2[y][x]) {
				return false
			}
			for c := range arr1[y][x] {
				// Check if arr2 is a brightened version of arr1 by the factor "value"
				if arr2[y][x][c] != arr1[y][x][c]*value {
					return false
				}
			}
		}
	}
	return true
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
