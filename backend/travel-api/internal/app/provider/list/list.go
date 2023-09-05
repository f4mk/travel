package list

import (
	"context"
	"database/sql"
	"errors"

	"github.com/f4mk/travel/backend/travel-api/internal/app/usecase/list"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	res := []RepoList{}
	q := `SELECT * from lists where user_id=$1`
	if err := r.repo.SelectContext(ctx, &res, q, uID); err != nil {
		return []list.List{}, err
	}
	ls := []list.List{}
	for _, l := range res {
		s := fromPqToString(l.ItemsID)
		ll := list.List{
			ID:          l.ID,
			UserID:      l.UserID,
			Name:        l.Name,
			Description: l.Description,
			Private:     l.Private,
			Favorite:    l.Favorite,
			Completed:   l.Completed,
			ItemsID:     s,
			DateCreated: l.DateCreated,
			DateUpdated: l.DateUpdated,
		}
		ls = append(ls, ll)
	}
	return ls, nil
}

func (r *Repo) QueryListByID(ctx context.Context, uID string, lID string) (list.List, error) {
	q := `SELECT * FROM lists WHERE list_id = $1 AND user_id = $2;`
	res := RepoList{}
	if err := r.repo.GetContext(ctx, &res, q, lID, uID); err != nil {
		return list.List{}, err
	}
	s := fromPqToString(res.ItemsID)
	l := list.List{
		ID:          res.ID,
		UserID:      res.UserID,
		Name:        res.Name,
		Description: res.Description,
		Private:     res.Private,
		Favorite:    res.Favorite,
		Completed:   res.Completed,
		ItemsID:     s,
		DateCreated: res.DateCreated,
		DateUpdated: res.DateUpdated,
	}
	return l, nil
}

func (r *Repo) QueryItemsByListID(ctx context.Context, userID string, listID string) ([]list.Item, error) {
	q := `SELECT items.*, points.point_id,
					ST_Y(points.location) AS lat, ST_X(points.location) AS lng
				FROM lists
				INNER JOIN items ON items.list_id = lists.list_id
				INNER JOIN points ON points.item_id = items.item_id
				WHERE lists.user_id = $1 AND lists.list_id = $2;`

	rows, err := r.repo.QueryxContext(ctx, q, userID, listID)
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
	q := `SELECT items.*, points.point_id,
					ST_Y(points.location) AS lat, ST_X(points.location) AS lng
				FROM lists
				INNER JOIN items ON items.list_id = lists.list_id
				INNER JOIN points ON points.item_id = items.item_id
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

func (r *Repo) CreateList(ctx context.Context, l list.List) error {
	list := populateList(l)
	q := `INSERT INTO lists (
					list_id, user_id, list_name, description, 
					private, favorite, completed, items, 
					date_created, date_updated
				) 
				VALUES (
					:list_id, :user_id, :list_name, :description, 
					:private, :favorite, :completed, :items,
					:date_created, :date_updated
				)`
	err := handleRowsResult(r.repo.NamedExecContext(ctx, q, list))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateListAdmin(ctx context.Context, l list.List) error {
	list := populateList(l)
	q := `UPDATE lists SET
					list_name = :list_name,
					description = :description,
					private = :private,
					favorite = :favorite,
					completed = :completed,
					items = :items,
					date_created = :date_created,
					date_updated = :date_updated
				WHERE list_id = :list_id;`

	err := handleRowsResult(r.repo.NamedExecContext(ctx, q, list))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateList(ctx context.Context, l list.List) error {
	list := populateList(l)
	q := `UPDATE lists SET 
					list_name = :list_name,
					description = :description,
					private = :private,
					favorite = :favorite,
					completed = :completed,
					items = :items,
					date_created = :date_created,
					date_updated = :date_updated
				WHERE list_id = :list_id AND user_id = :user_id;`

	err := handleRowsResult(r.repo.NamedExecContext(ctx, q, list))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteListAdmin(ctx context.Context, listID string) error {
	q := `DELETE FROM lists WHERE list_id = $1;`
	err := handleRowsResult(r.repo.ExecContext(ctx, q, listID))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteList(ctx context.Context, userID string, listID string) error {
	q := `DELETE FROM lists WHERE user_id = $1 AND list_id = $2;`
	err := handleRowsResult(r.repo.ExecContext(ctx, q, userID, listID))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CreateItem(ctx context.Context, userID string, i list.Item) error {
	item := populateItem(i)
	itemInsert := itemWithUserID{
		RepoItem: item,
		UserID:   userID,
	}
	point := RepoPoint{
		ID:     i.Point.ID,
		ItemID: i.Point.ItemID,
		Lat:    i.Point.Lat,
		Lng:    i.Point.Lng,
	}
	tx, err := r.repo.Beginx()
	if err != nil {
		return err
	}
	qItem := `
	INSERT INTO items (
		item_id, list_id, item_name,
		description, address,	point,
		image_links, is_visited,
		date_created,	date_updated
	)
	SELECT 	:item_id, :list_id, :item_name,
					:description, :address, :point,
					:image_links, :is_visited,
					:date_created, :date_updated
	WHERE EXISTS (
			SELECT 1 FROM lists
			WHERE lists.list_id = :list_id
			AND lists.user_id = :user_id
	);`

	// TODO: refactor magical number
	qPoint := `INSERT INTO points (
							point_id, item_id, location
						)
						VALUES (
							:point_id, :item_id, ST_SetSRID(ST_MakePoint(:lng, :lat), 4326)
						);
	`
	err = handleRowsResult(tx.NamedExecContext(ctx, qItem, itemInsert))
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			r.log.Err(rErr).Msg("Failed to rollback after error")
		}
		return err
	}

	err = handleRowsResult(tx.NamedExecContext(ctx, qPoint, point))
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			r.log.Err(rErr).Msg("Failed to rollback after error")
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			r.log.Err(rErr).Msg("Failed to rollback after error")
		}
		return err
	}
	return nil
}

func (r *Repo) UpdateItemAdmin(ctx context.Context, i list.Item) (err error) {
	item := populateItem(i)
	point := RepoPoint{
		ID:     i.Point.ID,
		ItemID: i.Point.ItemID,
		Lat:    i.Point.Lat,
		Lng:    i.Point.Lng,
	}
	tx, err := r.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				r.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()
	qItem := `UPDATE items SET
							item_name = :item_name,
							description = :description,
							address = :address,
							point = :point,
							image_links = :image_links,
							is_visited = :is_visited,
							date_created = :date_created,
							date_updated = :date_updated
						WHERE item_id = :item_id;`
	qPoint := `UPDATE points SET
							location = ST_SetSRID(ST_MakePoint(:lng, :lat), 4326)
						WHERE point_id = :point_id;`
	err = handleRowsResult(tx.NamedExecContext(ctx, qItem, item))
	if err != nil {
		return err
	}
	err = handleRowsResult(tx.NamedExecContext(ctx, qPoint, point))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateItem(ctx context.Context, userID string, i list.Item) (err error) {
	item := populateItem(i)
	itemUpdate := itemWithUserID{
		RepoItem: item,
		UserID:   userID,
	}
	point := RepoPoint{
		ID:     i.Point.ID,
		ItemID: i.Point.ItemID,
		Lat:    i.Point.Lat,
		Lng:    i.Point.Lng,
	}
	tx, err := r.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				r.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()
	qItem := `UPDATE items SET
							item_name = :item_name,
							description = :description,
							address = :address,
							point = :point,
							image_links = :image_links,
							is_visited = :is_visited,
							date_created = :date_created,
							date_updated = :date_updated
						WHERE item_id = :item_id
						AND EXISTS (
							SELECT 1 
							FROM lists 
							WHERE lists.list_id = :list_id 
							AND lists.user_id = :user_id
						);`
	qPoint := `UPDATE points SET
							location = ST_SetSRID(ST_MakePoint(:lng, :lat), 4326)
						WHERE point_id = :point_id;`
	err = handleRowsResult(tx.NamedExecContext(ctx, qItem, itemUpdate))
	if err != nil {
		return err
	}
	err = handleRowsResult(tx.NamedExecContext(ctx, qPoint, point))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteItemAdmin(ctx context.Context, itemID string) error {
	q := `DELETE FROM items WHERE item_id = $1;`
	err := handleRowsResult(r.repo.ExecContext(ctx, q, itemID))
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) DeleteItem(ctx context.Context, userID string, listID string, itemID string) error {
	q := `DELETE FROM items WHERE item_id = $1 
				AND EXISTS (
					SELECT 1
					FROM lists
					WHERE lists.list_id = $2 
					AND lists.user_id = $3
				);`
	err := handleRowsResult(r.repo.ExecContext(ctx, q, itemID, listID, userID))
	if err != nil {
		return err
	}
	return nil
}

func handleRowsResult(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func fromStringToPq(s *[]string) *pq.StringArray {
	var pqSlice *pq.StringArray
	if s != nil {
		tmp := pq.StringArray(*s)
		pqSlice = &tmp
	}
	return pqSlice
}

func fromPqToString(p *pq.StringArray) *[]string {
	var s *[]string
	if p != nil {
		temp := []string(*p)
		s = &temp
	}
	return s
}

func populateList(l list.List) RepoList {
	p := fromStringToPq(l.ItemsID)
	list := RepoList{
		ID:          l.ID,
		UserID:      l.UserID,
		Name:        l.Name,
		Description: l.Description,
		Private:     l.Private,
		Favorite:    l.Favorite,
		Completed:   l.Completed,
		ItemsID:     p,
		DateCreated: l.DateCreated,
		DateUpdated: l.DateUpdated,
	}
	return list
}

func populateItem(i list.Item) RepoItem {
	im := fromStringToPq(i.ImageLinks)
	item := RepoItem{
		ID:          i.ID,
		ListID:      i.ListID,
		Name:        i.Name,
		Description: i.Description,
		Address:     i.Address,
		PointID:     i.Point.ID,
		ImageLinks:  im,
		Visited:     i.Visited,
		DateCreated: i.DateCreated,
		DateUpdated: i.DateUpdated,
	}
	return item
}

type itemWithUserID struct {
	RepoItem
	UserID string `db:"user_id"`
}
type rowItemsByListID struct {
	RepoItem
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
			im := fromPqToString(row.ImageLinks)
			itemsMap[row.RepoItem.ID] = &list.Item{
				ID:          row.RepoItem.ID,
				ListID:      row.ListID,
				Name:        row.RepoItem.Name,
				Description: row.Description,
				Address:     row.Address,
				Point:       p,
				ImageLinks:  im,
				Visited:     row.Visited,
				DateCreated: row.DateCreated,
				DateUpdated: row.DateUpdated,
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return itemsMap, nil
}
