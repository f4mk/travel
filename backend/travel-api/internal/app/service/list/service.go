package list

import (
	"context"
	"net/http"

	listUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/list"
	"github.com/rs/zerolog"
)

// TODO: remove when done
//
//revive:disable
type Service struct {
	core *listUsecase.Core
	log  *zerolog.Logger
}

func NewService(l *zerolog.Logger, c *listUsecase.Core) *Service {

	return &Service{
		core: c,
		log:  l,
	}
}

// GetLists retrieves all lists.
func (s *Service) GetLists(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// CreateList creates a new list.
func (s *Service) CreateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// GetList retrieves a specific list by its ID.
func (s *Service) GetList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// UpdateList updates a specific list by its ID.
func (s *Service) UpdateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// DeleteList deletes a specific list by its ID.
func (s *Service) DeleteList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// GetItems retrieves all items for a specific list.
func (s *Service) GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// CreateItem creates a new item for a specific list.
func (s *Service) CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// GetItem retrieves a specific item by its ID for a specific list.
func (s *Service) GetItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// UpdateItem updates a specific item by its ID for a specific list.
func (s *Service) UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// DeleteItem deletes a specific item by its ID for a specific list.
func (s *Service) DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}
