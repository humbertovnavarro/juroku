package juroku

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
	"text/template"
)

// GenerateCode generates the ComputerCraft code to render the given image
// that must have an underlying color type of color.RGBA and is assumed to have
// ChunkImage already called on it.
func GenerateCode(img image.Image) ([]byte, error) {
	palette := GetPalette(img)
	if len(palette) > 16 {
		return nil, errors.New("juroku: palette must have <= 16 colors")
	}

	colorsCodes := "0123456789abcdef"

	paletteToColor := make(map[color.RGBA]byte)
	for i, col := range palette {
		paletteToColor[col.(color.RGBA)] = colorsCodes[i]
	}

	type rowData struct {
		Text      []string
		TextColor []byte
		BgColor   []byte
	}

	var rows []rowData

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 3 {
		var row rowData
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x += 2 {
			chunk := make([]byte, 0, 6)
			for dy := 0; dy < 3; dy++ {
				for dx := 0; dx < 2; dx++ {
					chunk = append(chunk,
						paletteToColor[img.At(x+dx, y+dy).(color.RGBA)])
				}
			}

			text, textColor, bgColor := chunkToBlit(chunk)
			row.Text = append(row.Text, text)
			row.TextColor = append(row.TextColor, textColor)
			row.BgColor = append(row.BgColor, bgColor)
		}
		rows = append(rows, row)
	}

	buf := new(bytes.Buffer)
	err := cc.Execute(buf, struct {
		Rows    []rowData
		Palette color.Palette
		Width   int
		Height  int
	}{
		Rows:    rows,
		Palette: palette,
		Width:   img.Bounds().Dx() / 2,
		Height:  img.Bounds().Dy() / 3,
	})

	return buf.Bytes(), err
}

func chunkToBlit(chunk []byte) (text string, textColor byte, bgColor byte) {
	bgColor = chunk[5]

	var b byte
	var i uint
	for i = 0; i < 6; i++ {
		if chunk[i] != bgColor {
			textColor = chunk[i]
			b |= 1 << i
		} else {
			b |= 0 << i
		}
	}

	if textColor == 0 {
		textColor = '0'
	}

	text = "\\" + strconv.Itoa(int(b)+128)

	return
}

var cc = template.Must(template.New("cc").Funcs(template.FuncMap{
	"colorToHex": func(c color.Color) string {
		r, g, b, _ := c.RGBA()
		return fmt.Sprintf("%02X%02X%02X", r>>8, g>>8, b>>8)
	},
	"strJoin": func(strs []string) string {
		return strings.Join(strs, "")
	},
	"bToString": func(b []byte) string {
		return string(b)
	},
}).Parse(`-- This code was automatically generated by
--        _                  __
--       (_)_  ___________  / /____  __
--      / / / / / ___/ __ \/ //_/ / / /
--     / / /_/ / /  / /_/ / ,< / /_/ /
--  __/ /\__,_/_/   \____/_/|_|\__,_/
-- /___/  by 1lann - github.com/tmpim/juroku
--
-- Usage:
-- local image = require("image")
-- image.draw(term) or image.draw(monitor)
local image = {}
function image.draw(t)
	local x, y = t.getCursorPos()
	{{range $index, $color := .Palette -}}
	t.setPaletteColor(2^{{$index}}, 0x{{colorToHex $color}})
	{{end}}
	{{range $index, $row := .Rows -}}
	t.setCursorPos(x, y + {{$index}})
	t.blit("{{strJoin $row.Text}}", "{{bToString $row.TextColor}}", "{{bToString $row.BgColor}}")
	{{end}}
end
function image.getColors()
	local colors = {}
	{{range $index, $color := .Palette -}}
	table.insert(colors, 2^{{$index}}, 0x{{colorToHex $color}})
	{{end}}

	return colors
end
function image.getSize()
	return {{.Width}}, {{.Height}}
end
return image
`))
