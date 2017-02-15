package main

import (
	"github.com/lazywei/go-opencv/opencv"
	"testing"
)

func BenchmarkCapture(b *testing.B) {
	cam := opencv.NewCameraCapture(1)
	for i := 0; i < b.N; i++ {
		capture(cam)
	}
}

func BenchmarkCompare(b *testing.B) {
	cam := opencv.NewCameraCapture(1)
	img1, _ := capture(cam)
	img2, _ := capture(cam)

	for i := 0; i < b.N; i++ {
		compare(img1, img2)
	}
}

func BenchmarkShado(b *testing.B) {
	cam := opencv.NewCameraCapture(1)
	for i := 0; i < b.N; i++ {
		shado(cam)
	}
}
