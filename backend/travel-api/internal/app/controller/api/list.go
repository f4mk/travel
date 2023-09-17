package api

import (
	"net/http"

	listService "github.com/f4mk/travel/backend/travel-api/internal/app/service/list"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/middleware"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type ListController struct {
	Log         *zerolog.Logger
	ListService *listService.Service
	Auth        *auth.Auth
	RateLimit   int
}

func NewListController(
	l *zerolog.Logger,
	ls *listService.Service,
	a *auth.Auth,
	rl int,
) *ListController {
	return &ListController{
		Log:         l,
		ListService: ls,
		Auth:        a,
		RateLimit:   rl,
	}
}

func (lc *ListController) RegisterRoutes(app *web.App) {
	app.Handle(http.MethodGet, "/lists", lc.ListService.GetLists, middleware.Authenticate(lc.Auth))
	app.Handle(http.MethodPost, "/lists", lc.ListService.CreateList, middleware.Authenticate(lc.Auth))

	app.Handle(http.MethodGet, "/lists/:listID", lc.ListService.GetList, middleware.Authenticate(lc.Auth))
	app.Handle(http.MethodPut, "/lists/:listID", lc.ListService.UpdateList, middleware.Authenticate(lc.Auth))
	app.Handle(http.MethodDelete, "/lists/:listID", lc.ListService.DeleteList, middleware.Authenticate(lc.Auth))

	app.Handle(http.MethodGet, "/lists/:listID/items", lc.ListService.GetItems, middleware.Authenticate(lc.Auth))
	app.Handle(http.MethodPost, "/lists/:listID/items", lc.ListService.CreateItem, middleware.Authenticate(lc.Auth))

	app.Handle(http.MethodGet, "/lists/:listID/items/:itemID", lc.ListService.GetItem, middleware.Authenticate(lc.Auth))
	app.Handle(http.MethodPut, "/lists/:listID/items/:itemID", lc.ListService.UpdateItem, middleware.Authenticate(lc.Auth))
	app.Handle(http.MethodDelete, "/lists/:listID/items/:itemID", lc.ListService.DeleteItem, middleware.Authenticate(lc.Auth))
}
