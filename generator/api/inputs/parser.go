package inputs

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/icza/gox/imagex/colorx"
	"github.com/mathias-boulay/generator/shapes"
	"github.com/mathias-boulay/generator/utils"
)

/* Create the pattern from an api color input */
func FromColorInput(input *Pattern) gg.Pattern {
	if input.PatternType == INPUT_LINEAR_GRADIENT {
		gradient := gg.NewLinearGradient(utils.ToLinerarGradientCoordinates(input.Angle))
		for k, v := range input.Colors {
			color, _ := colorx.ParseHexColor(v)

			gradient.AddColorStop(float64(k)/math.Max(float64(len(input.Colors))-1, 1), color)
		}

		return gradient
	}

	if input.PatternType == INPUT_PLAIN_COLOR {
		color, _ := colorx.ParseHexColor(input.Colors[0])
		return gg.NewSolidPattern(color)
	}

	// Shouldn't happen in practice, but watch out
	return nil
}

func FromPointInput(input [][]int) []gg.Point {
	points := make([]gg.Point, 0)
	for _, v := range input {
		points = append(points, gg.Point{X: float64(v[0]), Y: float64(v[1])})
	}

	return points
}

/* Create a blob from the api input */
func FromBlobInput(input *Blob) shapes.Blob {
	blob := shapes.Blob{}

	blob.Position = gg.Point{X: input.Center[0], Y: input.Center[1]}
	blob.Elevation = input.Elevation
	blob.Pattern = FromColorInput(&input.Pattern)
	blob.Scale = input.Scale
	blob.Rotation = input.Rotation
	blob.StrokeWidth = input.StrokeWidth

	if input.ShapeType == INPUT_TYPE_BLOB {
		blob.SetBlob(FromPointInput(input.Coordinates))
	} else if input.ShapeType == INPUT_TYPE_POLYGON {
		blob.SetPolygon(FromPointInput(input.Coordinates))
	} else if input.ShapeType == INPUT_TYPE_CIRCLE {
		blob.SetCircle(math.Abs(float64(input.Coordinates[0][0])))
	}

	return blob
}
