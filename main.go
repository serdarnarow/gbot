package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"

	"github.com/joho/godotenv"
)

const stickerID = "CAACAgIAAxkBAAEN0B9nsvFpiTWKBZLJ_WUnZTSN7s5DxwACfAADO2AkFCs4iGx6rGDrNgQ"

//const photoID = "fb197a65"

func messageHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	button := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Delete", CallbackData: "delete"},
				{Text: "Sticker", CallbackData: "send_sticker"},
				{Text: "Photo", CallbackData: "send_photo"},
			},
		},
	}
	_, err := ctx.EffectiveMessage.Reply(b, "wyberite deystwie", &gotgbot.SendMessageOpts{
		ReplyMarkup: button,
	})
	return err
}

// udalenie message
func deleteMessageHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.CallbackQuery == nil || ctx.CallbackQuery.Message == nil {
		return nil
	}
	_, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
	if err != nil {
		return err
	}
	_, err = ctx.CallbackQuery.Answer(b, nil)
	return err
}

// otprawka stickera
func sendStickerHandeler(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendSticker(ctx.EffectiveChat.Id, gotgbot.InputFileByID(stickerID), nil)
	if err != nil {
		return err
	}
	_, err = ctx.CallbackQuery.Answer(b, nil)
	return err
}

var photoState = make(map[int64]int)

func photoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	photos := []string{
		"/home/serdar/Downloads/ps.png",
		"/home/serdar/Downloads/zontic.png",
	}
	chatID := ctx.EffectiveChat.Id
	index := photoState[chatID] % len(photos)
	photoState[chatID]++

	file, err := os.Open(photos[index])
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = b.SendPhoto(chatID, gotgbot.InputFileByReader("file", file), &gotgbot.SendPhotoOpts{
		Caption: fmt.Sprintf("Фото %d", index+1),
	})
	if err != nil {
		return err
	}
	_, err = ctx.CallbackQuery.Answer(b, nil)
	return err

}

// func photoIDHandler(b *gotgbot.Bot, ctx *ext.Context) error {
// 	if ctx.EffectiveMessage.Photo != nil {
// 		// Берем file_id последнего (самого качественного) фото
// 		photoID := ctx.EffectiveMessage.Photo[len(ctx.EffectiveMessage.Photo)-1].FileId
// 		_, err := ctx.EffectiveMessage.Reply(b, " File ID фото:\n\n"+photoID, nil)
// 		return err
// 	}
// 	return nil
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		//fmt.Println("Ошибка загрузки .env файла")
		log.Fatal("Ошибка загрузки .env файла")
	}
	//fmt.Println(err)
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		fmt.Println("Токен не найден в .env")
		log.Fatal("Токен не найден в .env")
	}
	//sozdanie bot
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	//fmt.Println(bot)

	dispatcher := ext.NewDispatcher(nil)
	//fmt.Println(dispatcher)
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewMessage(nil, messageHandler))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("delete"), deleteMessageHandler))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("send_sticker"), sendStickerHandeler))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("send_photo"), photoHandler))

	err = updater.StartPolling(bot, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска обновлений: %v", err)
	}

	fmt.Println("Бот запущен! Нажмите Ctrl+C to stop.")
	updater.Idle()
}
