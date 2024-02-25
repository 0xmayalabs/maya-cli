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

var contrastFactor int

// contrastConfig specifies the configuration for contrasting an image by a contrasting factor.
type contrastConfig struct {
	originalImg    string
	finalImg       string
	contrastFactor int // TODO(xenowits): Convert it to floating-point
	proofDir       string
}

// newContrastCmd returns a new cobra.Command for contrasting an image by a contrasting factor.
func newContrastCmd() *cobra.Command {
	var conf contrastConfig

	cmd := &cobra.Command{
		Use: "contrast",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proveContrast(conf)
		},
	}

	bindContrastFlags(cmd, &conf)

	return cmd
}

// bindContrastFlags binds the contrasting configuration flags.
func bindContrastFlags(cmd *cobra.Command, conf *contrastConfig) {
	cmd.Flags().StringVar(&conf.originalImg, "original-image", "", "The path to the original image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")
	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().IntVar(&conf.contrastFactor, "contrasting-factor", 2, "The factor with which image is Contrasted.")
}

// proveContrast generates the zk proof of contrasting an image by a contrasting factor.
func proveContrast(config contrastConfig) error {
	contrastFactor = config.contrastFactor

	fmt.Println("Contrast factor", contrastFactor)

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

	proof, vk, err := generateContrastProof(originalPixels, finalPixels)
	if err != nil {
		return err
	}

	contrastDir := path.Join(config.proofDir, "contrast")
	if err = os.MkdirAll(contrastDir, 0o777); err != nil {
		return err
	}

	proofFile, err := os.Create(path.Join(contrastDir, "proof.bin"))
	if err != nil {
		return err
	}
	defer proofFile.Close()

	n, err := proof.WriteTo(proofFile)
	if err != nil {
		return err
	}

	fmt.Println("Proof size: ", n)

	vkFile, err := os.Create(path.Join(contrastDir, "vkey.bin"))
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

// generateContrastProof returns the zk proof of contrasting an image by a contrasting factor.
func generateContrastProof(original, contrasted [][][]uint8) (groth16.Proof, groth16.VerifyingKey, error) {
	var circuit ContrastCircuit
	circuit.Original = make([][][]frontend.Variable, len(original)) // First dimension
	for i := range original {
		circuit.Original[i] = make([][]frontend.Variable, len(original[i])) // Second dimension
		for j := range circuit.Original[i] {
			circuit.Original[i][j] = make([]frontend.Variable, len(original[i][j])) // Third dimension
		}
	}

	circuit.Contrasted = make([][][]frontend.Variable, len(contrasted)) // First dimension
	for i := range contrasted {
		circuit.Contrasted[i] = make([][]frontend.Variable, len(contrasted[i])) // Second dimension
		for j := range circuit.Contrasted[i] {
			circuit.Contrasted[i][j] = make([]frontend.Variable, len(contrasted[i][j])) // Third dimension
		}
	}

	t0 := time.Now()
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("contrast compilation time:", time.Since(t0).Seconds())

	t0 = time.Now()
	witness, err := frontend.NewWitness(&ContrastCircuit{
		Original:   convertToFrontendVariable(original),
		Contrasted: convertToFrontendVariable(contrasted),
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

// ContrastCircuit represents the arithmetic circuit to prove contrast transformation.
type ContrastCircuit struct {
	Original   [][][]frontend.Variable `gnark:",secret"`
	Contrasted [][][]frontend.Variable `gnark:",public"`
}

func (c *ContrastCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(len(c.Original), len(c.Contrasted))
	api.AssertIsEqual(len(c.Original[0]), len(c.Contrasted[0]))
	api.AssertIsEqual(len(c.Original[0][0]), len(c.Contrasted[0][0]))

	// The pixel values for the original and Contrasted images must match exactly.
	// for i := 0; i < len(c.Original); i++ {
	// 	for j := 0; j < len(c.Original[0]); j++ {
	// 		r := api.Add(c.Original[i][j][0], contrastFactor)
	// 		r = api.Select(cmp.IsLess(api, r, MaxPixelValue), r, MaxPixelValue)
	// 		r = api.Select(cmp.IsLess(api, r, MinPixelValue), MinPixelValue, r)
	//
	// 		g := api.Add(c.Original[i][j][1], contrastFactor)
	// 		g = api.Select(cmp.IsLess(api, g, MaxPixelValue), g, MaxPixelValue)
	// 		g = api.Select(cmp.IsLess(api, g, MinPixelValue), MinPixelValue, g)
	//
	// 		b := api.Add(c.Original[i][j][2], contrastFactor)
	// 		b = api.Select(cmp.IsLess(api, b, MaxPixelValue), b, MaxPixelValue)
	// 		b = api.Select(cmp.IsLess(api, b, MinPixelValue), MinPixelValue, b)
	//
	// 		api.AssertIsEqual(c.Contrasted[i][j][0], r) // R
	// 		api.AssertIsEqual(c.Contrasted[i][j][1], g) // G
	// 		api.AssertIsEqual(c.Contrasted[i][j][2], b) // B
	// 	}
	// }

	return nil
}

// RUST: image-0.24.8/src/imageops/colorops.rs:89
/**
  let max = S::DEFAULT_MAX_VALUE;
  let max: f32 = NumCast::from(max).unwrap();

  let percent = ((100.0 + contrast) / 100.0).powi(2);

  for (x, y, pixel) in image.pixels() {
      let f = pixel.map(|b| {
          let c: f32 = NumCast::from(b).unwrap();

          let d = ((c / max - 0.5) * percent + 0.5) * max;
          let e = clamp(d, 0.0, max);

          NumCast::from(e).unwrap()
      });
      out.put_pixel(x, y, f);
  }
*/

// verifyContrastConfig specifies the verification configuration for contrasting an image by 90 degrees.
type verifyContrastConfig struct {
	proofDir string
	finalImg string
}

// newVerifyContrastCmd returns a new cobra.Command for contrasting an image.
func newVerifyContrastCmd() *cobra.Command {
	var conf verifyContrastConfig

	cmd := &cobra.Command{
		Use: "contrast",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyContrast(conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")
	cmd.Flags().StringVar(&conf.finalImg, "final-image", "", "The path to the final image. Supported image formats: PNG.")

	return cmd
}

// verifyContrast verifies the zk proof of contrasting an image by a contrasting factor.
func verifyContrast(config verifyContrastConfig) error {
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

	witness, err := frontend.NewWitness(&ContrastCircuit{
		Contrasted: convertToFrontendVariable(finalPixels),
	}, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	contrastDir := path.Join(config.proofDir, "contrast")
	if err = os.MkdirAll(contrastDir, 0o777); err != nil {
		return err
	}

	proof, err := readProof(path.Join(contrastDir, "proof.bin"))
	if err != nil {
		return err
	}

	vk, err := readVerifyingKey(path.Join(contrastDir, "vkey.bin"))
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
