package shapes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"src/main/api/inputs"
	"src/main/shapes"

	"github.com/dsnet/golib/memfile"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(gin *gin.Engine) {
	gin.POST("/image", createImage)
}

/*
 *
 */
func createImage(gc *gin.Context) {
	fmt.Println(gc.Request.Body)

	body := inputs.PostCreateImage{}
	if err := gc.ShouldBindJSON(&body); err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": err.Error()})
		return
	}

	// Then create the image
	file := memfile.New([]byte{})
	dc := gg.NewContext(body.Width, body.Height)

	// Draw the background, if any
	backgroundBlob := shapes.Blob{
		Position: gg.Point{X: 0.5, Y: 0.5},
		Pattern:  inputs.FromColorInput(&body.Background),
		Scale:    1,
	}

	backgroundBlob.SetPolygon([]gg.Point{
		{X: -float64(body.Width) / 2, Y: -float64(body.Height) / 2},
		{X: float64(body.Width) / 2, Y: -float64(body.Height) / 2},
		{X: float64(body.Width) / 2, Y: float64(body.Height) / 2},
		{X: -float64(body.Width) / 2, Y: float64(body.Height) / 2},
	})

	backgroundBlob.Draw(dc)

	// Draw blobs in order
	for _, drawingInstructions := range body.Blobs {
		blob := inputs.FromBlobInput(&drawingInstructions.Blob)
		if drawingInstructions.Options == nil {
			blob.Draw(dc)
		} else {
			drawOp := drawingInstructions.Options
			options := shapes.SteppedDrawingOptions{
				Steps:         drawOp.Steps,
				ScaleStep:     drawOp.ScaleStep,
				RotationStep:  drawOp.RotationStep,
				ElevationStep: drawOp.ElevationStep,
				TranslateStep: gg.Point{X: float64(drawOp.TranslateStep[0]), Y: float64(drawOp.TranslateStep[1])},
				StartPattern:  inputs.FromColorInput(drawOp.StartPattern),
				EndPattern:    inputs.FromColorInput(drawOp.EndPattern),
			}

			blob.DrawStepped(dc, &options)
		}
	}

	dc.EncodePNG(file)

	gc.JSON(200, gin.H{
		"image": base64.StdEncoding.EncodeToString(file.Bytes()),
	})
}
