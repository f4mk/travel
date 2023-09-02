package api

import (
	"net/http"

	listService "github.com/f4mk/travel/backend/travel-api/internal/app/service/list"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type ListController struct {
	Log         *zerolog.Logger
	ListService *listService.Service
	Auth        *auth.Auth
	RateLimit   int
}

func (lc *ListController) RegisterRoutes(app *web.App) {
	app.Handle(http.MethodGet, "/lists", lc.ListService.GetLists)
	app.Handle(http.MethodPost, "/lists", lc.ListService.CreateList)

	app.Handle(http.MethodGet, "/lists/:id", lc.ListService.GetList)
	app.Handle(http.MethodPut, "/lists/:id", lc.ListService.UpdateList)
	app.Handle(http.MethodDelete, "/lists/:id", lc.ListService.DeleteList)

	app.Handle(http.MethodGet, "/lists/:listID/items", lc.ListService.GetItems)
	app.Handle(http.MethodPost, "/lists/:listID/items", lc.ListService.CreateItem)

	app.Handle(http.MethodGet, "/lists/:listID/items/:itemID", lc.ListService.GetItem)
	app.Handle(http.MethodPut, "/lists/:listID/items/:itemID", lc.ListService.UpdateItem)
	app.Handle(http.MethodDelete, "/lists/:listID/items/:itemID", lc.ListService.DeleteItem)
}
