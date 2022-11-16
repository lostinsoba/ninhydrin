package postgres

import (
	"context"

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

func (s *Storage) ListPools(ctx context.Context, tagIDs ...string) (pools []*model.Pool, err error) {
	var query = `select id, description, tag_ids from pool where $1 or tag_ids = any($2)`
	rows, err := s.db.QueryContext(ctx, query, len(tagIDs) == 0, pq.Array(tagIDs))
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}

	pools = make([]*model.Pool, 0)
	for rows.Next() {
		var (
			id          string
			description string
			tagIDs      []string
		)
		err = rows.Scan(&id, &description, pq.Array(&tagIDs))
		if err != nil {
			return
		}
		pools = append(pools, &model.Pool{
			ID:          id,
			Description: description,
			TagIDs:      tagIDs,
		})
	}
	return
}

func (s *Storage) UpdatePool(ctx context.Context, pool *model.Pool) error {
	var query = `update pool set description = $1, tag_ids = $2 where id = $3`
	_, err := s.db.ExecContext(ctx, query, pool.Description, pool.TagIDs, pool.ID)
	return err
}
