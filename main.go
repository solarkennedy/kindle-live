package main

import (
	"log"
	_ "net/http"
	"os"

	"bytes"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"image/png"
	"math"
)

type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	} else {
		return 255
	}
}

func render_image() image.Image {
	w := 1072
	h := 1448
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 40.0
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				cr.Brightness(float64(x), float64(y)),
				cr.Brightness(float64(x), float64(y)),
				cr.Brightness(float64(x), float64(y)),
				255}
			img.Set(x, y, c)
		}
	}
	img.Set(2, 3, color.RGBA{255, 0, 0, 255})
	return img
}

func serve_image(c *gin.Context) {
	buffer := new(bytes.Buffer)

	img := render_image()
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	c.Data(200, "image/png", buffer.Bytes())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})
	router.GET("/image.png", serve_image)
	router.Run(":" + port)
}
