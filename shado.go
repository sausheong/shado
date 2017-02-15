package main

import (
	// "bytes"
	// "compress/flate"
	"errors"
	"fmt"
	"github.com/lazywei/go-opencv/opencv"
	"image"
	// "image/color"
	"image/jpeg"
	// "io"
	// "math"
	"os"
	"strconv"
	"time"
)

func main() {
	cam := opencv.NewCameraCapture(1)
	if cam == nil {
		panic("cannot open camera")
	}
	defer cam.Release()

	for {
		time.Sleep(1000 * time.Millisecond)
		t := time.Now().UnixNano()
		diff := shado(cam, strconv.FormatInt(t, 10))
		fmt.Println(".." + strconv.FormatInt(t, 10) + " - " + strconv.FormatFloat(diff, 'f', 3, 64))
	}

}

//
func shado(cam *opencv.Capture, name string) (diff float64) {

	img1, err := capture(cam, 200, 150)
	if err != nil {
		fmt.Println("Cannot capture image1 from camera:", err)
	}
	time.Sleep(200 * time.Millisecond)
	img2, err := capture(cam, 200, 150)
	if err != nil {
		fmt.Println("Cannot capture image 2from camera:", err)
	}
	diff = difference(img1, img2)
	if diff > 1.5 {
		fmt.Println("** DETECTED **")
		file, err := os.Create("./img" + name + ".jpg")
		if err != nil {
			fmt.Println("Cannot create image file:", err)
		}
		jpeg.Encode(file, img2, &jpeg.Options{Quality: 10})
	}
	return
}

// Capture image from camera, resize it using cubic interpolation
func capture(cam *opencv.Capture, w int, h int) (img image.Image, err error) {
	if cam.GrabFrame() {
		iplimg := cam.RetrieveFrame(1)
		resized := opencv.Resize(iplimg, w, h, opencv.CV_INTER_CUBIC)
		if iplimg != nil {
			img = resized.ToImage()
		} else {
			err = errors.New("Cannot retrieve image from camera")
		}
	}
	return
}

func difference(img1 image.Image, img2 image.Image) float64 {
	b := img1.Bounds()
	var sum int64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			if r1 > r2 {
				sum += int64(r1 - r2)
			} else {
				sum += int64(r2 - r1)
			}
			if g1 > g2 {
				sum += int64(g1 - g2)
			} else {
				sum += int64(g2 - g1)
			}
			if b1 > b2 {
				sum += int64(b1 - b2)
			} else {
				sum += int64(b2 - b1)
			}
		}
	}
	nPixels := (b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y)
	return float64(sum*100) / (float64(nPixels) * 0xffff * 3)
}

// Convert image to grayscale
func gray(src image.Image) image.Image {
	bounds := src.Bounds()
	gray := image.NewGray(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			gray.Set(x, y, src.At(x, y))
		}
	}
	return gray.SubImage(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))
}
