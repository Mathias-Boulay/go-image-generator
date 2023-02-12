package main

import (
	"image/color"
	"math"
	"src/main/shapes"
	"src/main/utils"

	"github.com/fogleman/gg"
)

/*
func main() {
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

		Scale: 2,

		Pattern: gg.NewSolidPattern(color.RGBA{255, 255, 0, 255}),

		Rotation: 0,

		StrokeWidth: 5,
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

	blob.Draw(dc)

	blob.Scale = 1
	blob.Elevation = 20

	blob.Draw(dc)

	dc.SavePNG("out.png")
}

func DrawCircleBackground(dc *gg.Context, centerColor, outerColor color.Color, layerCount, layerWidth int) {

	for i := 1; i <= layerCount; i++ {
		dc.SetColor(utils.ColorLerp(centerColor, outerColor, (1. / float64(layerCount) * float64(i))))
		dc.MoveTo(250, 0)

		dc.DrawCircle(float64(dc.Width()/2), float64(dc.Height()/2),
			math.Max(float64(dc.Width()-i*layerWidth), float64(dc.Height()-i*layerWidth))/2)
		dc.Fill()
	}
}
