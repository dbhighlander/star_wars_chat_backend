package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"star_wars/m/internal/db"
	"star_wars/m/internal/helpers"
	"star_wars/m/internal/models"

	"github.com/go-chi/chi/v5"
)

type PublicBot struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func GetBots(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetBots called") // log
	var bots []models.Bot
	if err := db.DB.Find(&bots).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response helpers.Response
	var publicBots []PublicBot

	for _, bot := range bots {
		publicBot := PublicBot{
			Name: bot.Name,
			Slug: bot.Slug,
		}

		publicBots = append(publicBots, publicBot)
	}
	response.Result = "success"
	response.Details = publicBots

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetBotBySlug(w http.ResponseWriter, r *http.Request) {
	slugParam := chi.URLParam(r, "slug")

	var bot models.Bot
	if err := db.DB.Where("slug = ? ", slugParam).First(&bot).Error; err != nil {
		http.Error(w, "Bot not found", http.StatusNotFound)
		return
	}

	var response helpers.Response
	response.Result = "success"
	response.Details = bot

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AssignBotToChat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AssignBotToChat called") // log
	slugParam := chi.URLParam(r, "bot_slug")
	chatRef := chi.URLParam(r, "chat_ref")

	var bot models.Bot
	if err := db.DB.Where("slug = ? ", slugParam).First(&bot).Error; err != nil {
		http.Error(w, "Bot not found", http.StatusNotFound)
		return
	}

	var chat models.Chat

	if err := db.DB.Where("reference = ? ", chatRef).First(&chat).Error; err != nil {
		http.Error(w, "chat not found", http.StatusNotFound)
		return
	}

	var botSystemMessage = "You are chatting to " + bot.Name
	var customerMessage models.ChatMessage

	customerMessage.Message = botSystemMessage
	customerMessage.ChatID = chat.ID
	customerMessage.MessageType = "system"

	if err := db.DB.Create(&customerMessage).Error; err != nil {
		http.Error(w, "Cannot create chat message", 500)
		return
	}

	chat.BotID = &bot.ID

	if err := db.DB.Save(&chat).Error; err != nil {
		http.Error(w, "Cannot update chat", 500)
		return
	}

	var response helpers.Response
	response.Result = "success"
	response.Details = botSystemMessage

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
