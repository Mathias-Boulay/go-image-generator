package shapes

import (
	"image/color"
	"math"
	"src/main/utils"

	"github.com/esimov/stackblur-go"
	"github.com/fogleman/gg"
	"github.com/icza/gog"
)

type Blob struct {
	/* Internal stuff, used by the drawing step */

	/** Current painting context, is it dangerous to leak ? */
	currentContext *gg.Context

	/** Bounds, used to remap gradient shaders */
	leftBound, rightBound gg.Point

	/* Exposed stuff, made to be tweaked as one would like */

	/**
	Represent all the Points of the blob,
	the sequence is considered open ended and will be closed automatically
	*/
	Points []gg.Point

	/** Position withing the drawing context,
	between 0 and 1 on each axis to accomodate for various formats */
	Position gg.Point

	/** Scale of the blob, a higher than 1 scale will produce a bigger blob */
	Scale float64

	/** Pattern to fill the blob with, can be a gradient */
	Pattern gg.Pattern

	/** Clockwise rotation in degrees  */
	Rotation int
}

func (blob *Blob) Draw(dc *gg.Context) {
	blob.currentContext = dc
	blob.findBounds()

	dc.Push()
	center := gg.Point{X: getScaledWidth(dc, blob.Position), Y: getScaledHeight(dc, blob.Position)}

	dc.ScaleAbout(blob.Scale, blob.Scale, center.X, center.Y)
	dc.RotateAbout(gg.Radians(float64(blob.Rotation)), center.X, center.Y)

	dc.SetFillStyle(blob)
	dc.SetStrokeStyle(blob)

	dc.SetLineWidth(5)

	for i := 0; i < len(blob.Points); i += 3 {
		// Fill all paths, taking into account to wrap around the list
		var midIndex, endIndex int
		midIndex = gog.If(i+1 < len(blob.Points), i+1, 0)
		endIndex = gog.If(i+2 < len(blob.Points), i+2, midIndex)

		dc.CubicTo(
			center.X+blob.Points[i].X, center.Y+blob.Points[i].Y,
			center.X+blob.Points[midIndex].X, center.Y+blob.Points[midIndex].Y,
			center.X+blob.Points[endIndex].X, center.Y+blob.Points[endIndex].Y)
	}

	// When points are a multiple of 3, we need to close out of the loop
	if len(blob.Points)%3 == 0 {
		end := len(blob.Points) - 1
		dc.CubicTo(
			center.X+blob.Points[end].X,
			center.Y+blob.Points[end].Y,
			// Yes, there is a bug about the averages, it is a feature now
			center.X+(blob.Points[end].X+blob.Points[0].X)/2,
			center.Y+(blob.Points[0].Y+blob.Points[0].Y)/2,

			center.X+blob.Points[0].X,
			center.Y+blob.Points[0].Y,
		)
	}

	dc.ClosePath() // Should not be used
	dc.Stroke()

	blurred, _ := stackblur.Process(dc.Image(), 2)
	dc.Pop()

	dc.DrawImage(blurred, 0, 0)

}

/*
Compute the top left and bottom right bounds from current points
set internal fields
*/
func (blob *Blob) findBounds() {
	left := math.Inf(1)
	top := math.Inf(1)
	right := math.Inf(-1)
	bottom := math.Inf(-1)

	for _, v := range blob.Points {
		left = math.Min(left, v.X)
		top = math.Min(top, v.Y)
		right = math.Max(right, v.X)
		bottom = math.Max(bottom, v.Y)
	}
	blob.leftBound = gg.Point{X: left, Y: top}
	blob.rightBound = gg.Point{X: right, Y: bottom}
}

/*
Translation layer for gradients used by Blob
The translation maps x and y coordinates on a 0-1000 space (since they are integers)
and passes it down to the real pattern implementation.
*/
func (blob *Blob) ColorAt(x, y int) color.Color {
	center := blob.getScaledPosition(blob.currentContext)

	x -= int(center.X)
	y -= int(center.Y)

	x /= int(blob.Scale)
	y /= int(blob.Scale)

	mappedX := utils.Map(float64(x), blob.leftBound.X, blob.rightBound.X, 0, 1000)
	mappedY := utils.Map(float64(y), blob.leftBound.Y, blob.rightBound.Y, 0, 1000)
	return blob.Pattern.ColorAt(
		int(mappedX),
		int(mappedY),
	)
}

/* Util functions for scaling back the [0,1] coords */
func (blob *Blob) getScaledPosition(dc *gg.Context) gg.Point {
	return gg.Point{X: getScaledWidth(dc, blob.Position), Y: getScaledHeight(dc, blob.Position)}
}

func getScaledWidth(dc *gg.Context, position gg.Point) float64 {
	return float64(dc.Width()) * position.X
}

func getScaledHeight(dc *gg.Context, position gg.Point) float64 {
	return float64(dc.Height()) * position.Y
}
