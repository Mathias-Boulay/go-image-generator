package main

import (
	//shapes "src/main/api"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/mathias-boulay/generator/api/shapes"
)

/*
func toto() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		fmt.Println("ping")
		file := memfile.New([]byte{})

		dc := gg.NewContext(720, 720)

		DrawSun(dc, 100)

		dc.SetRGB(0, 0, 0)
		dc.Fill()
		dc.EncodePNG(file)
		file.Seek(0, 0)

		c.JSON(200, gin.H{
			"message": "pong",
			"image":   base64.StdEncoding.EncodeToString(file.Bytes()),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}*/
/*
func main() {
	dc := gg.NewContext(500, 500)
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.Fill()
	//outerColor := color.RGBA{255, 255, 255, 255}
	//centerColor := color.RGBA{0, 0, 0, 255}
	//DrawCircleBackground(dc, centerColor, outerColor, 50, 10)

	blob := shapes.Blob{

		Position: gg.Point{X: 0.5, Y: 0.5},

		Scale: 1,

		Pattern: gg.NewSolidPattern(color.RGBA{255, 255, 0, 255}),

		Rotation: 0,

		StrokeWidth: 0,
	}

	blob.SetPolygon([]gg.Point{
		{X: -100, Y: -100},
		{X: 100, Y: -100},
		{X: 100, Y: 100},
		{X: -100, Y: 100},
	})

	// blob.SetBlob([]gg.Point{
	// 	{X: -100, Y: -100},
	// 	{X: 100, Y: -100},
	// 	{X: 100, Y: 100},
	// 	{X: -100, Y: 100},
	// })

	//blob.SetCircle(100)

	grad := gg.NewLinearGradient(utils.ToLinerarGradientCoordinates(0))
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	blob.Pattern = grad
	blob.Elevation = 20
	blob.Scale = 2
	blob.DrawStepped(dc, &shapes.SteppedDrawingOptions{
		Steps:         3,
		ScaleStep:     -0.15,
		TranslateStep: gg.Point{X: 0.02, Y: 0.02},
		RotationStep:  10,
		EndPattern:    grad,
		StartPattern:  gg.NewSolidPattern(color.RGBA{0, 0, 255, 255}),
	})

	dc.SavePNG("out.png")
}*/

func main() {
	// Prepare API system
	engine := gin.Default()
	shapes.RegisterRoute(engine)
	port, found := os.LookupEnv("SERVER_PORT")
	if !found {
		port = "8080"
	}

	engine.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080
}
