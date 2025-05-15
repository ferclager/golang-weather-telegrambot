# Golang Weather Telegram Bot
This repository contains the code for building a Weather Telegram Bot using Golang. The bot fetches weather information from the OpenWeatherMap API and sends it to a specified Telegram chat.

## Prerequisites
Before you begin, make sure you have the following:
* **Telegram Bot**: create a Telegram bot by following the instructions [here](https://core.telegram.org/bots).
* **Telegram Chat ID**: obtain the chat ID from your bot. Refer to [this guide](https://core.telegram.org/bots#3-how-do-i-create-a-bot) for more information. Also you can GET it from https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates, and see more recent chats.
* **OpenWeatherMap API Token**: register on the OpenWeatherMap website to obtain an API token. Visit [here](https://openweathermap.org/appid) to register and get your API token.
* **OpenWeatherMap City ID**: find the city ID for your desired location on the OpenWeatherMap website. For example, the city ID for Mexico City is 3530597.

* **Configuration**: configure your enviroment (file .env)
```.env
		KEY_WEATHER=YOUR_KEY
		TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN
		TELEGRAM_CHAT_ID=YOUR_CHAT_ID
```
## Usage
To use the Weather Telegram Bot, follow these steps:
* Clone the repository:

```bash
	git clone https://github.com/your-username/golangWeatherTelegramBot.git
``````
* Navigate to the project directory:
```bash
	cd golang-weather-telegramBot
````
* Run the `weatherBot.go` file with your own information. For example, to get the weather forecast for Mexico City, run the following command:
```go
	go run weatherBot.go -cityName=MexicoCity -request=WF
```

That's it! The bot will fetch the weather information from the OpenWeatherMap API and send it to the specified Telegram chat.

Feel free to modify the code according to your requirements and enhance the functionality of the bot.

Enjoy using the Weather Telegram Bot!
## License

This project is licensed under the MIT License.
