package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterPool(ctx context.Context, pool *model.Pool) error {
	var query = `insert into pool (id, description, tag_ids) values ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, pool.ID, pool.Description, pq.Array(pool.TagIDs))
	return err
}

func (s *Storage) DeregisterPool(ctx context.Context, poolID string) error {
	var query = `delete from pool where id = $1`
	_, err := s.db.ExecContext(ctx, query, poolID)
	return err
}

func (s *Storage) ReadPool(ctx context.Context, poolID string) (pool *model.Pool, err error) {
	var query = `select id, description, tag_ids from pool where id = $1`
	var (
		id          string
		description string
		tagIDs      []string
	)
	err = s.db.QueryRowContext(ctx, query, poolID).Scan(&id, &description, pq.Array(&tagIDs))
	switch err {
	case nil:
		return &model.Pool{
			ID:          id,
			Description: description,
			TagIDs:      tagIDs,
		}, nil
	case sql.ErrNoRows:
		return nil, model.ErrNotFound{}
	default:
		return nil, err
	}
}

func (s *Storage) ListPoolIDs(ctx context.Context, tagIDs ...string) (poolIDs []string, err error) {
	var query = `select id from pool where $1 or tag_ids = any($2)`
	rows, err := s.db.QueryContext(ctx, query, len(tagIDs) == 0, pq.Array(tagIDs))
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}

	poolIDs = make([]string, 0)
	for rows.Next() {
		var (
			id string
		)
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		poolIDs = append(poolIDs, id)
	}
	return
}

func (s *Storage) UpdatePool(ctx context.Context, pool *model.Pool) error {
	var query = `update pool set description = $1, tag_ids = $2 where id = $3`
	_, err := s.db.ExecContext(ctx, query, pool.Description, pool.TagIDs, pool.ID)
	return err
}
