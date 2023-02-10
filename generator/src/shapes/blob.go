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

	/** Current drawing context, is it dangerous to leak ? */
	currentContext *gg.Context
	/** drawing context isolated from the main one, used to draw the shape itself **/
	shapeContext *gg.Context

	/** Bounds, used to remap gradient shaders */
	leftBound, rightBound gg.Point

	/* Exposed stuff, made to be tweaked as one would like */

	/** Represent all the Points of the blob,
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

	/** Stroke width, if 0 the shape is filled instead */
	StrokeWidth float64

	/** Clockwise rotation in degrees  */
	Rotation int

	/** Elevation, unitless for now */
	Elevation uint32
}

/* Draw the shape within the main drawing context */
func (blob *Blob) Draw(dc *gg.Context) {
	// Setup context for this draw
	dc.Push()
	blob.currentContext = dc

	// Create and setup the shape context, whitin which the shape is drawn
	blob.findBounds()
	blob.createShapeContext()
	sc := blob.shapeContext

	if blob.Elevation != 0 {
		// Draw the elevation shadow first within the context
		sc.SetColor(color.Black)
		blob.fillPath()
		blob.applyDraw(sc)
		originalImage := sc.Image()
		blurryImage, _ := stackblur.Process(originalImage, blob.Elevation)

		sc.DrawImage(blurryImage, 0, 0)
	}

	sc.SetFillStyle(blob)
	sc.SetStrokeStyle(blob)
	blob.fillPath()
	blob.applyDraw(sc)

	// Draw the image from the shape context to the main drawing context
	center := gg.Point{X: getScaledWidth(dc, blob.Position), Y: getScaledHeight(dc, blob.Position)}
	dc.RotateAbout(gg.Radians(float64(blob.Rotation)), center.X, center.Y)
	dc.DrawImage(sc.Image(), int(center.X)-sc.Width()/2, int(center.Y)-sc.Height()/2) // Note the offset
	dc.Pop()
}

/** Actually apply the draw, with the proper method */
func (blob *Blob) applyDraw(dc *gg.Context) {
	dc.ClosePath() // Avoid this having a visual effect
	if blob.StrokeWidth != 0 {
		dc.Stroke()
	} else {
		dc.Fill()
	}
}

/** Fill all paths on the current shape context, taking into account to wrap around the list */
func (blob *Blob) fillPath() {
	sc := blob.shapeContext
	center := gg.Point{X: float64(sc.Width() / 2), Y: float64(sc.Height() / 2)}

	for i := 0; i < len(blob.Points); i += 3 {
		var midIndex, endIndex int
		midIndex = gog.If(i+1 < len(blob.Points), i+1, 0)
		endIndex = gog.If(i+2 < len(blob.Points), i+2, midIndex)

		sc.CubicTo(
			center.X+blob.Points[i].X, center.Y+blob.Points[i].Y,
			center.X+blob.Points[midIndex].X, center.Y+blob.Points[midIndex].Y,
			center.X+blob.Points[endIndex].X, center.Y+blob.Points[endIndex].Y)
	}

	// When points are a multiple of 3, we need to close out of the loop manually
	if len(blob.Points)%3 == 0 {
		end := len(blob.Points) - 1
		sc.CubicTo(
			center.X+blob.Points[end].X,
			center.Y+blob.Points[end].Y,
			// Yes, there is a bug about the averages, it is a feature now
			center.X+(blob.Points[end].X+blob.Points[0].X)/2,
			center.Y+(blob.Points[0].Y+blob.Points[0].Y)/2,

			center.X+blob.Points[0].X,
			center.Y+blob.Points[0].Y,
		)
	}
}

/** Given the current blob settings, create a perfectly sized drawing context */
func (blob *Blob) createShapeContext() {
	width := (blob.rightBound.X - blob.leftBound.X + blob.StrokeWidth) * blob.Scale
	height := (blob.rightBound.Y - blob.leftBound.Y + blob.StrokeWidth) * blob.Scale

	// Elevation Radius is uncoupled from actual scale
	width += float64(blob.Elevation) * 2
	height += float64(blob.Elevation) * 2

	sc := gg.NewContext(int(width), int(height))
	blob.shapeContext = sc

	// Then setup it with the appropriate settings
	center := gg.Point{X: float64(sc.Width() / 2), Y: float64(sc.Height() / 2)}
	sc.ScaleAbout(blob.Scale, blob.Scale, center.X, center.Y)
	sc.SetLineWidth(blob.StrokeWidth)
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
Implement the Pattern interface ! Translation layer for gradients.
The translation maps x and y coordinates on a 0-1000 space (since they are integers)
and passes it down to the real pattern implementation.
*/
func (blob *Blob) ColorAt(x, y int) color.Color {
	center := gg.Point{
		X: float64(blob.shapeContext.Width() / 2),
		Y: float64(blob.shapeContext.Height() / 2),
	}

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
