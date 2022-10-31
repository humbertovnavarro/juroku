package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"math"
	"strconv"
	"syscall/js"

	"github.com/anthonynsimon/bild/transform"
	"github.com/humbertovnavarro/juroku/pkg/juroku"
)

func main() {
	lua := ""
	c := make(chan struct{})
	println("Go WebAssembly Initialized")
	document := js.Global().Get("document")
	size := document.Call("getElementById", "size")
	copy := document.Call("getElementById", "copy")
	clipboard := js.Global().Get("navigator").Get("clipboard")
	scale := document.Call("getElementById", "scale")
	preview := document.Call("getElementById", "preview")
	convert := document.Call("getElementById", "convert")
	filter := document.Call("getElementById", "filter")
	copy.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) any {
		if lua == "" {
			return nil
		}
		clipboard.Call("writeText", lua)
		js.Global().Call("alert", "lua code copied to clipboard")
		return nil
	}))

	fileInput := document.Call("getElementById", "image-file")

	// file input change handler
	fileInput.Set("onchange", js.FuncOf(func(v js.Value, x []js.Value) any {
		if fileInput.Get("files").Call("item", 0).IsUndefined() {
			return nil
		}
		preview.Get("classList").Call("remove", "invisible")
		preview.Set("src", "https://raw.githubusercontent.com/SamHerbert/SVG-Loaders/master/svg-loaders/spinning-circles.svg")
		fileInput.Get("files").Call("item", 0).Call("arrayBuffer").Call("then", js.FuncOf(func(v js.Value, x []js.Value) any {
			data := js.Global().Get("Uint8Array").New(x[0])
			dst := make([]byte, data.Get("length").Int())
			js.CopyBytesToGo(dst, data)
			preview.Set("src", fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(dst)))
			return nil
		}))
		return nil
	}))

	// convert button callback
	convert.Set("onclick", js.FuncOf(func(v js.Value, x []js.Value) any {
		if fileInput.Get("files").Call("item", 0).IsUndefined() {
			js.Global().Call("alert", "you need to upload an image first!")
			return nil
		}
		filterName := filter.Get("value").String()
		fileInput.Get("files").Call("item", 0).Call("arrayBuffer").Call("then", js.FuncOf(func(v js.Value, x []js.Value) any {
			data := js.Global().Get("Uint8Array").New(x[0])
			dst := make([]byte, data.Get("length").Int())
			js.CopyBytesToGo(dst, data)
			img, _, err := image.Decode(bytes.NewReader(dst))
			if err != nil {
				js.Global().Call("alert", err.Error())
			}
			scalePercentString := scale.Get("value").String()
			scalePercent, _ := strconv.ParseFloat(scalePercentString, 64)
			ix := int(math.Round(float64(img.Bounds().Dx()) * scalePercent))
			iy := int(math.Round(float64(img.Bounds().Dy()) * scalePercent))
			resized := transform.Resize(img, ix-ix%2, iy-iy%3, transform.Linear)
			_lua, finalImage, err := imageToLua(resized, filterName)
			lua = _lua
			size.Get("classList").Call("remove", "invisible")
			size.Set("innerText", fmt.Sprintf("%d kb", len(lua)/1024))
			if err != nil {
				js.Global().Call("alert", err.Error())
				println(err.Error())
			} else {
				buf := new(bytes.Buffer)
				err := png.Encode(buf, finalImage)
				if err != nil {
					js.Global().Call("alert", err.Error())
					return nil
				}
				preview.Set("src", fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buf.Bytes())))
				copy.Get("classList").Call("remove", "invisible")
			}
			return nil
		}))
		return nil
	}))
	<-c
}

func imageToLua(img image.Image, filter string) (string, image.Image, error) {
	if img.Bounds().Dy()%3 != 0 {
		return "", nil, errors.New("image height must be a multiple of 3")
	}

	if img.Bounds().Dx()%2 != 0 {
		return "", nil, errors.New("image width must be a multiple of 2")
	}

	println("Image loaded, quantizing...")
	quant := juroku.Quantize(img, filter)
	println("Image quantized, chunking and generating code...")
	chunked, err := juroku.ChunkImage(quant)
	if err != nil {
		println("Failed to chunk image:", err.Error())
		return "", nil, err
	}
	code, err := juroku.GenerateCode(chunked)
	if err != nil {
		println("Failed to generate code:", err)
		return "", nil, err
	}
	println("Done generating code")
	return string(code), chunked, nil
}
