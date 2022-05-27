# Pomodoro Go Bot

Pomodoro Go Bot is a telegram bot that helps you to keep track of working and chilling time.

## Installation

Just pull the repository, or get the binary. Btw, the bot requires a telegram bot token to be exported as an environment variable `TOKEN`

```bash
export TOKEN=123456789:AABBCCddEF4Nood_SO_H3R3ISMyTokEN
```

## Usage
You can either run the binary file:

```bash
./pomodoro
```
Or you can launch it using docker:
```bash
docker build -t pomodoro_bot . && docker run -d --name pomodoro
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://github.com/ytanne/pomodoro_bot/blob/main/LICENSE.txt)
