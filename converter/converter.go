package converter

import (
	"bytes"

	"image"
	"image/draw"
	"image/png"

	"github.com/krt/aomeganizer/bindata"
	"github.com/lazywei/go-opencv/opencv"
	"github.com/nfnt/resize"
)

// Convert ...
func Convert(im image.Image, haarcascadeFile string) (*image.RGBA, int) {
	canvas := image.NewRGBA(im.Bounds())
	draw.Draw(canvas, im.Bounds(), im, image.ZP, draw.Src)
	faces := DetectFace(im, haarcascadeFile)

	for _, value := range faces {
		// fmt.Println(value.X())
		// fmt.Println(value.Y())
		faceBound := image.Rect(value.X(), value.Y(), value.X()+value.Width(), value.Y()+value.Height())
		resizedFace := resize.Resize(uint(value.Width()), uint(value.Height()), GetGoodMask(), resize.Lanczos3)
		draw.Draw(canvas, faceBound, resizedFace, image.ZP, draw.Over)
	}

	return canvas, len(faces)
}

// DetectFace ...
func DetectFace(im image.Image, haarCascadeFile string) []*opencv.Rect {
	cvImage := opencv.FromImage(im)
	cascade := opencv.LoadHaarClassifierCascade(haarCascadeFile)
	faces := cascade.DetectObjects(cvImage)
	return faces
}

//GetGoodMask ...
func GetGoodMask() image.Image {
	faceData, _ := bindata.Asset("data/aomeganize_area.png")
	goodMask, _ := png.Decode(bytes.NewReader(faceData))
	return goodMask
}
