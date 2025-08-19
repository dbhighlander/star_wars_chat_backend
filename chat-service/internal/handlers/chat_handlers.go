package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"star_wars/m/internal/db"
	"star_wars/m/internal/gemini"
	"star_wars/m/internal/helpers"
	"star_wars/m/internal/models"
	"time"

	"github.com/go-chi/chi/v5"
)

type CreateChatResponse struct {
	UserRef        string          `json:"user_ref"`
	ChatRef        string          `json:"chat_ref"`
	BotSlug        string          `json:"bot_slug"`
	PublicMessages []PublicMessage `json:"messages"`
}

type CreateMessageResponse struct {
	CustomerMessage string `json:"customer_message"`
	BotMessage      string `json:"message"`
}

type PublicMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type MessageRequest struct {
	Message string `json:"message"`
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateChat called") // log
	botSlugParam := chi.URLParam(r, "bot_slug")

	// helpers.PrintJSON(botSlugParam)

	var bot models.Bot
	if err := db.DB.Where("slug = ? ", botSlugParam).First(&bot).Error; err != nil {
		http.Error(w, "Bot not found", http.StatusNotFound)
		return
	}

	var user models.User
	var chat models.Chat

	user.Reference = helpers.RandomString(26) + helpers.IntToStr(time.Now().UnixMilli())

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "User cannot be created", 500)
		return
	}

	chat.BotID = &bot.ID
	chat.UserID = user.ID
	chat.Reference = helpers.RandomString(26) + helpers.IntToStr(time.Now().UnixMilli())

	if err := db.DB.Create(&chat).Error; err != nil {
		http.Error(w, "Chat cannot be created", 500)
		return
	}

	var startMessage models.ChatMessage

	startMessage.BotID = &bot.ID
	startMessage.ChatID = chat.ID
	startMessage.Message = "You are chatting to " + bot.Name
	startMessage.MessageType = "system"

	if err := db.DB.Create(&startMessage).Error; err != nil {
		http.Error(w, "Chat cannot be created", 500)
		return
	}

	var PublicMessages []PublicMessage
	var publicMessage PublicMessage

	publicMessage.Type = startMessage.MessageType
	publicMessage.Message = startMessage.Message

	PublicMessages = append(PublicMessages, publicMessage)

	var createChatResponse CreateChatResponse

	createChatResponse.ChatRef = chat.Reference
	createChatResponse.UserRef = user.Reference
	createChatResponse.BotSlug = bot.Slug
	createChatResponse.PublicMessages = PublicMessages

	var response helpers.Response
	response.Result = "success"
	response.Details = createChatResponse

	// helpers.PrintJSON(response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetChat called") // log
	chatSlugParam := chi.URLParam(r, "chat_ref")
	userSlugParam := chi.URLParam(r, "user_ref")

	var user models.User
	var chat models.Chat

	if err := db.DB.Where("reference = ?", userSlugParam).First(&user).Error; err != nil {
		http.Error(w, "User cannot be found", 500)
		return
	}

	if err := db.DB.Preload("ChatMessages").Preload("Bot").Where("reference = ?", chatSlugParam).First(&chat).Error; err != nil {
		http.Error(w, "Chat cannot be found", 500)
		return
	}

	var createChatResponse CreateChatResponse

	createChatResponse.ChatRef = chat.Reference
	createChatResponse.UserRef = user.Reference
	createChatResponse.BotSlug = chat.Bot.Slug
	createChatResponse.PublicMessages = messagesToPublic(chat.ChatMessages)

	var response helpers.Response
	response.Result = "success"
	response.Details = createChatResponse

	// helpers.PrintJSON(response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Send Message called") // log
	chatSlugParam := chi.URLParam(r, "chat_ref")

	var data MessageRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var chat models.Chat

	if err := db.DB.Preload("ChatMessages").Preload("Bot").Where("reference = ?", chatSlugParam).First(&chat).Error; err != nil {
		http.Error(w, "Chat cannot be found", 500)
		return
	}

	var customerMessage models.ChatMessage

	customerMessage.Message = data.Message
	customerMessage.ChatID = chat.ID
	customerMessage.MessageType = "user"

	if err := db.DB.Create(&customerMessage).Error; err != nil {
		http.Error(w, "Cannot create chat message", 500)
		return
	}

	aiMessage := gemini.MakeApiRequest(chat, data.Message)

	var aiDbMessage models.ChatMessage

	aiDbMessage.Message = aiMessage
	aiDbMessage.ChatID = chat.ID
	aiDbMessage.MessageType = "bot"
	aiDbMessage.BotID = chat.BotID

	if err := db.DB.Create(&aiDbMessage).Error; err != nil {
		http.Error(w, "Cannot create chat message", 500)
		return
	}

	var response helpers.Response
	response.Result = "success"

	var sendMessageResponse CreateMessageResponse
	sendMessageResponse.CustomerMessage = customerMessage.Message
	sendMessageResponse.BotMessage = aiMessage

	response.Details = sendMessageResponse

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func messagesToPublic(messages []models.ChatMessage) []PublicMessage {

	var publicMessages []PublicMessage

	for _, message := range messages {
		var publicMessage PublicMessage
		publicMessage.Type = message.MessageType
		publicMessage.Message = message.Message
		publicMessages = append(publicMessages, publicMessage)
	}

	return publicMessages
}
