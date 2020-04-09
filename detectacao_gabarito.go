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
	gray := gocv.NewMat()
	blurred := gocv.NewMat()
	edged := gocv.NewMat()
	paper := gocv.NewMat()
	warped := gocv.NewMat()
	thresh := gocv.NewMat()

	gocv.CvtColor(bubbleSheet, &gray, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &blurred, image.Point{X: 5, Y: 5}, 0, 0, 0)
	gocv.Canny(blurred, &edged, 75, 80)
	ctr := gocv.FindContours(edged, gocv.RetrievalExternal, gocv.ChainApproxSimple)

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
	imutils.FourPointTransform(bubbleSheet, docCnt, &paper)
	imutils.FourPointTransform(gray, docCnt, &warped)

	gocv.Threshold(warped, &thresh, 0, 255, gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)
	ctr = gocv.FindContours(thresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var cnts [][]image.Point
	for _, c := range ctr {
		rect := gocv.BoundingRect(c)
		w := rect.Max.X - rect.Min.X
		h := rect.Max.Y - rect.Min.Y
		ar := float32(w / h)
		if w >= 20 && h >= 20 && ar >= 0.9 && ar <= 1.1 {
			cnts = append(cnts, c)
		}
	}

	cnts = imutils.SortContours(cnts, "top-to-bottom")
	questionCnts = cnts[0]

	original.IMShow(paper)
	exam.IMShow(warped)
	original.WaitKey(0)
}
