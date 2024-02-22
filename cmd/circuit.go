package cmd

import "github.com/consensys/gnark/frontend"

var heightStartNew, widthStartNew int

// Circuit represents the arithmetic circuit to prove crop transformations.
type Circuit struct {
	Original [][][]frontend.Variable `gnark:",secret"`
	Cropped  [][][]frontend.Variable `gnark:",public"`
}

func (c *Circuit) Define(api frontend.API) error {
	// The pixel values for the original and cropped images must match exactly.
	for i := 0; i < len(c.Cropped); i++ {
		for j := 0; j < len(c.Cropped[i]); j++ {
			api.AssertIsEqual(c.Cropped[i][j][0], c.Original[i+heightStartNew][j+widthStartNew][0]) // R
			api.AssertIsEqual(c.Cropped[i][j][1], c.Original[i+heightStartNew][j+widthStartNew][1]) // G
			api.AssertIsEqual(c.Cropped[i][j][2], c.Original[i+heightStartNew][j+widthStartNew][2]) // B
		}
	}

	return nil
}
