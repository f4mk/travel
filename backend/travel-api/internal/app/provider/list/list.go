package list

import (
	"context"
	"database/sql"
	"errors"

	"github.com/f4mk/travel/backend/travel-api/internal/app/usecase/list"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Repo struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewRepo(l *zerolog.Logger, r *sqlx.DB) *Repo {
	return &Repo{repo: r, log: l}
}

func (r *Repo) QueryListsByUserID(ctx context.Context, uID string) ([]list.List, error) {
	res := []list.List{}
	q := `SELECT * from lists where user_id=$1`
	if err := r.repo.SelectContext(ctx, &res, q, uID); err != nil {
		return []list.List{}, err
	}
	return res, nil
}

func (r *Repo) QueryListByID(ctx context.Context, uID string, lID string) (list.List, error) {
	q := `SELECT * FROM lists WHERE list_id = $1 AND user_id = $2;`
	res := list.List{}
	if err := r.repo.GetContext(ctx, &res, q, lID, uID); err != nil {
		return list.List{}, err
	}
	return res, nil
}

func (r *Repo) QueryItemsByListID(ctx context.Context, userID string, listID string) ([]list.Item, error) {
	q := `SELECT items.*, points.point_id, points.lat, 
				points.lng, links.link_id, links.name, links.url 
				FROM lists
				INNER JOIN items ON items.list_id = lists.list_id
				INNER JOIN points ON points.item_id = items.item_id
				LEFT JOIN links ON links.item_id = items.item_id
				WHERE lists.list_id = $1 AND lists.user_id = $2;`

	rows, err := r.repo.QueryxContext(ctx, q, listID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemsMap, err := fromRowsToMap(rows)
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, *item)
	}

	return items, nil
}

func (r *Repo) QueryItemByID(ctx context.Context, userID string, listID string, itemID string) (list.Item, error) {
	q := `SELECT items.*, points.point_id, points.lat, 
				points.lng, links.link_id, links.name, links.url 
				FROM lists
				INNER JOIN items ON items.list_id = lists.list_id
				INNER JOIN points ON points.item_id = items.item_id
				LEFT JOIN links ON links.item_id = items.item_id
				WHERE lists.user_id = $1 AND lists.list_id = $2 AND items.item_id = $3;`

	rows, err := r.repo.QueryxContext(ctx, q, userID, listID, itemID)
	if err != nil {
		return list.Item{}, err
	}
	defer rows.Close()

	itemsMap, err := fromRowsToMap(rows)
	if err != nil {
		return list.Item{}, err
	}

	if len(itemsMap) > 1 {
		return list.Item{}, errors.New("error querying item: duplicate records")
	}

	res, ok := itemsMap[itemID]
	if !ok {
		return list.Item{}, sql.ErrNoRows
	}
	return *res, nil
}

type rowItemsByListID struct {
	RepoItem
	RepoLink
	RepoPoint
}

func fromRowsToMap(rows *sqlx.Rows) (map[string]*list.Item, error) {
	itemsMap := make(map[string]*list.Item)

	for rows.Next() {
		row := rowItemsByListID{}

		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		if _, exists := itemsMap[row.RepoItem.ID]; !exists {

			p := list.Point{
				ID:     row.RepoPoint.ID,
				ItemID: row.RepoItem.ID,
				Lat:    row.Lat,
				Lng:    row.Lng,
			}

			itemsMap[row.RepoItem.ID] = &list.Item{
				ID:          row.RepoItem.ID,
				ListID:      row.ListID,
				Name:        row.RepoItem.Name,
				Description: row.Description,
				Address:     row.Address,
				Point:       p,
				ImageLinks:  row.ImageLinks,
				Links:       []list.Link{},
				Visited:     row.Visited,
				DateCreated: row.DateCreated,
				DateUpdated: row.DateUpdated,
			}
		}
		if row.RepoLink.ID != "" {
			link := list.Link{
				ID:     row.RepoLink.ID,
				ItemID: row.RepoLink.ItemID,
				Name:   row.RepoLink.Name,
				URL:    row.RepoLink.URL,
			}
			itemsMap[row.RepoItem.ID].Links = append(itemsMap[row.RepoItem.ID].Links, link)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return itemsMap, nil
}
