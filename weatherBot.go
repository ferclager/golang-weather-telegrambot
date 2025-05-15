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

	"github.com/ferclager/golang-weather-telegrambot/models"

	"github.com/joho/godotenv"
)

const MAX = 3

var cities = map[string]string{
	"Madrid":     "3117735",
	"MexicoCity": "3530597",
	"NewYork":    "5128581",
	"Toronto":    "6167865",
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	KEY_WEATHER := os.Getenv("KEY_WEATHER")
	TELEGRAM_BOT_TOKEN := os.Getenv("TELEGRAM_BOT_TOKEN")
	TELEGRAM_CHAT_ID := os.Getenv("TELEGRAM_CHAT_ID")

	cityName := flag.String("cityName", "MexicoCity", "City")
	option := flag.String("request", "WF", "Type")
	flag.Parse()

	var idCity = cities[*cityName]
	if idCity != "" {
		now := time.Now()
		message := "WeatherBot " + now.Format(time.Kitchen)
		switch *option {
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
	var s = new(models.WeatherAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}
	var message = "\n" + s.Name + " (" + s.Sys.Country + ") " + s.Weather[0].Description + ". 🌡  Temperature (ºC) " + fmt.Sprintf("%.0f", toCelsius(s.Main.Temp)) + ". Feels like " + fmt.Sprintf("%.0f", toCelsius(s.Main.FeelsLike)) + " (L " + fmt.Sprintf("%.0f", toCelsius(s.Main.TempMin)) + " - H " + fmt.Sprintf("%.0f", toCelsius(s.Main.TempMax)) + "), " + fmt.Sprintf("%2d", s.Main.Humidity) + " humidity. "
	return message
}

func parseResponseForecast(body []byte) string {
	var s = new(models.ForecastAPIResponse)
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
