package auth

import "github.com/osamaNazieh/smartDesinger/internal/app"


type Handler struct {
	app *app.App
}  

func NewHandler(app *app.App) *Handler {
	return &Handler{
		app: app,  
	}
}

type Body struct {
	Username string `json:"username"`
	Password string `json:"password"`
}