package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"star_wars/m/internal/db"
	"star_wars/m/internal/handlers"
	"star_wars/m/internal/helpers"
	"star_wars/m/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Ensure all tables exist
	err = testDB.AutoMigrate(
		&models.User{},
		&models.Bot{},
		&models.Chat{},
		&models.ChatMessage{},
	)
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	db.DB = testDB
}

// setupRouter sets up chi routes for both handlers
func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/bots/{slug}", handlers.GetBotBySlug)
	r.Get("/bots", handlers.GetBots)
	return r
}

// --- GetBotBySlug Tests ---

func TestGetBotBySlug_Success(t *testing.T) {
	setupTestDB(t)

	// insert a test bot
	bot := models.Bot{Slug: "test-bot", Name: "Testy"}
	assert.NoError(t, db.DB.Create(&bot).Error)

	req := httptest.NewRequest("GET", "/bots/test-bot", nil)
	rr := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response helpers.Response
	assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

	// Decode Details into Bot struct
	detailsBytes, _ := json.Marshal(response.Details)
	var botData models.Bot
	assert.NoError(t, json.Unmarshal(detailsBytes, &botData))

	assert.Equal(t, "test-bot", botData.Slug)
	assert.Equal(t, "Testy", botData.Name)
}

func TestGetBotBySlug_NotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("GET", "/bots/does-not-exist", nil)
	rr := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Bot not found")
}

// --- GetBots Tests ---

func TestGetBots_Success(t *testing.T) {
	setupTestDB(t)

	// Insert test bots
	bots := []models.Bot{
		{Slug: "bot-1", Name: "Bot One"},
		{Slug: "bot-2", Name: "Bot Two"},
	}
	for _, b := range bots {
		assert.NoError(t, db.DB.Create(&b).Error)
	}

	req := httptest.NewRequest("GET", "/bots", nil)
	rr := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response helpers.Response
	assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

	// Decode Details into slice of PublicBot
	detailsBytes, _ := json.Marshal(response.Details)
	var publicBots []handlers.PublicBot
	assert.NoError(t, json.Unmarshal(detailsBytes, &publicBots))

	assert.Len(t, publicBots, 2)
	assert.Equal(t, "bot-1", publicBots[0].Slug)
	assert.Equal(t, "Bot One", publicBots[0].Name)
	assert.Equal(t, "bot-2", publicBots[1].Slug)
	assert.Equal(t, "Bot Two", publicBots[1].Name)
}

func TestGetBots_Empty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("GET", "/bots", nil)
	rr := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response helpers.Response
	assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

	detailsBytes, _ := json.Marshal(response.Details)
	var publicBots []handlers.PublicBot
	assert.NoError(t, json.Unmarshal(detailsBytes, &publicBots))

	assert.Len(t, publicBots, 0) // no bots
}
