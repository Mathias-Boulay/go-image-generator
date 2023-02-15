package inputs

const (
	INPUT_TYPE_BLOB    = "BLOB"
	INPUT_TYPE_CIRCLE  = "CIRCLE"
	INPUT_TYPE_POLYGON = "POLYGON"
)

/* Represent blob parameters */
type Blob struct {
	/* Scale of the blob */
	Scale float64 `json:"scale" binding:"required,min=0,max=10"`

	/* Elevation of the blob */
	Elevation uint32 `json:"elevation" binding:"min=0,max=40"`

	/* If set, will switch to a stroke instead of filling */
	StrokeWidth float64 `json:"stroke_width" binding:"min=0,max=100"`

	/* Rotation of the blob */
	Rotation int `json:"rotation" binding:"min=0,max=360"`

	/* Color of the blob */
	Pattern Pattern `json:"pattern" binding:"required"`

	/* Coordinates, as array to keep the data sent compact */
	Coordinates [][]int `json:"coordinates" binding:"required,dive,len=2"`

	/* Center of the blob */
	Center []float64 `json:"center" binding:"required,len=2"`

	/* The shape type */
	ShapeType string `json:"shape_type" binding:"required,oneof=CIRCLE BLOB POLYGON"`
}
