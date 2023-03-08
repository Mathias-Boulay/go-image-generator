package shapes

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/dsnet/golib/memfile"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/mathias-boulay/generator/api/inputs"
	"github.com/mathias-boulay/generator/shapes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterRoute(gin *gin.Engine) {
	gin.POST("/image", createImage)
}

func logInput(gc *gin.Context) {
	if os.Getenv("MONGO_URL") == "" {
		return
	}

	// Log the body content to mongo db
	// Prepare mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Close connection when shutting down the server
	fmt.Println(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		panic(err)
	}

	// Store the body in a buffer to send
	buf := new(bytes.Buffer)
	buf.ReadFrom(gc.Request.Body)

	collection := client.Database("logs").Collection("inputs")
	collection.InsertOne(ctx, bson.M{
		"input": buf.String(),
		"date":  time.Now().Unix(),
	})
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

	logInput(gc)

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
				TranslateStep: gg.Point{X: 0, Y: 0},
				StartPattern:  inputs.FromColorInput(drawOp.StartPattern),
				EndPattern:    inputs.FromColorInput(drawOp.EndPattern),
			}

			if len(drawOp.TranslateStep) > 0 {
				options.TranslateStep = gg.Point{X: float64(drawOp.TranslateStep[0]), Y: float64(drawOp.TranslateStep[1])}
			}

			blob.DrawStepped(dc, &options)
		}
	}

	dc.EncodePNG(file)

	gc.JSON(200, gin.H{
		"image": base64.StdEncoding.EncodeToString(file.Bytes()),
	})
}
