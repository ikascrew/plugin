package file_test

import (
	"testing"

	"gocv.io/x/gocv"
)

func TestProbeSample(t *testing.T) {
	cap, err := gocv.VideoCaptureFile("sample.mp4")
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	defer cap.Close()
	t.Logf("frames=%v fps=%v w=%v h=%v", cap.Get(gocv.VideoCaptureFrameCount), cap.Get(gocv.VideoCaptureFPS), cap.Get(gocv.VideoCaptureFrameWidth), cap.Get(gocv.VideoCaptureFrameHeight))
}
