package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	helpCommand  = "help"
	startCommand = "start"
)

func (b *Bot) handleHelpCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	txt := `
Hello! I'm open-source bot which provides securely personal email addresses right here, in telegram.  
I don't keep your messages on server, all messages forwards immediately to you.  
Check out my source code [here](https://github.com/requiemofthesouls/pigeomail).  

*Available domains*
	- %s

*My commands*
	/create - Create new email
	/list   - Show your email
	/delete - Delete your email
	/help   - Get help message
`

	msg.Text = fmt.Sprintf(txt, b.domain)
	msg.ParseMode = "markdown"

	_, _ = b.api.Send(msg)
}
