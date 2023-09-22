# Building a weather Telegram Bot using golang
All you will need to start is:
* Your telegram bot: [create one](https://core.telegram.org/bots).
* Your telegram chat ID from your bot: [check this out for more information](https://core.telegram.org/bots#3-how-do-i-create-a-bot).
* An OpenWeatherMap API token: [register & get one](https://openweathermap.org/appid).
* Your OpenWeatherMap city ID; for example, Mexico City 3530597.
* Configure your enviroment (file .env)
```.env
	KEY_WEATHER=YOUR_KEY
	TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN
	TELEGRAM_CHAT_ID=YOUR_CHAT_ID
```
Now, just use your own information and run the code. Example:
```
go run weatherBot.go -cityName=MexicoCity -request=WF
```
Enjoy it!
