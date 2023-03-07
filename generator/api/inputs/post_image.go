package inputs

/* Struct representing the post API call */
type PostCreateImage struct {
	/* Background of the image */
	Background Pattern `json:"background" binding:"required"`

	/* Width of the generated image */
	Width int `json:"width" binding:"required,min=0,max=1920"`

	/* Height of the generated image */
	Height int `json:"height" binding:"required,min=0,max=1920"`

	/* The list of shapes to draw */
	Blobs []DrawingInstructions `json:"blobs" binding:"dive,min=1,max=10"`
}

type DrawingInstructions struct {
	/* The blob to draw */
	Blob Blob `json:"blob" binding:"required"`

	/* If desired, draw the shape repeatedly */
	Options *SteppedDrawingOptions `json:"options" binding:"omitempty"`
}

type SteppedDrawingOptions struct {
	/* How many steps will be executed */
	Steps int `json:"steps" binding:"required,min=1,max=20"`

	/* The difference in scale between steps */
	ScaleStep float64 `json:"scale_step"`

	/* The difference in rotation between steps */
	RotationStep int `json:"rotation_step"`

	/* The difference in elevation between steps */
	ElevationStep uint32 `json:"elevation_step"`

	/* The difference in coordinates between steps
	* Note that the coordinates interpretation will differ according to Blob settings
	 */
	TranslateStep []float64 `json:"translate_step" binding:"omitempty,len=2,omitempty,dive,min=-1,max=1"`

	/* Start and end patterns, lerping from one to another */
	StartPattern *Pattern `json:"start_pattern" binding:"omitempty"`
	EndPattern   *Pattern `json:"end_pattern"`
}
