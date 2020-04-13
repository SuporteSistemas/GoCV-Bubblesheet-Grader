package main

import (
	"image"
	imutils "./imutils"
	"gocv.io/x/gocv"
	"image/color"
	"fmt"
)

func main() {
	answerKey := map[int]int{0: 1, 1: 4, 2: 0, 3: 3, 4: 1}
	rightColor := color.RGBA{R:0,G:255,B:0,A:0}
	wrongColor := color.RGBA{R:255,G:0,B:0,A:0}
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

	imutils.SortContours(cnts, "top-to-bottom")
	var correct int = 0
	var answerSheet map[int]int;
	answerSheet = make(map[int]int)

	for i:=0; i<len(cnts); i+=len(answerKey) {
		imutils.SortContours(cnts[i:i+5], "left-to-right")
		bubbleCount := 0 
		for j:=i; j<i+len(answerKey); j++ {
			mask:= gocv.NewMatWithSize(thresh.Rows(),thresh.Cols(),thresh.Type())
			bubble:= gocv.NewMatWithSize(thresh.Rows(),thresh.Cols(),thresh.Type())
			c := [][]image.Point{cnts[j]}
			gocv.DrawContours(&mask, c, -1, color.RGBA{R:255,G:255,B:255,A:255}, -1)
			gocv.BitwiseAndWithMask(thresh,thresh,&bubble,mask)
			newCount := gocv.CountNonZero(bubble)
			if(newCount > bubbleCount){
				bubbleCount = newCount
				answerSheet[i/5] = j-i
			}
		}
		c := [][]image.Point{cnts[answerSheet[i/len(answerKey)]+i]}
		if(answerSheet[i/len(answerKey)]==answerKey[i/len(answerKey)]) { 
			gocv.DrawContours(&paper, c, 0, rightColor, 2)
			correct++
		} else { gocv.DrawContours(&paper, c, 0, wrongColor, 2)}
	}

	grade := fmt.Sprintf("%.0f%%",float32(correct)/float32(len(answerKey))*100)
	gocv.PutText(&paper, grade, image.Point{X:10,Y:70},gocv.FontHersheyPlain, 5, wrongColor, 3)
	original.IMShow(bubbleSheet)
	exam.IMShow(paper)
	original.WaitKey(0)
}
