package cmd

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/spf13/cobra"
	"os"
	"path"
)

func newVerifyCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "verify",
		Short: "Verifies proof for the specified edit.",
		Long:  "Verifies the zero knowledge proof of edit on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}

// verifyCropConfig specifies the verification configuration for cropping an image.
type verifyCropConfig struct {
	proofDir   string
	croppedImg string
}

// newVerifyCropCmd returns a new cobra.Command for cropping.
func newVerifyCropCmd() *cobra.Command {
	var conf verifyCropConfig

	cmd := &cobra.Command{
		Use: "crop",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyCrop(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.croppedImg, "cropped-image", "", "The path to the cropped image. Supported image formats: PNG.")

	return cmd
}

// verifyCrop verifies the zk proof of crop transformation.
func verifyCrop(config verifyCropConfig) error {
	// Open the cropped image file.
	croppedImage, err := os.Open(config.croppedImg)
	if err != nil {
		return err
	}
	defer croppedImage.Close()

	// Get the pixel values for the cropped image.
	croppedPixels, err := convertImgToPixels(croppedImage)
	if err != nil {
		return err
	}

	witness, err := frontend.NewWitness(&Circuit{
		Cropped: convertToFrontendVariable(croppedPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	proof, err := readProof(path.Join(config.proofDir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(config.proofDir, "vkey.bin"))
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

// readProof returns the zk proof by reading it from the disk.
func readProof(proofPath string) (groth16.Proof, error) {
	file, err := os.Open(proofPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resp := groth16.NewProof(ecc.BN254)
	_, err = resp.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// readVerifyingKey returns the verifying key by reading it from the disk.
func readVerifyingKey(verifyingKeyPath string) (groth16.VerifyingKey, error) {
	file, err := os.Open(verifyingKeyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resp := groth16.NewVerifyingKey(ecc.BN254)
	_, err = resp.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
