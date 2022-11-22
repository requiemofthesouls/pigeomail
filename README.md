<p align="center">
      <img src="https://i.ibb.co/sPYGXGK/photo-2022-11-22-13-05-41.jpg" width="400">
</p>

<p align="center">
   <img src="https://img.shields.io/github/go-mod/go-version/requiemofthesouls/pigeomail" alt="Go version">
   <img src="https://img.shields.io/github/last-commit/requiemofthesouls/pigeomail" alt=Last commit">
   <img src="https://img.shields.io/github/license/requiemofthesouls/pigeomail" alt="License">
</p>

## About

- Service which provides securely personal email addresses written in pure Go.
- Using this service, through our telegram bot, one can create an email, receive incoming emails.
- Currently, you can only receive emails, but in future we will add sending emails via the bot as well.
- We don't store the emails on the server, messages sends directly in your telegram, check out my source code. 

## Documentation

### Setting up the project locally:

1. Take a copy of your config from default config (in .deploy folder)

      ``` cp config.dev.yaml config.yaml ```

2. Up the required containers( docker-compose located in .deploy folder)

      ``` docker-compose.yml up ```

3. Generate token from telegram by creating a bot, using [@BotFather](https://t.me/botfather)  

      Check out [this tutorial](https://docs.microsoft.com/en-us/azure/bot-service/bot-service-channel-connect-telegram?view=azure-bot-service-4.0 )

4. Build the project

      ``` go build -o pigeomail main.go ```

5. launch the services: We have two services receiver( to get incoming mail) and tg_bot (to interact with telegram bot API)

      ``` ./pigeomail receiver -c .deploy/config.yaml ```

      ``` ./pigeomail tg_bot -c .deploy/config.yaml ```


## Developers

- [Arv](https://github.com/arvryna)
- [Konstantin](https://github.com/requiemofthesouls)

