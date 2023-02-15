package inputs

const (
	INPUT_PLAIN_COLOR     = "PLAIN_COLOR"
	INPUT_LINEAR_GRADIENT = "LINEAR_GRADIENT"
)

/* Represent a color */
type Pattern struct {
	/* The pattern type */
	PatternType string `json:"pattern_type" binding:"required,oneof=PLAIN_COLOR LINEAR_GRADIENT"`

	/* The color list, as hex colors */
	Colors []string `json:"colors" binding:"required,min=1,max=10,dive,hexcolor"`

	/* Rotation for the linear gradient */
	Angle int `json:"angle" binding:"min=0,max=360"`
}
