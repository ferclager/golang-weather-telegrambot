package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const MAX = 3

var cities = map[string]string{
	"Madrid":     "3117735",
	"MexicoCity": "3530597",
	"NewYork":    "5128581",
	"Toronto":    "6167865",
}

type ForecastAPIResponse struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Visibility int     `json:"visibility"`
		Pop        float64 `json:"pop"`
		Sys        struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
		Rain  struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

type WeatherAPIResponse struct {
	Base   string `json:"base"`
	Clouds struct {
		All int64 `json:"all"`
	} `json:"clouds"`
	Cod   int64 `json:"cod"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Dt   int64 `json:"dt"`
	ID   int64 `json:"id"`
	Main struct {
		FeelsLike float64 `json:"feels_like"`
		Humidity  int64   `json:"humidity"`
		Pressure  int64   `json:"pressure"`
		Temp      float64 `json:"temp"`
		TempMax   float64 `json:"temp_max"`
		TempMin   float64 `json:"temp_min"`
	} `json:"main"`
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		ID      int64  `json:"id"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
		Type    int64  `json:"type"`
	} `json:"sys"`
	Timezone   int64 `json:"timezone"`
	Visibility int64 `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		ID          int64  `json:"id"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Deg   int64   `json:"deg"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	KEY_WEATHER := os.Getenv("KEY_WEATHER")
	TELEGRAM_BOT_TOKEN := os.Getenv("TELEGRAM_BOT_TOKEN")
	TELEGRAM_CHAT_ID := os.Getenv("TELEGRAM_CHAT_ID")

	cityName := *flag.String("cityName", "MexicoCity", "City")
	option := *flag.String("request", "WF", "Type")
	flag.Parse()

	var idCity = cities[cityName]
	if idCity != "" {
		now := time.Now()
		message := "WeatherBot " + now.Format(time.Kitchen)
		switch option {
		case "W":
			message += callAPI(KEY_WEATHER, "weather", idCity)
		case "F":
			message += callAPI(KEY_WEATHER, "forecast", idCity)
		case "WF":
			message += callAPI(KEY_WEATHER, "weather", idCity) + callAPI(KEY_WEATHER, "forecast", idCity)
		default:
			message = ""
		}
		if message != "" {
			sendMessage(TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID, message)
		}
	}
}

func callAPI(key string, endPoint string, idCity string) string {
	url := "https://api.openweathermap.org/data/2.5/" + endPoint + "?id=" + idCity + "&APPID=" + key + "&lang=en"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Weather request was not OK: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var result = ""
	if endPoint == "weather" {
		result = parseResponseWeather(body)
	} else if endPoint == "forecast" {
		result = parseResponseForecast(body)
	}
	return result
}

func sendMessage(telBotToken string, telChat string, message string) {
	url := "https://api.telegram.org/bot" + telBotToken + "/sendMessage?chat_id=" + telChat + "&text=" + url.QueryEscape(message)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Telegram request was not OK: " + resp.Status)
	}
}

func parseResponseWeather(body []byte) string {
	var s = new(WeatherAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}
	var message = "\n" + s.Name + " (" + s.Sys.Country + ") " + s.Weather[0].Description + ". 🌡  Temperature (ºC) " + fmt.Sprintf("%.0f", toCelsius(s.Main.Temp)) + ". Feels like " + fmt.Sprintf("%.0f", toCelsius(s.Main.FeelsLike)) + " (L " + fmt.Sprintf("%.0f", toCelsius(s.Main.TempMin)) + " - H " + fmt.Sprintf("%.0f", toCelsius(s.Main.TempMax)) + "), " + fmt.Sprintf("%2d", s.Main.Humidity) + " humidity. "
	return message
}

func parseResponseForecast(body []byte) string {
	var s = new(ForecastAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}
	var message = ""
	for _, value := range s.List[0:MAX] {
		var partialMessage = " ⌚️ " + value.DtTxt + " " + value.Weather[0].Description + ". 🌡  Temperature (ºC) " + fmt.Sprintf("%.0f", toCelsius(value.Main.Temp)) + ". Feels like " + fmt.Sprintf("%.0f", toCelsius(value.Main.FeelsLike)) + " (L " + fmt.Sprintf("%.0f", toCelsius(value.Main.TempMin)) + " - H " + fmt.Sprintf("%.0f", toCelsius(value.Main.TempMax)) + "), " + fmt.Sprintf("%2d", value.Main.Humidity) + " humidity."
		message += "\n" + partialMessage
	}
	return message
}

func toCelsius(kelvin float64) float64 {
	return math.Round(kelvin - 273.15)
}

func getStrCities(cities map[string]string) string {
	var message = ""
	for key, _ := range cities {
		message += key + "|"
	}
	message = message[0 : len(message)-1]
	return message
}
