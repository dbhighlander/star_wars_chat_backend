package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"star_wars/m/internal/db"
	"star_wars/m/internal/handlers"
	"star_wars/m/internal/helpers"
	"star_wars/m/internal/models"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Helper to create a router with the chat routes
func setupChatRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/chat/{bot_slug}", handlers.CreateChat)
	r.Get("/chat/{chat_ref}/{user_ref}", handlers.GetChat)
	r.Post("/chat/message/{chat_ref}", handlers.SendMessage)
	return r
}

func TestCreateChat_Success(t *testing.T) {
	setupTestDB(t)

	// Insert a bot for testing
	bot := models.Bot{Slug: "test-bot", Name: "Test Bot"}
	assert.NoError(t, db.DB.Create(&bot).Error)

	req := httptest.NewRequest("POST", "/chat/test-bot", nil)
	rr := httptest.NewRecorder()

	router := setupChatRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	var response helpers.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, rr.Body.String())

	// Parse nested CreateChatResponse
	detailsBytes, _ := json.Marshal(response.Details)
	var chatResp handlers.CreateChatResponse
	err = json.Unmarshal(detailsBytes, &chatResp)
	assert.NoError(t, err, rr.Body.String())

	assert.Equal(t, "test-bot", chatResp.BotSlug)
	assert.NotEmpty(t, chatResp.ChatRef)
	assert.NotEmpty(t, chatResp.UserRef)
	assert.Len(t, chatResp.PublicMessages, 1)
	assert.Equal(t, "system", chatResp.PublicMessages[0].Type)
}

func TestGetChat_Success(t *testing.T) {
	setupTestDB(t)

	// Seed bot, user, chat, and message
	bot := models.Bot{Slug: "test-bot", Name: "Test Bot"}
	db.DB.Create(&bot)

	user := models.User{Reference: "user123"}
	db.DB.Create(&user)

	chat := models.Chat{
		BotID:     &bot.ID,
		UserID:    user.ID,
		Reference: "chat123",
	}
	db.DB.Create(&chat)

	msg := models.ChatMessage{
		ChatID:      chat.ID,
		BotID:       &bot.ID,
		Message:     "Hello",
		MessageType: "system",
	}
	db.DB.Create(&msg)

	req := httptest.NewRequest("GET", "/chat/chat123/user123", nil)
	rr := httptest.NewRecorder()

	router := setupChatRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	var response helpers.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	detailsBytes, _ := json.Marshal(response.Details)
	var chatResp handlers.CreateChatResponse
	err = json.Unmarshal(detailsBytes, &chatResp)
	assert.NoError(t, err)

	assert.Equal(t, "test-bot", chatResp.BotSlug)
	assert.Equal(t, "chat123", chatResp.ChatRef)
	assert.Equal(t, "user123", chatResp.UserRef)
	assert.Len(t, chatResp.PublicMessages, 1)
	assert.Equal(t, "system", chatResp.PublicMessages[0].Type)
}

func TestCreateChat_BotNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("POST", "/chat/nonexistent-bot", nil)
	rr := httptest.NewRecorder()

	router := setupChatRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
