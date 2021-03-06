package main

import (
	"flag"
	"log"
	"os"

	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	_ "github.com/schachmat/wego/backends"
	_ "github.com/schachmat/wego/frontends"
	"github.com/schachmat/wego/iface"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"time"
)

var weather iface.Backend

func addLabel(img *image.Gray, x, y int, size float64, label string) {
	black := image.NewUniform(color.Black)
	b, err := ioutil.ReadFile("./ubuntu.ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: size,
		DPI:  600,
	})
	d := &font.Drawer{
		Dst:  img,
		Src:  black,
		Face: face,
	}
	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d.DrawString(label)
	defer face.Close()
}

func addWeatherIcon(img *image.Gray, x, y int, size float64, label string) {
	black := image.NewUniform(color.Black)
	b, err := ioutil.ReadFile("./weathericons-regular-webfont.ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: size,
		DPI:  600,
	})
	d := &font.Drawer{
		Dst:  img,
		Src:  black,
		Face: face,
	}
	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d.DrawString(label)
	defer face.Close()
}

func fetchForecast(location string) iface.Data {
	numdays := 3
	fmt.Println("Going to fetch the weather....")
	r := weather.Fetch(location, numdays)
	return r
}

func weatherBackendSetup() {
	for _, be := range iface.AllBackends {
		be.Setup()
	}
	for _, fe := range iface.AllFrontends {
		fe.Setup()
	}
	api_key := os.Getenv("forecast_api_key")
	flag.Set("forecast-api-key", api_key)
	flag.Set("forecast-lang", "en")
	flag.Set("forecast-debug", "false")
	flag.Set("frontend", "emoji")
	flag.Parse()
	ok := false
	weather, ok = iface.AllBackends["forecast.io"]
	if !ok {
		log.Fatalf("Could not find selected backend forecast.io")
	}
}

func codeToIcon(code iface.WeatherCode, day bool) string {
	if day {
		codes := map[iface.WeatherCode]string{
			iface.CodeUnknown:             "\uf07b",
			iface.CodeCloudy:              "\uf002",
			iface.CodeFog:                 "\uf003",
			iface.CodeHeavyRain:           "\uf008",
			iface.CodeHeavyShowers:        "\uf009",
			iface.CodeHeavySnow:           "\uf00a",
			iface.CodeHeavySnowShowers:    "\uf06b",
			iface.CodeLightRain:           "\uf006",
			iface.CodeLightShowers:        "\uf006",
			iface.CodeLightSleet:          "\uf0b2",
			iface.CodeLightSleetShowers:   "\uf068",
			iface.CodeLightSnow:           "\uf00a",
			iface.CodeLightSnowShowers:    "\uf06b",
			iface.CodePartlyCloudy:        "\uf07d",
			iface.CodeSunny:               "\uf00d",
			iface.CodeThunderyHeavyRain:   "\uf00e",
			iface.CodeThunderyShowers:     "\uf010",
			iface.CodeThunderySnowShowers: "\uf06b",
			iface.CodeVeryCloudy:          "\uf013",
		}
		return codes[code]
	} else {
		codes := map[iface.WeatherCode]string{
			iface.CodeUnknown:             "\uf07b",
			iface.CodeCloudy:              "\uf086",
			iface.CodeFog:                 "\uf04a",
			iface.CodeHeavyRain:           "\uf028",
			iface.CodeHeavyShowers:        "\uf028",
			iface.CodeHeavySnow:           "\uf02a",
			iface.CodeHeavySnowShowers:    "\uf06d",
			iface.CodeLightRain:           "\uf026",
			iface.CodeLightShowers:        "\uf029",
			iface.CodeLightSleet:          "\uf0b4",
			iface.CodeLightSleetShowers:   "\uf06a",
			iface.CodeLightSnow:           "\uf02a",
			iface.CodeLightSnowShowers:    "\uf06d",
			iface.CodePartlyCloudy:        "\uf086",
			iface.CodeSunny:               "\uf02e",
			iface.CodeThunderyHeavyRain:   "\uf02d",
			iface.CodeThunderyShowers:     "\uf02d",
			iface.CodeThunderySnowShowers: "\uf06d",
			iface.CodeVeryCloudy:          "\uf013",
		}
		return codes[code]
	}
}

func hourToBool(h int) bool {
	return h > 6 && h < 18
}

func drawForecast(img *image.Gray, y int, day iface.Day) {
	date_string := fmt.Sprintf(day.Date.Format("Monday, 2006-01-02 (MST)"))
	addLabel(img, 50, y, 4, fmt.Sprintf("Forcast for %s", date_string))
	for i := 0; i < 24; i++ {
		x := i*40 + 20
		addLabel(img, x, y+10, 1, fmt.Sprintf("Hour %d", i))
	}
	for _, slot := range day.Slots {
		slot_hour := slot.Time.Hour()
		x := slot_hour*40 + 20
		addLabel(img, x, y+20, 1, fmt.Sprintf("Slot %d", slot_hour))
		addWeatherIcon(img, x, y+45, 3, codeToIcon(slot.Code, hourToBool(slot_hour)))

		var chance int = 0
		if slot.ChanceOfRainPercent != nil {
			chance = *slot.ChanceOfRainPercent
		}
		addWeatherIcon(img, x, y+60, 1, "\uf04e")
		addLabel(img, x+10, y+60, 1, fmt.Sprintf("%d%%", chance))

		addWeatherIcon(img, x, y+70, 1, "\uf055")
		addLabel(img, x+10, y+70, 1, fmt.Sprintf("%.1f °C", *slot.FeelsLikeC))
	}
}

func speedToBeaufort(kmph float32) (string, string) {
	if kmph < 1 {
		return "\uf0b7", "Calm"
	} else if kmph < 5 {
		return "\uf0b8", "Light air"
	} else if kmph < 11 {
		return "\uf0b9", "Light breeze"
	} else if kmph < 19 {
		return "\uf0ba", "Gentle breeze"
	} else if kmph < 28 {
		return "\uf0bb", "Moderate breeze" //4
	} else if kmph < 38 {
		return "\uf0bc", "Fresh breeze"
	} else if kmph < 49 {
		return "\uf0bd", "Strong breeze"
	} else if kmph < 61 {
		return "\uf0be", "High winds" //7
	} else if kmph < 74 {
		return "\uf0bf", "Fresh gale"
	} else if kmph < 88 {
		return "\uf0c0", "Strong gale"
	} else if kmph < 102 {
		return "\uf0c1", "Storm" //10
	} else if kmph < 117 {
		return "\uf0c2", "Violent storm"
	} else {
		return "\uf0c3", "Hurricane force"
	}
}

func directionToIcon(direction int) (string, string) {
	if direction < 22 {
		return "\uf05c", "Southerly"
	} else if direction < 67 {
		return "\uf05a", "Southwesterly"
	} else if direction < 112 {
		return "\uf059", "Westerly"
	} else if direction < 157 {
		return "\uf05d", "Northwesterly"
	} else if direction < 202 {
		return "\uf060", "Northerly"
	} else if direction < 247 {
		return "\uf053", "Northeasterly"
	} else if direction < 292 {
		return "\uf061", "Easterly"
	} else if direction < 337 {
		return "\uf05b", "Southeasterly"
	} else {
		return "\uf05c", "Southerly"
	}
}

func renderForecast(img *image.Gray, r iface.Data) {
	c := r.Current
	addLabel(img, 50, 50, 6, "Current Weather")

	hour := c.Time.Hour()
	addWeatherIcon(img, 70, 230, 18, codeToIcon(c.Code, hourToBool(hour)))
	addLabel(img, 100, 320, 4, c.Desc)

	addWeatherIcon(img, 368, 130, 4, "\uf055")
	addLabel(img, 400, 130, 4, fmt.Sprintf("Temperature: %.1f °C (Feels like %.1f °C)", *c.TempC, *c.FeelsLikeC))

	addWeatherIcon(img, 365, 170, 4, "\uf07a")
	addLabel(img, 400, 170, 4, fmt.Sprintf("Humidity: %d%%", *r.Current.Humidity))

	addWeatherIcon(img, 370, 210, 4, "\uf04e")
	var chance int = 0
	if c.ChanceOfRainPercent != nil {
		chance = *c.ChanceOfRainPercent
	}
	addLabel(img, 400, 210, 4, fmt.Sprintf("Rain chance: %d%%", chance))

	wind_icon, wind_desc := speedToBeaufort(*c.WindspeedKmph)
	addWeatherIcon(img, 360, 250, 4, wind_icon)
	addLabel(img, 400, 250, 4, fmt.Sprintf("Windspeed: %.1f km/h (%s)", *c.WindspeedKmph, wind_desc))

	windd_icon, windd_desc := directionToIcon(*c.WinddirDegree)
	addWeatherIcon(img, 370, 290, 4, windd_icon)
	addLabel(img, 400, 290, 4, fmt.Sprintf("Wind direction: %d° (%s)", *c.WinddirDegree, windd_desc))

	for i, d := range r.Forecast {
		y := i*200 + 400
		drawForecast(img, y, d)
	}
}

func render_image(location string, ip string) image.Image {
	w := 1072
	h := 1448
	img := image.NewGray(image.Rect(0, 0, w, h))

	for x := 1; x < w-1; x++ {
		for y := 1; y < h-1; y++ {
			c := color.RGBA{255, 255, 255, 255}
			img.Set(x, y, c)
		}
	}

	label_string := fmt.Sprintf("Generated on: %s", time.Now().String())
	addLabel(img, 50, 1420, 2, label_string)

	addLabel(img, 50, 1400, 2, fmt.Sprintf("Client IP: %s (%s)", ip, location))

	forecast := fetchForecast(location)
	renderForecast(img, forecast)

	return img
}

func serve_image(c *gin.Context) {
	buffer := new(bytes.Buffer)
	location := c.DefaultQuery("location", "37.676878,-122.459695")
	ip := c.ClientIP()
	img := render_image(location, ip)
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

	weatherBackendSetup()

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
