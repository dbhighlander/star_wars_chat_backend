package routes

import (
	"star_wars/m/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/bots", handlers.GetBots)
	r.Get("/bots/{slug}", handlers.GetBotBySlug)
	r.Get("/chats/{user_ref}/{chat_ref}", handlers.GetChat)
	r.Get("/chats/create", handlers.CreateChat)
	r.Post("/message/{chat_ref}", handlers.SendMessage)

}
