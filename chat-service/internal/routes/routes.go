package routes

import (
	"star_wars/m/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/bots", handlers.GetBots)
	r.Get("/bots/{slug}", handlers.GetBotBySlug)
	r.Get("/chats/create/{bot_slug}", handlers.CreateChat)
	r.Get("/chats/{user_ref}/{chat_ref}", handlers.GetChat)
	r.Post("/message/{chat_ref}", handlers.SendMessage)
	r.Get("/chats/switch_bot/{bot_slug}/{chat_ref}", handlers.AssignBotToChat)

}
