package cmd

import (
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
	"io"
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
