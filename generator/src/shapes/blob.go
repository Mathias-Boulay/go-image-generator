package shapes

import (
	"image/color"
	"math"
	"src/main/utils"

	"github.com/esimov/stackblur-go"
	"github.com/fogleman/gg"
	"github.com/icza/gog"
)

const (
	SHAPE_TYPE_BLOB    = 0
	SHAPE_TYPE_POLYGON = 1
	SHAPE_TYPE_CIRCLE  = 2
)

type Blob struct {
	/* Internal stuff, used by the drawing step */

	/** Current drawing context, is it dangerous to leak ? */
	currentContext *gg.Context
	/** drawing context isolated from the main one, used to draw the shape itself **/
	shapeContext *gg.Context

	/** Bounds, used to remap gradient shaders */
	leftBound, rightBound gg.Point

	/** Represent all the points of the blob,
	the sequence is considered open ended and will be closed automatically
	*/
	points []gg.Point

	/** Represent the shape type. Affects how points are used */
	shapeType int

	/* Exposed stuff, made to be tweaked as one would like */

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

/* Set the shape as a blob */
func (blob *Blob) SetBlob(points []gg.Point) {
	blob.points = points
	blob.shapeType = SHAPE_TYPE_BLOB
}

/* Set the shape as a polygon */
func (blob *Blob) SetPolygon(points []gg.Point) {
	blob.points = points
	blob.shapeType = SHAPE_TYPE_POLYGON
}

/* Set the shape as a circle */
func (blob *Blob) SetCircle(radius float64) {
	blob.points = []gg.Point{
		{X: -radius, Y: -radius},
		{X: radius, Y: radius},
	}
	blob.shapeType = SHAPE_TYPE_CIRCLE
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
	center := gg.Point{X: float64(sc.Width() / 2), Y: float64(sc.Height() / 2)}
	sc.ScaleAbout(blob.Scale, blob.Scale, center.X, center.Y)

	if blob.Elevation != 0 {
		// Draw the elevation shadow first within the context
		sc.SetColor(color.Black)
		blob.fillPath()
		blob.applyDraw(sc)
		originalImage := sc.Image()
		blurryImage, _ := stackblur.Process(originalImage, blob.Elevation)

		// Temporary restore the scale context, else the image will get scaled twice
		sc.ScaleAbout(1/blob.Scale, 1/blob.Scale, center.X, center.Y)
		sc.DrawImage(blurryImage, 0, 0)
		sc.ScaleAbout(blob.Scale, blob.Scale, center.X, center.Y)
	}

	sc.SetFillStyle(blob)
	sc.SetStrokeStyle(blob)
	blob.fillPath()
	blob.applyDraw(sc)

	// Draw the image from the shape context to the main drawing context
	center = gg.Point{X: getScaledWidth(dc, blob.Position), Y: getScaledHeight(dc, blob.Position)}
	dc.RotateAbout(gg.Radians(float64(blob.Rotation)), center.X, center.Y)
	dc.DrawImage(sc.Image(), int(center.X)-sc.Width()/2, int(center.Y)-sc.Height()/2) // Note the offset
	dc.Pop()
}

/** Actually apply the draw, with the proper method */
func (blob *Blob) applyDraw(dc *gg.Context) {

	// Draw a circle, if necessary
	if blob.shapeType == SHAPE_TYPE_CIRCLE {
		center := gg.Point{X: getScaledWidth(dc, blob.Position), Y: getScaledHeight(dc, blob.Position)}
		dc.DrawCircle(center.X, center.Y, blob.points[1].X)
	} else {
		dc.ClosePath() // Avoid this having a visual effect
	}

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

	if blob.shapeType == SHAPE_TYPE_CIRCLE {
		return
	}

	if blob.shapeType == SHAPE_TYPE_POLYGON {
		for _, v := range blob.points {
			sc.LineTo(center.X+v.X, center.Y+v.Y)
		}
		return
	}

	if blob.shapeType == SHAPE_TYPE_BLOB {
		for i := 0; i < len(blob.points); i += 3 {
			var midIndex, endIndex int
			midIndex = gog.If(i+1 < len(blob.points), i+1, 0)
			endIndex = gog.If(i+2 < len(blob.points), i+2, midIndex)

			sc.CubicTo(
				center.X+blob.points[i].X, center.Y+blob.points[i].Y,
				center.X+blob.points[midIndex].X, center.Y+blob.points[midIndex].Y,
				center.X+blob.points[endIndex].X, center.Y+blob.points[endIndex].Y)
		}

		// When points are a multiple of 3, we need to close out of the loop manually
		if len(blob.points)%3 == 0 {
			end := len(blob.points) - 1
			sc.CubicTo(
				center.X+blob.points[end].X,
				center.Y+blob.points[end].Y,
				// Yes, there is a bug about the averages, it is a feature now
				center.X+(blob.points[end].X+blob.points[0].X)/2,
				center.Y+(blob.points[0].Y+blob.points[0].Y)/2,

				center.X+blob.points[0].X,
				center.Y+blob.points[0].Y,
			)
		}
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

	for _, v := range blob.points {
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

	x = int(float64(x) / blob.Scale)
	y = int(float64(y) / blob.Scale)

	mappedX := utils.Map(float64(x), blob.leftBound.X, blob.rightBound.X, 0, 1000)
	mappedY := utils.Map(float64(y), blob.leftBound.Y, blob.rightBound.Y, 0, 1000)
	return blob.Pattern.ColorAt(
		int(mappedX),
		int(mappedY),
	)
}

/* Draw the blob multiple times */
func (blob *Blob) DrawStepped(dc *gg.Context, options *SteppedDrawingOptions) {
	// Replace the pattern implementation
	options.resetStepCount()
	blob.Pattern = options
	for i := 0; i < options.Steps; i++ {
		blob.Draw(dc)

		blob.Scale += options.ScaleStep
		blob.Rotation += options.RotationStep
		blob.Position.X += options.TranslateStep.X
		blob.Position.Y += options.TranslateStep.Y
		blob.Elevation += options.ElevationStep
		options.incStepCount()
	}
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
