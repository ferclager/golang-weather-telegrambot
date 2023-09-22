package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const KEY_WEATHER = "YOUR_KEY"
const ID_CITY = "3530597" //EXAMPLE: MEXICO CITY
const TELEGRAM_BOT_TOKEN = "YOUR_BOT_TOKEN"
const TELEGRAM_CHAT_ID = "YOUR_CHAT_ID"
const MAX = 9

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
	now := time.Now()
	message := "WeatherBot " + now.Format(time.Kitchen) + "\n" + callAPI("weather") + callAPI("forecast")
	fmt.Printf(message)
	sendMessage(message)
}

func callAPI(endPoint string) string {
	url := "https://api.openweathermap.org/data/2.5/" + endPoint + "?id=" + ID_CITY + "&APPID=" + KEY_WEATHER + "&lang=en"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Request was not OK: " + resp.Status)
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

func sendMessage(message string) {
	url := "https://api.telegram.org/bot" + TELEGRAM_BOT_TOKEN + "/sendMessage?chat_id=" + TELEGRAM_CHAT_ID + "&text=" + url.QueryEscape(message)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Request was not OK: " + resp.Status)
	}
}

func parseResponseWeather(body []byte) string {
	var s = new(WeatherAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}
	var message = s.Name + " (" + s.Sys.Country + ") " + s.Weather[0].Description + ". 🌡  Temperature (ºC) " + fmt.Sprintf("%.2f", toCelsius(s.Main.Temp)) + " (min " + fmt.Sprintf("%.2f", toCelsius(s.Main.TempMin)) + " - max " + fmt.Sprintf("%.2f", toCelsius(s.Main.TempMax)) + "). Feels like " + fmt.Sprintf("%.2f", toCelsius(s.Main.FeelsLike)) + ", " + fmt.Sprintf("%2d", s.Main.Humidity) + " humidity. "
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
		var partialMessage = " ⌚️ " + value.DtTxt + " " + value.Weather[0].Description + ". 🌡  Temperature (ºC) " + fmt.Sprintf("%.2f", toCelsius(value.Main.Temp)) + " (min " + fmt.Sprintf("%.2f", toCelsius(value.Main.TempMin)) + " - max " + fmt.Sprintf("%.2f", toCelsius(value.Main.TempMax)) + "). Feels like " + fmt.Sprintf("%.2f", toCelsius(value.Main.FeelsLike)) + ", " + fmt.Sprintf("%2d", value.Main.Humidity) + " humidity."
		message += "\n" + partialMessage
	}
	return message
}

func toCelsius(kelvin float64) float64 {
	return (kelvin - 273.15)
}
