package imutils

import (
	"image"
	"math"
	"sort"

	"gocv.io/x/gocv"
)

//GrabContours is the golang versioin of the function
//grab_contours from the python lib imutils
func GrabContours(cnts [][]image.Point) []image.Point {

	if len(cnts) == 2 {
		cnts := cnts[0]
		return cnts
	} else if len(cnts) == 3 {
		cnts := cnts[1]
		return cnts
	}
	return nil
}

//OrderPoints is the golang versioin of the function
//order_points from the python lib imutils
func OrderPoints(pts []image.Point) []image.Point {
	sort.Slice(pts[:], func(a, b int) bool {
		return pts[a].X < pts[b].X
	})
	leftMost := pts[:2]
	rightMost := pts[2:]
	sort.Slice(leftMost[:], func(a, b int) bool {
		return leftMost[a].Y < leftMost[b].Y
	})
	sort.Slice(rightMost[:], func(a, b int) bool {
		pythagoreanDist := func(i, j image.Point) float64 {
			var x, y int
			if i.X > j.X {
				x = (i.X - j.X) ^ 2
			} else {
				x = (j.X - i.X) ^ 2
			}
			if i.Y > j.Y {
				y = (i.Y - j.Y) ^ 2
			} else {
				y = (j.Y - i.Y) ^ 2
			}

			return math.Sqrt((float64(x) + float64(y)))
		}
		return pythagoreanDist(leftMost[0], rightMost[0]) <
			pythagoreanDist(leftMost[0], rightMost[1])
	})
	return []image.Point{leftMost[0], rightMost[0], rightMost[1], leftMost[1]}
}

//FourPointTransform is the golang versioin of the function
//four_point_transform from the python lib imutils
func FourPointTransform(img gocv.Mat, pts []image.Point) gocv.Mat {
	rect := OrderPoints(pts)
	tl := rect[0]
	tr := rect[1]
	br := rect[2]
	bl := rect[3]

	widthA := math.Sqrt(math.Pow(float64(br.X)-float64(bl.X), 2) +
		math.Pow(float64(br.Y)-float64(bl.Y), 2))
	widthB := math.Sqrt(math.Pow(float64(tr.X)-float64(tl.X), 2) +
		math.Pow(float64(tr.Y)-float64(tl.Y), 2))
	maxWidth := math.Max(widthA, widthB)

	heightA := math.Sqrt(math.Pow(float64(tr.X)-float64(br.X), 2) +
		math.Pow(float64(tr.Y)-float64(br.Y), 2))
	heightB := math.Sqrt(math.Pow(float64(tl.X)-float64(bl.X), 2) +
		math.Pow(float64(tl.Y)-float64(bl.Y), 2))
	maxHeight := math.Max(heightA, heightB)

	dst := []image.Point{
		image.Pt(0, 0),
		image.Pt(int(maxWidth)-1, 0),
		image.Pt(int(maxWidth)-1, int(maxHeight)-1),
		image.Pt(0, int(maxHeight)-1)}

	m := gocv.GetPerspectiveTransform(rect, dst)
	var returnMat gocv.Mat
	gocv.WarpPerspective(img, &returnMat, m, image.Pt(int(maxWidth), int(maxHeight)))
	return returnMat
}
