package api

import "github.com/osamaNazieh/smartDesinger/internal/app"

type Handler struct {
	app *app.App
}

// NewHandler creates and returns a new Handler instance.
// It initializes the Handler with the provided App instance.
func NewHandler(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}


// Transformation Options 
type Resize struct {
	Width int `json:"width"`
	Height int `json:"height"`
}

type Crop struct {
	Width int `json:"width"`
	Height int `json:"height"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Filters struct {
	Grayscale bool `json:"grayscale"`
	Dither bool `json:"dither"`
}

type Transformations struct {
	ResizeOption Resize `json:"resize"`
	CropOption Crop `json:"crop"`
	Rotate int `json:"rotate"`
	Format string `json:"format"`
	FiltersOptions Filters `json:"filters"`
} 