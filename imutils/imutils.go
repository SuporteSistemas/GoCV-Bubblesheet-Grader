package imutils

import (
	"fmt"
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
	return []image.Point{leftMost[0], rightMost[1], rightMost[0], leftMost[1]}
}

//FourPointTransform is the golang versioin of the function
//four_point_transform from the python lib imutils
func FourPointTransform(img gocv.Mat, pts []image.Point, dst *gocv.Mat) {
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

	dt := []image.Point{
		image.Pt(0, 0),
		image.Pt(int(maxWidth)-1, 0),
		image.Pt(int(maxWidth)-1, int(maxHeight)-1),
		image.Pt(0, int(maxHeight)-1)}

	m := gocv.GetPerspectiveTransform(rect, dt)
	gocv.WarpPerspective(img, dst, m, image.Pt(int(maxWidth), int(maxHeight)))

}

//SortContours is the golang versioin of the function
//sort_contours from the python lib imutils
func SortContours(cnts [][]image.Point, method string) [][]image.Point {
	reverse := false
	vertical := false

	if method == "right-to-left" || method == "bottom-to-top" {
		reverse = true
	}

	if method == "bottom-to-top" || method == "top-to-bottom" {
		vertical = true
	}

	var boundingBoxes []image.Rectangle
	for _, c := range cnts {
		boundingBoxes = append(boundingBoxes, gocv.BoundingRect(c))
	}
	for i := range boundingBoxes {
		for j := 0; j < i; j++ {
			if vertical {
				if !reverse {
					if boundingBoxes[i].Min.Y < boundingBoxes[j].Min.Y {
						boundingBoxes[i], boundingBoxes[j] = boundingBoxes[j], boundingBoxes[i]
						cnts[i], cnts[j] = cnts[j], cnts[i]
					}
				} else if reverse {
					if boundingBoxes[i].Min.Y > boundingBoxes[j].Min.Y {
						boundingBoxes[i], boundingBoxes[j] = boundingBoxes[j], boundingBoxes[i]
						cnts[i], cnts[j] = cnts[j], cnts[i]
					}
				}
			} else if !vertical {
				if !reverse {
					if boundingBoxes[i].Min.X < boundingBoxes[j].Min.X {
						boundingBoxes[i], boundingBoxes[j] = boundingBoxes[j], boundingBoxes[i]
						cnts[i], cnts[j] = cnts[j], cnts[i]
					}
				} else if reverse {
					if boundingBoxes[i].Min.X > boundingBoxes[j].Min.X {
						boundingBoxes[i], boundingBoxes[j] = boundingBoxes[j], boundingBoxes[i]
						cnts[i], cnts[j] = cnts[j], cnts[i]
					}
				}
			}
		}
	}
	fmt.Println(boundingBoxes)
	return cnts
}
