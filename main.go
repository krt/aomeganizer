package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/krt/aomeganizer/converter"
	"github.com/minodisk/go-fix-orientation/processor"
)

// 色調調整とか
// 標準入力で受けたときは標準出力で返す
// orientationの補正

func main() {
	var im image.Image
	var canvas *image.RGBA
	var err error
	// var ioFromFile bool
	var fileName string
	var rbytes []byte
	var faceCount int

	if len(os.Args) >= 2 {
		// ioFromFile = true
		fileName = os.Args[1]
		rbytes, err = ioutil.ReadFile(GetFilePathFromWd(fileName))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}
	} else {
		// ioFromFile = false
		rbytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}
	}
	if im, err = Input(rbytes); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	haarcascadeFile := GetFilePath("haarcascade_frontalface_alt.xml")
	canvas, faceCount = converter.Convert(im, haarcascadeFile)
	fmt.Printf("Num of Aomeganes: %v\n", faceCount)
	Output(canvas)
}

// Input ...
func Input(b []byte) (image.Image, error) {
	var im, fixedImage image.Image
	var orientation int
	var err error

	if im, _, err = image.Decode(bytes.NewReader(b)); err != nil {
		return nil, fmt.Errorf("Error while decoding image: %s\n", err)
	}

	orientation, err = processor.ReadOrientation(bytes.NewReader(b))
	if err == nil {
		fixedImage = processor.ApplyOrientation(im, orientation)
	} else {
		fixedImage = im
	}

	return fixedImage, nil
}

// GetFilePath ...
func GetFilePath(fileName string) string {
	_, currentfile, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(currentfile), fileName)
}

// GetFilePathFromWd ...
func GetFilePathFromWd(fileName string) string {
	var curDir, _ = os.Getwd()
	curDir += "/"
	return curDir + fileName
}

// Output ...
func Output(canvas *image.RGBA) {
	out, err := os.Create("./output.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	option := &jpeg.Options{Quality: 80}
	err = jpeg.Encode(out, canvas, option)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
