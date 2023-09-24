package list

import (
	"context"
	"database/sql"
	"errors"

	"github.com/f4mk/travel/backend/travel-api/internal/app/usecase/list"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/images"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

const EPSG = 4326

type Storer struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewStorer(l *zerolog.Logger, r *sqlx.DB) *Storer {
	return &Storer{repo: r, log: l}
}

func (s *Storer) QueryListsByUserID(ctx context.Context, uID string) ([]list.List, error) {
	res := []StorerList{}
	q := `SELECT * from lists where user_id=$1`
	if err := s.repo.SelectContext(ctx, &res, q, uID); err != nil {
		return []list.List{}, err
	}
	ls := []list.List{}
	for _, l := range res {
		s := []string(l.ItemsID)
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

func (s *Storer) QueryListByID(ctx context.Context, uID string, lID string) (list.List, error) {
	q := `SELECT * FROM lists WHERE list_id = $1 AND user_id = $2;`
	res := StorerList{}
	if err := s.repo.GetContext(ctx, &res, q, lID, uID); err != nil {
		return list.List{}, err
	}
	r := []string(res.ItemsID)
	l := list.List{
		ID:          res.ID,
		UserID:      res.UserID,
		Name:        res.Name,
		Description: res.Description,
		Private:     res.Private,
		Favorite:    res.Favorite,
		Completed:   res.Completed,
		ItemsID:     r,
		DateCreated: res.DateCreated,
		DateUpdated: res.DateUpdated,
	}
	return l, nil
}

func (s *Storer) QueryItemsByListID(ctx context.Context, userID string, listID string) ([]list.Item, error) {
	q := `
		SELECT items.*, points.point_id,
			ST_Y(points.location) AS lat, ST_X(points.location) AS lng
		FROM lists
		INNER JOIN items ON items.list_id = lists.list_id
		INNER JOIN points ON points.item_id = items.item_id
		WHERE lists.user_id = $1 AND lists.list_id = $2;
	`

	rows, err := s.repo.QueryxContext(ctx, q, userID, listID)
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

func (s *Storer) QueryItemByID(ctx context.Context, userID string, listID string, itemID string) (list.Item, error) {
	q := `
		SELECT items.*, points.point_id,
			ST_Y(points.location) AS lat, ST_X(points.location) AS lng
		FROM lists
		INNER JOIN items ON items.list_id = lists.list_id
		INNER JOIN points ON points.item_id = items.item_id
		WHERE lists.user_id = $1 AND lists.list_id = $2 AND items.item_id = $3;
	`

	rows, err := s.repo.QueryxContext(ctx, q, userID, listID, itemID)
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

func (s *Storer) CreateList(ctx context.Context, l list.List) error {
	list := populateList(l)
	q := `
		INSERT INTO lists (
			list_id, user_id, list_name, description, 
			private, favorite, completed, items, 
			date_created, date_updated
		) 
		VALUES (
			:list_id, :user_id, :list_name, :description, 
			:private, :favorite, :completed, :items,
			:date_created, :date_updated
		);
	`

	err := handleRowsResult(s.repo.NamedExecContext(ctx, q, list))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) UpdateListAdmin(ctx context.Context, l list.List) error {
	list := populateList(l)
	q := `
		UPDATE lists SET
			list_name = :list_name,
			description = :description,
			private = :private,
			favorite = :favorite,
			completed = :completed,
			items = :items,
			date_created = :date_created,
			date_updated = :date_updated
		WHERE list_id = :list_id;
	`

	err := handleRowsResult(s.repo.NamedExecContext(ctx, q, list))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) UpdateList(ctx context.Context, l list.List) error {
	list := populateList(l)
	q := `
		UPDATE lists SET 
			list_name = :list_name,
			description = :description,
			private = :private,
			favorite = :favorite,
			completed = :completed,
			items = :items,
			date_created = :date_created,
			date_updated = :date_updated
		WHERE list_id = :list_id AND user_id = :user_id;
	`

	err := handleRowsResult(s.repo.NamedExecContext(ctx, q, list))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) DeleteListAdmin(ctx context.Context, listID string) error {
	q := `DELETE FROM lists WHERE list_id = $1;`
	err := handleRowsResult(s.repo.ExecContext(ctx, q, listID))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) DeleteList(ctx context.Context, userID string, listID string) error {
	q := `DELETE FROM lists WHERE user_id = $1 AND list_id = $2;`
	err := handleRowsResult(s.repo.ExecContext(ctx, q, userID, listID))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) CreateItem(ctx context.Context, i list.Item) (err error) {
	item := populateItem(i)
	point := StorerPoint{
		ID:     i.Point.ID,
		ItemID: i.Point.ItemID,
		Lat:    i.Point.Lat,
		Lng:    i.Point.Lng,
		EPSG:   EPSG,
	}
	tx, err := s.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				s.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()
	qImages := `
		UPDATE images SET
			item_id = $1,
			status = $2,
		WHERE image_id = $3 
		AND list_id = $4;
	`
	qItem := `
		INSERT INTO items (
			item_id, list_id, user_id, item_name,
			description, address,	point,
			images_id, is_visited,
			date_created,	date_updated
		)
		SELECT 	:item_id, :list_id, :user_id, :item_name,
						:description, :address, :point,
						:images_id, :is_visited,
						:date_created, :date_updated
		WHERE EXISTS (
				SELECT 1 FROM lists
				WHERE lists.list_id = :list_id
				AND lists.user_id = :user_id
		);
	`
	qPoint := `
		INSERT INTO points (
			point_id, item_id, location
		)
		VALUES (
			:point_id, :item_id, ST_SetSRID(ST_MakePoint(:lng, :lat), :epsg)
		);
	`
	err = handleRowsResult(tx.NamedExecContext(ctx, qItem, item))
	if err != nil {
		return err
	}

	for _, imageID := range item.ImagesID {
		_, err := s.repo.ExecContext(ctx, qImages, item.ID, images.Loaded, imageID, item.ListID)
		if err != nil {
			return err
		}
	}

	err = handleRowsResult(tx.NamedExecContext(ctx, qPoint, point))
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (s *Storer) UpdateItemAdmin(ctx context.Context, i list.Item, toDelete []string) (err error) {
	item := populateItem(i)
	point := StorerPoint{
		ID:     i.Point.ID,
		ItemID: i.Point.ItemID,
		Lat:    i.Point.Lat,
		Lng:    i.Point.Lng,
		EPSG:   EPSG,
	}
	tx, err := s.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				s.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()
	qImages := `
		UPDATE images SET status = $1 WHERE image_id = $2 AND item_id = $3;
	`
	qItem := `
		UPDATE items SET
			item_name = :item_name,
			description = :description,
			address = :address,
			point = :point,
			images_id = :images_id,
			is_visited = :is_visited,
			date_created = :date_created,
			date_updated = :date_updated
		WHERE item_id = :item_id;
	`
	qPoint := `
		UPDATE points SET
			location = ST_SetSRID(ST_MakePoint(:lng, :lat), :epsg)
		WHERE point_id = :point_id;
	`

	for _, imgID := range item.ImagesID {
		_, err = tx.ExecContext(ctx, qImages, images.Loaded, imgID, item.ID)
		if err != nil {
			return err
		}
	}
	for _, imgID := range toDelete {
		_, err = tx.ExecContext(ctx, qImages, images.Deleted, imgID, item.ID)
		if err != nil {
			return err
		}
	}
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

func (s *Storer) UpdateItem(ctx context.Context, i list.Item, toDelete []string) (err error) {
	item := populateItem(i)
	point := StorerPoint{
		ID:     i.Point.ID,
		ItemID: i.Point.ItemID,
		Lat:    i.Point.Lat,
		Lng:    i.Point.Lng,
		EPSG:   EPSG,
	}
	tx, err := s.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				s.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()
	qImages := `
		UPDATE images SET status = $1 
		WHERE image_id = $2 
		AND item_id = $3
		AND list_id = $4
		AND user_id = $5;
	`
	qItem := `
		UPDATE items SET
			item_name = :item_name,
			description = :description,
			address = :address,
			point = :point,
			images_id = :images_id,
			is_visited = :is_visited,
			date_created = :date_created,
			date_updated = :date_updated
		WHERE item_id = :item_id
		AND EXISTS (
			SELECT 1 
			FROM lists 
			WHERE lists.list_id = :list_id 
			AND lists.user_id = :user_id
		);
	`
	qPoint := `
		UPDATE points SET
			location = ST_SetSRID(ST_MakePoint(:lng, :lat), :epsg)
		WHERE point_id = :point_id;
	`

	for _, imgID := range item.ImagesID {
		_, err = tx.ExecContext(ctx, qImages, images.Loaded, imgID, item.ID, item.ListID, item.UserID)
		if err != nil {
			return err
		}
	}
	for _, imgID := range toDelete {
		_, err = tx.ExecContext(ctx, qImages, images.Deleted, imgID, item.ID, item.ListID, item.UserID)
		if err != nil {
			return err
		}
	}
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

func (s *Storer) DeleteItemAdmin(ctx context.Context, itemID string) (err error) {

	qItem := `DELETE FROM items WHERE item_id = $1;`
	qImages := `UPDATE images SET status = $1 WHERE item_id = $2;`
	tx, err := s.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				s.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()
	if err := handleRowsResult(tx.ExecContext(ctx, qItem, itemID)); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, qImages, images.Deleted, itemID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) DeleteItem(ctx context.Context, userID string, listID string, itemID string) error {
	qItem := `
		DELETE FROM items WHERE item_id = $1 
		AND EXISTS (
			SELECT 1
			FROM lists
			WHERE lists.list_id = $2 
			AND lists.user_id = $3
		);
	`
	qImages := `
		UPDATE images SET status = $1 
		WHERE item_id = $2 
		AND list_id = $3
		AND user_id = $4;
		`
	tx, err := s.repo.Beginx()
	if err != nil {
		return err
	}
	err = handleRowsResult(tx.ExecContext(ctx, qItem, itemID, listID, userID))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, qImages, images.Deleted, itemID, listID, userID)
	if err != nil {
		return err
	}
	err = tx.Commit()
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

func populateList(l list.List) StorerList {
	p := pq.StringArray(l.ItemsID)
	list := StorerList{
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

func populateItem(i list.Item) StorerItem {
	im := pq.StringArray(i.ImagesID)
	item := StorerItem{
		ID:          i.ID,
		ListID:      i.ListID,
		UserID:      i.UserID,
		Name:        i.Name,
		Description: i.Description,
		Address:     i.Address,
		PointID:     i.Point.ID,
		ImagesID:    im,
		Visited:     i.Visited,
		DateCreated: i.DateCreated,
		DateUpdated: i.DateUpdated,
	}
	return item
}

type rowItemsByListID struct {
	StorerItem
	StorerPoint
}

func fromRowsToMap(rows *sqlx.Rows) (map[string]*list.Item, error) {
	itemsMap := make(map[string]*list.Item)
	for rows.Next() {
		row := rowItemsByListID{}
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		if _, exists := itemsMap[row.StorerItem.ID]; !exists {
			p := list.Point{
				ID:     row.StorerPoint.ID,
				ItemID: row.StorerItem.ID,
				Lat:    row.Lat,
				Lng:    row.Lng,
			}
			im := []string(row.ImagesID)
			itemsMap[row.StorerItem.ID] = &list.Item{
				ID:          row.StorerItem.ID,
				ListID:      row.ListID,
				UserID:      row.UserID,
				Name:        row.StorerItem.Name,
				Description: row.Description,
				Address:     row.Address,
				Point:       p,
				ImagesID:    im,
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
