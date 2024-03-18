package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
	"github.com/spf13/cobra"
	"image"
	"io"
	"os"
)

// New returns a new cobra command that handles maya cli commands and subcommands.
func New() *cobra.Command {
	return newRootCmd(
		newProveCmd(
			newCropCmd(),
			newRotate90Cmd(),
			newRotate180Cmd(),
			newRotate270Cmd(),
			newFlipVerticalCmd(),
			newFlipHorizontalCmd(),
			newBrightenCmd(),
		),
		newVerifyCmd(
			newVerifyCropCmd(),
			newVerifyRotate90Cmd(),
			newVerifyRotate180Cmd(),
			newVerifyRotate270Cmd(),
			newVerifyFlipVerticalCmd(),
			newVerifyFlipHorizontalCmd(),
			newVerifyBrightenCmd(),
		),
	)
}

func newRootCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "maya",
		Short: "Maya CLI",
		Long:  "Command line tool to create zero-knowledge proof of image transformations.",
	}

	root.AddCommand(cmds...)

	return root
}

func newProveCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "prove",
		Short: "Generates proof for the specified transformation.",
		Long:  "Generates zero knowledge proof of transformation on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}

func newVerifyCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "verify",
		Short: "Verifies proof for the specified transformation.",
		Long:  "Verifies the zero knowledge proof of transformation on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}

func compileCircuit(backend string, circuit frontend.Circuit) (constraint.ConstraintSystem, error) {
	switch backend {
	case "groth16":
		return frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	case "plonk": // TODO(dhruv): add plonkfri when its serialisation is supported.
		return frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, circuit)
	default:
		return nil, errors.New(fmt.Sprintf("invalid backend, %s", backend))
	}
}

func generateProofByBackend(backend string, cs constraint.ConstraintSystem, witness witness.Witness) (io.WriterTo, io.WriterTo, error) {
	switch backend {
	case "groth16":
		pk, vk, err := groth16.Setup(cs)
		if err != nil {
			return nil, nil, err
		}

		proof, err := groth16.Prove(cs, pk, witness)
		if err != nil {
			return nil, nil, err
		}

		return proof, vk, nil
	case "plonk":
		// TODO(dhruv): replace this with actual trusted setup ceremony.
		kzgSrs, err := test.NewKZGSRS(cs)
		if err != nil {
			return nil, nil, err
		}

		pk, vk, err := plonk.Setup(cs, kzgSrs)
		if err != nil {
			return nil, nil, err
		}

		proof, err := plonk.Prove(cs, pk, witness)
		if err != nil {
			return nil, nil, err
		}

		return proof, vk, nil
	default:
		return nil, nil, errors.New(fmt.Sprintf("invalid backend, %s", backend))
	}
}

// VerifyProofByBackend verifies the given proof by provided proof system backend.
func VerifyProofByBackend(backend, transformation string, proof, vk []byte, finalImg image.Image) error {
	pubWit, err := publicWitness(transformation, finalImg)
	if err != nil {
		return err
	}

	switch backend {
	case "groth16":
		grothProof := groth16.NewProof(ecc.BN254)
		_, err := grothProof.ReadFrom(bytes.NewBuffer(proof))
		if err != nil {
			return err
		}

		grothVk := groth16.NewVerifyingKey(ecc.BN254)
		_, err = grothVk.ReadFrom(bytes.NewBuffer(vk))
		if err != nil {
			return err
		}

		return groth16.Verify(grothProof, grothVk, pubWit)
	case "plonk":
		plonkProof := plonk.NewProof(ecc.BN254)
		_, err := plonkProof.ReadFrom(bytes.NewBuffer(proof))
		if err != nil {
			return err
		}

		plonkVk := plonk.NewVerifyingKey(ecc.BN254)
		_, err = plonkVk.ReadFrom(bytes.NewBuffer(vk))
		if err != nil {
			return err
		}

		return plonk.Verify(plonkProof, plonkVk, pubWit)
	default:
		return errors.New(fmt.Sprintf("invalid backend, %s", backend))
	}
}

// publicWitness returns public witness for the given transformation.
func publicWitness(transformation string, finalImg image.Image) (witness.Witness, error) {
	pixels, err := convertImgToPixels(finalImg)
	if err != nil {
		return nil, err
	}

	switch transformation {
	case "crop":
		wt, err := frontend.NewWitness(&CropCircuit{
			Cropped: convertToFrontendVariable(pixels),
		}, ecc.BN254.ScalarField())
		if err != nil {
			return nil, err
		}

		return wt.Public()
	case "flip_horizontal":
		wt, err := frontend.NewWitness(&FlipHorizontalCircuit{
			Flipped: convertToFrontendVariable(pixels),
		}, ecc.BN254.ScalarField())
		if err != nil {
			return nil, err
		}

		return wt.Public()
	default:
		return nil, errors.New("invalid transformation")
	}
}

func loadImage(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	return img, nil
}
