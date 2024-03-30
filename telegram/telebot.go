package telebot

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TELE_API_TOKEN    = "TELE_API_TOKEN"
	TELE_HOST_CHATBOT = "TELE_HOST_CHATBOT"
)

type Telebot struct {
	Bot         *tgbotapi.BotAPI
	HostChatbot string
}

func New() *Telebot {
	tat := ""
	if at := os.Getenv(TELE_API_TOKEN); at != "" {
		tat = at
	}

	bot, err := tgbotapi.NewBotAPI(tat)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	botUpdate := tgbotapi.NewUpdate(0)
	botUpdate.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates, err := bot.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Println(err)
	}

	hostChatbot := ""
	if hcb := os.Getenv(TELE_HOST_CHATBOT); hcb != "" {
		hostChatbot = hcb
	}

	retFT := &Telebot{
		Bot:         bot,
		HostChatbot: hostChatbot,
	}
	go receiveUpdates(retFT, ctx, updates)
	log.Println("Start listening for updates. Press ctrl+c to stop")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	cancel()

	return retFT
}

func receiveUpdates(h *Telebot, ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			h.handleUpdate(update)
		}
	}
}

func (h *Telebot) handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		h.handleMessage(update.Message)
		break

		// handle telegram menu
		// case update.CallbackQuery != nil:
		// 	handleButton(update.CallbackQuery)
		// 	break
	}
}

// message from user is here ..
func (h *Telebot) handleMessage(message *tgbotapi.Message) {

	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	const cred string = `{
							"email": "",
							"pin": ""
						}`
	// token := h.requestPINChatbot(cred)

	// chatResponse := h.requestMsgChatbot(message.Text, token)
	chatResponse := h.requestMsgChatbotText(message.Text, message)

	msg := tgbotapi.NewMessage(message.Chat.ID, chatResponse)
	_, err := h.Bot.Send(msg)

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}

}

func (h *Telebot) requestPINChatbot(credential string) string {
	var err error
	var client = &http.Client{}
	payload := bytes.NewBufferString(credential)
	// request token
	request, err := http.NewRequest("POST", h.HostChatbot+"/v1/auth/login-pin", payload)
	if err != nil {
		log.Println(err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("App-Name", "Park Spectrum Ventures")
	request.Header.Set("Character-Name", "Sauna")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
	}
	defer response.Body.Close()

	lgnpinResp := LoginpinResp{}
	err = json.NewDecoder(response.Body).Decode(&lgnpinResp)
	if err != nil {
		log.Println(err.Error())
	}

	return lgnpinResp.Data.Token
}

func (h *Telebot) requestMsgChatbot(message, token string) string {
	var err error
	var client = &http.Client{}

	payload := bytes.NewBufferString(`{"content":"` + message + `", "name":"lola"}`)
	// request token
	request, err := http.NewRequest("POST", h.HostChatbot+"/v1/ai-text/message", payload)
	if err != nil {
		log.Println(err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("App-Name", "Park Spectrum Ventures")
	request.Header.Set("Character-Name", "Sauna")
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
	}
	defer response.Body.Close()

	chatResp := ChatResp{}
	err = json.NewDecoder(response.Body).Decode(&chatResp)
	if err != nil {
		log.Println(err.Error())
	}

	return chatResp.Data.Data.Content
}

func (h *Telebot) requestMsgChatbotText(message string, messageApi *tgbotapi.Message) string {
	var err error
	var client = &http.Client{}

	payload := bytes.NewBufferString(`{"content":"` + message + `", "name":"` + messageApi.From.FirstName + `"}`)
	// request token
	request, err := http.NewRequest("POST", h.HostChatbot+"/v1/ai-text/message", payload)
	if err != nil {
		log.Println(err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
	}
	defer response.Body.Close()

	chatResp := ChatTextResp{}
	err = json.NewDecoder(response.Body).Decode(&chatResp)
	if err != nil {
		log.Println(err.Error())
	}

	return chatResp.Data
}
