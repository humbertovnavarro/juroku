package juroku

import (
	"image"

	"github.com/esimov/colorquant"
)

var dither map[string]colorquant.Dither = map[string]colorquant.Dither{
	"FloydSteinberg": {
		Filter: [][]float32{
			{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
			{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
			{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
		},
	},
	"Burkes": {
		Filter: [][]float32{
			{0.0, 0.0, 0.0, 8.0 / 32.0, 4.0 / 32.0},
			{2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{4.0 / 32.0, 8.0 / 32.0, 0.0, 0.0, 0.0},
		},
	},
	"Stucki": {
		Filter: [][]float32{
			{0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0},
			{2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0},
			{1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0},
		},
	},
	"Atkinson": {
		Filter: [][]float32{
			{0.0, 0.0, 1.0 / 8.0, 1.0 / 8.0},
			{1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0.0},
			{0.0, 1.0 / 8.0, 0.0, 0.0},
		},
	},
	"Sierra-3": {
		Filter: [][]float32{
			{0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0},
			{2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
			{0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0},
		},
	},
	"Sierra-2": {
		Filter: [][]float32{
			{0.0, 0.0, 0.0, 4.0 / 16.0, 3.0 / 16.0},
			{1.0 / 16.0, 2.0 / 16.0, 3.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
		},
	},
	"Sierra-Lite": {
		Filter: [][]float32{
			{0.0, 0.0, 2.0 / 4.0},
			{1.0 / 4.0, 1.0 / 4.0, 0.0},
			{0.0, 0.0, 0.0},
		},
	},
}

func Quantize(img image.Image, filterName string) image.Image {
	dst := image.NewPaletted(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()), ccPalette)
	if filter, ok := dither[filterName]; ok {
		return colorquant.Dither.Quantize(filter, img, dst, 16, true, true)
	}
	return colorquant.NoDither.Quantize(img, dst, 16, true, true)
}
