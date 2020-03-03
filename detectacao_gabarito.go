package main

import (
	"image"

	imutils "./imutils"

	"gocv.io/x/gocv"
)

func main() {
	original := gocv.NewWindow("ORIGINAL")
	exam := gocv.NewWindow("EXAM")
	bubbleSheet := gocv.IMRead("image.png", 1)
	filter := gocv.NewMat()
	gocv.CvtColor(bubbleSheet, &filter, gocv.ColorBGRToGray)
	gocv.GaussianBlur(filter, &filter, image.Point{X: 5, Y: 5}, 0, 0, 0)
	gocv.Canny(filter, &filter, 75, 200)
	ctr := gocv.FindContours(filter, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var docCnt []image.Point
	if len(ctr) > 0 {
		for _, c := range ctr {
			peri := gocv.ArcLength(c, true)
			approx := gocv.ApproxPolyDP(c, 0.02*peri, true)
			if len(approx) == 4 {
				docCnt = approx
				break
			}
		}
	}
	paper := imutils.FourPointTransform(bubbleSheet, docCnt)

	original.IMShow(bubbleSheet)
	exam.IMShow(paper)
	original.WaitKey(0)
}
