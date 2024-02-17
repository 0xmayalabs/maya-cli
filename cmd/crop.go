package cmd

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/spf13/cobra"
	"image/png"
	"io"
	"os"
	"path"
	"time"
)

// cropConfig specifies the configuration for cropping an image.
type cropConfig struct {
	originalImg    string
	croppedImg     string
	widthStartNew  uint
	heightStartNew uint
	proofDir       string
}

// newCropCmd returns a new cobra.Command for cropping.
func newCropCmd() *cobra.Command {
	var conf cropConfig

	cmd := &cobra.Command{
		Use: "crop",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveCrop(conf)
		},
	}

	bindFlags(cmd, &conf)

	return cmd
}

// bindFlags binds the crop configuration flags.
func bindFlags(cmd *cobra.Command, conf *cropConfig) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.croppedImg, "cropped-image", "", "The path to the cropped image. Supported image formats: PNG.")
	cmd.Flags().UintVar(&conf.widthStartNew, "width-start-new", 0, "The Original-coordinate for the top-left corner of the cropped image, relative to the original image's width.")
	cmd.Flags().UintVar(&conf.heightStartNew, "height-start-new", 0, "The Cropped-coordinate for the top-left corner of the cropped image, relative to the original image's height.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
}

// proveCrop generates the zk proof of crop transformation.
func proveCrop(config cropConfig) error {
	// Open the original image file.
	originalImage, err := os.Open(config.originalImg)
	if err != nil {
		return err
	}
	defer originalImage.Close()

	// Open the cropped image file.
	croppedImage, err := os.Open(config.croppedImg)
	if err != nil {
		return err
	}
	defer croppedImage.Close()

	// Get the pixel values for the original image.
	originalPixels, err := convertImgToPixels(originalImage)
	if err != nil {
		return err
	}

	// Get the pixel values for the cropped image.
	croppedPixels, err := convertImgToPixels(croppedImage)
	if err != nil {
		return err
	}

	proof, vk, err := generateProof(croppedPixels, originalPixels)
	if err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(config.proofDir, "proof.bin"))
	if err != nil {
		return err
	}
	defer proofFile.Close()

	_, err = proof.WriteTo(proofFile)
	if err != nil {
		return err
	}

	vkFile, err := os.Create(path.Join(config.proofDir, "vkey.bin"))
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

// convertImgToPixels returns a 3D array of pixel values for the provided image.
func convertImgToPixels(file io.Reader) ([][][]uint8, error) {
	// Decode the image.
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	// Get the image bounds.
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Printf("Image has width %d and height %d\n", width, height)

	// Create a 2D slice (which is effectively a 3D slice when considering RGB values).
	pixels := make([][][]uint8, height) // height x width x rgb
	for y := 0; y < height; y++ {
		pixels[y] = make([][]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Divide color values by 256 to scale from 0-65535 to 0-255
			pixels[y][x] = []uint8{uint8(r / 256), uint8(g / 256), uint8(b / 256)}
		}
	}

	return pixels, nil
}

// generateProof returns the proof of crop transformation.
func generateProof(cropped, original [][][]uint8) (groth16.Proof, groth16.VerifyingKey, error) {
	var circuit Circuit
	circuit.Original = make([][][]frontend.Variable, len(original)) // First dimension
	for i := range original {
		circuit.Original[i] = make([][]frontend.Variable, len(original[i])) // Second dimension
		for j := range circuit.Original[i] {
			circuit.Original[i][j] = make([]frontend.Variable, len(original[i][j])) // Third dimension
		}
	}

	circuit.Cropped = make([][][]frontend.Variable, len(cropped)) // First dimension
	for i := range cropped {
		circuit.Cropped[i] = make([][]frontend.Variable, len(cropped[i])) // Second dimension
		for j := range circuit.Cropped[i] {
			circuit.Cropped[i][j] = make([]frontend.Variable, len(cropped[i][j])) // Third dimension
		}
	}

	t0 := time.Now()
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("Circuit compilation time:", time.Since(t0).Seconds())

	witness, err := frontend.NewWitness(&Circuit{
		Original: convertToFrontendVariable(original),
		Cropped:  convertToFrontendVariable(cropped),
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

func convertToFrontendVariable(arr [][][]uint8) [][][]frontend.Variable {
	var resp [][][]frontend.Variable
	resp = make([][][]frontend.Variable, len(arr)) // First dimension
	for i := range arr {
		resp[i] = make([][]frontend.Variable, len(arr[i])) // Second dimension
		for j := range arr[i] {
			resp[i][j] = make([]frontend.Variable, len(arr[i][j])) // Third dimension
			for k := 0; k < 3; k++ {
				resp[i][j][k] = frontend.Variable(arr[i][j][k])
			}
		}
	}

	return resp
}
