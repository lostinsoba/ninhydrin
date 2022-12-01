package postgres

import (
	"context"
	"database/sql"
	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterTag(ctx context.Context, tagID string) error {
	var query = `insert into tag (id) values ($1) returning id`
	_, err := s.db.ExecContext(ctx, query, tagID)
	return err
}

func (s *Storage) DeregisterTag(ctx context.Context, tagID string) error {
	var query = `delete from tag where id = $1`
	_, err := s.db.ExecContext(ctx, query, tagID)
	return err
}

func (s *Storage) ReadTag(ctx context.Context, tagID string) (tag string, err error) {
	var query = `select id from tag where id = $1`
	err = s.db.QueryRowContext(ctx, query, tagID).Scan(&tag)
	if err == sql.ErrNoRows {
		err = model.ErrNotFound{}
	}
	return
}

func (s *Storage) ListTagIDs(ctx context.Context) (tagIDs []string, err error) {
	var query = `select id from tag`
	rows, err := s.db.QueryContext(ctx, query)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}

	tagIDs = make([]string, 0)
	for rows.Next() {
		var (
			id string
		)
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		tagIDs = append(tagIDs, id)
	}
	return
}
