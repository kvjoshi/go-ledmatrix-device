package main

import (
	"flag"
	"fmt"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/canvas"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/matrix"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	rows                   = flag.Int("led-rows", 64, "number of rows supported")
	cols                   = flag.Int("led-cols", 64, "number of columns supported")
	parallel               = flag.Int("led-parallel", 2, "number of daisy-chained panels")
	chain                  = flag.Int("led-chain", 7, "number of displays daisy-chained")
	brightness             = flag.Int("brightness", 100, "brightness (0-100)")
	gpio_slowdown          = flag.Int("led-gpio-slowdown", 3, "GPIO SLOWDOWN")
	pwm_lsb                = flag.Int("led-pwm-lsb-nanoseconds", 80, "lsb nanosec")
	pwm_bits               = flag.Int("led-pwm-bits", 11, "pwm bits")
	hardwareMapping        = flag.String("led-gpio-mapping", "regular", "Name of GPIO mapping used.")
	showRefresh            = flag.Bool("led-show-refresh", true, "Show refresh rate.")
	inverseColors          = flag.Bool("led-inverse", false, "Switch if your matrix has inverse colors on.")
	disableHardwarePulsing = flag.Bool("led-no-hardware-pulse", false, "Don't use hardware pin-pulse generation.")
	pixelMapping           = flag.String("led-pixel-mapper", "U-mapper", "Pixel mapping from api")
	img                    = flag.String("image", "/home/dietpi/cc/i2.jpg", "image path")

	rotate = flag.Int("rotate", 0, "rotate angle, 90, 180, 270")
)

var (
	fileName    string
	fullURLFile string
)

func main() {
	//f, err := os.Open(*img)
	//fatal(err)
	//img, _, err := image.Decode(f)
	config := &matrix.DefaultConfig
	config.Rows = *rows
	config.Cols = *cols
	config.SlowdownGPIO = *gpio_slowdown
	config.PWMBits = *pwm_bits
	config.PWMLSBNanoseconds = *pwm_lsb
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness
	config.HardwareMapping = *hardwareMapping
	config.ShowRefreshRate = *showRefresh
	config.InverseColors = *inverseColors
	config.DisableHardwarePulsing = *disableHardwarePulsing
	//	config.PixelMapping = *pixelMapping

	m, err := matrix.NewRGBLedMatrix(config)
	fatal(err)

	tk := canvas.NewToolKit(m)
	defer tk.Close()

	fullURLFile = "http://api.pumpguard.net/api/dota/download/6.jpg"
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Printf("err inn praseing url")
		log.Fatalln(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]
	log.Printf(fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("file create err")
		log.Fatalln(err)
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Printf("get req err")
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded file %s with size %d", fileName, size)

	img1, _, err := image.Decode(file)
	if err != nil {
		log.Printf("image decode err")
		log.Fatalln(err)
	}

	/*resp1, err := http.Get("http://api.tankoncloud.com/api/")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {

		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
	*/
	switch *rotate {
	case 90:
		tk.Transform = imaging.Rotate90
	case 180:
		tk.Transform = imaging.Rotate180
	case 270:
		tk.Transform = imaging.Rotate270
	}
	var dur time.Duration
	dur = 30

	tk.PlayImage(img1, dur)
	//	fatal(err)
	time.Sleep(time.Second * 1000000)
	//	close <- true
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
