package gemini

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"star_wars/m/internal/helpers"
	"star_wars/m/internal/models"

	"github.com/joho/godotenv"
)

type GeminiRequest struct {
	Contents          []GeminiMessage         `json:"contents"`
	SystemInstruction GeminiSystemInstruction `json:"system_instruction"`
}

type GeminiSystemInstruction struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiMessage struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiResponseCandidate `json:"candidates"`
}

type GeminiResponseCandidate struct {
	Content GeminiMessage `json:"content"`
}

// type GeminiResponseCandidateParts struct {
// 	Parts []GeminiMessage `json:"parts"`
// }

func MakeApiRequest(chat models.Chat, message string) string {
	requestBody := assembleRequestBody(chat, message)

	helpers.PrintJSON(requestBody)

	// jsonData := []byte(`{"name":"John", "age":30}`)
	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewBufferString(requestBody))
	if err != nil {
		panic(err)
	}

	err = godotenv.Load("./.env.secrets")
	if err != nil {
		log.Println("Warning: .env.secrets not loaded")
	}
	apiKey := os.Getenv("GEMINI_API_KEY")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Optionally parse the JSON error
		var apiError struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &apiError); err != nil {
			// fallback if JSON decoding fails
			log.Fatal(string(body))
		}
		log.Fatal(apiError)
		return ""
	}

	var geminiResponse GeminiResponse
	err = json.Unmarshal([]byte(body), &geminiResponse)
	if err != nil {
		log.Fatal(err)
	}

	geminiMessage := geminiResponse.Candidates[0].Content.Parts[0].Text

	return geminiMessage
}

func assembleRequestBody(chat models.Chat, message string) string {
	var geminiRequest GeminiRequest
	var systemInstruction GeminiSystemInstruction
	var systemInstructionPart GeminiPart

	systemInstructionPart.Text = chat.Bot.PersonalityPrompt
	systemInstruction.Parts = append(systemInstruction.Parts, systemInstructionPart)
	geminiRequest.SystemInstruction = systemInstruction

	for _, message := range chat.ChatMessages {
		var geminiParts []GeminiPart
		var geminiPart GeminiPart
		var geminiMessage GeminiMessage

		if message.MessageType != "system" {
			geminiPart.Text = message.Message
			geminiParts = append(geminiParts, geminiPart)

			switch message.MessageType {
			case "user":
				geminiMessage.Role = "user"
				geminiMessage.Parts = geminiParts
			case "bot":
				geminiMessage.Role = "model"
				geminiMessage.Parts = geminiParts
			}

			geminiRequest.Contents = append(geminiRequest.Contents, geminiMessage)

		}

	}

	var customerMessage GeminiMessage
	var geminiParts []GeminiPart
	var geminiPart GeminiPart

	geminiPart.Text = message
	geminiParts = append(geminiParts, geminiPart)
	customerMessage.Role = "user"
	customerMessage.Parts = geminiParts

	geminiRequest.Contents = append(geminiRequest.Contents, customerMessage)

	jsonData, err := json.Marshal(geminiRequest)
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
