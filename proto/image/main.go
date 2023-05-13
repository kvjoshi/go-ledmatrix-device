package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/canvas"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/matrix"
	"image"
	"io"
	"io/ioutil"
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
	parallel               = flag.Int("led-parallel", 3, "number of daisy-chained panels")
	chain                  = flag.Int("led-chain", 5, "number of displays daisy-chained")
	brightness             = flag.Int("brightness", 100, "brightness (0-100)")
	gpio_slowdown          = flag.Int("led-gpio-slowdown", 3, "GPIO SLOWDOWN")
	pwm_lsb                = flag.Int("led-pwm-lsb-nanoseconds", 70, "lsb nanosec")
	pwm_bits               = flag.Int("led-pwm-bits", 9, "pwm bits")
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

func fetchImg(imgUrl string) image.Image {
	fullURLFile = "http://api.pumpguard.net/api/dota/download/" + imgUrl
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Printf("err in prasing url")
		log.Fatalln(err)
	}
	path := fileURL.Path
	log.Printf(path)
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
	err = os.Chmod(fileName, 0777)
	if err != nil {
		log.Printf("err in chmod")
		log.Fatalln(err)
	}

	fmt.Printf("Downloaded file %s with size %d", fileName, size)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	img1, _, err := image.Decode(f)
	if err != nil {
		log.Printf("image decode err")
		log.Fatalln(err)
	}
	return img1

}

type Schedule struct {
	ContentName string
	ContentPath string
	Time        string
}

func getContentSchedule() []Schedule {
	getScheduleUrl := "http://142.93.198.132:3000/api/sch/getScheduleBySidD"
	//getScheduleUrl := "http://192.168.1.2:3000/api/sch/getScheduleBySidD"
	/* scheduleURL, err := url.Parse(getScheduleUrl)
	if err != nil {
		log.Printf("err in parsing schedule url")
	} */
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	/*
	   	scheduleRequestBody := `{
	       "scheduleId":"645255100283e16678c9e609"
	   	}`*/
	resp, err := client.PostForm(getScheduleUrl, url.Values{"scheduleId": {"645255100283e16678c9e609"}})
	if err != nil {
		log.Printf("get req err")
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//bodyString := string(body)
	//log.Printf("body:")
	//log.Print(bodyString)
	var schedule []Schedule
	err2 := json.Unmarshal(body, &schedule)
	if err2 != nil {
		fmt.Println("error:", err2)
		os.Exit(1)
	}
	//fmt.Print(schedule)
	for _, a := range schedule {
		fmt.Println(a)
	}
	return schedule
}

func show_img(schedule Schedule) {

}
func main() {
	//f, err := os.Open(*img)
	//fatal(err)
	//img, _, err := image.Decode(f)

	schedule := getContentSchedule()
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
	switch *rotate {
		case 90:
			tk.Transform = imaging.Rotate90
		case 180:
			tk.Transform = imaging.Rotate180
		case 270:
			tk.Transform = imaging.Rotate270
		}
	*/
	for _, a := range schedule {
		print(a.ContentPath)

	}

	img1 := fetchImg(schedule[0].ContentPath)
	time.Sleep(time.Second * 100)
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
