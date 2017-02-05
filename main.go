package main

import (
	"log"
	"os"

	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	//	"golang.org/x/image/math/fixed"

	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"time"
)

type Circle struct {
	X, Y, R float64
}

func addLabel(img *image.Gray, x, y int, size float64, label string) {
	col := color.Black
	b, err := loadFontFile()
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: size,
	})
	defer face.Close()
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
	}
	//	d.Dot = fixed.Point26_6{
	//		X: (fixed.I(img) - d.MeasureString(label)) / 2,
	//		Y: fixed.I(y),
	//	}

	//d.DrawString(label)

}

func loadFontFile() ([]byte, error) {
	return ioutil.ReadFile("./weathericons-regular-webfont.ttf")
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 255
	} else {
		return 0
	}
}

func render_image(ip string) image.Image {
	w := 1072
	h := 1448
	img := image.NewGray(image.Rect(0, 0, w, h))

	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 40.0
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				cr.Brightness(float64(x), float64(y)),
				cr.Brightness(float64(x), float64(y)),
				cr.Brightness(float64(x), float64(y)),
				0}
			img.Set(x, y, c)
		}
	}
	label_string := fmt.Sprintf("Generated on: %s", time.Now().String())
	addLabel(img, 100, 100, 26, label_string)
	addLabel(img, 100, 120, 26, fmt.Sprintf("Client IP: %s", ip))
	return img
}

func serve_image(c *gin.Context) {
	buffer := new(bytes.Buffer)
	img := render_image(c.ClientIP())
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
