package postgres

import (
	"context"
	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterWorker(ctx context.Context, worker *model.Worker) error {
	var query = `insert into worker (id, tag_ids) values ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, worker.ID, pq.Array(&worker.TagIDs))
	return err
}

func (s *Storage) DeregisterWorker(ctx context.Context, workerID string) error {
	var query = `delete from worker where id = $1`
	_, err := s.db.ExecContext(ctx, query, workerID)
	return err
}

func (s *Storage) ListWorkers(ctx context.Context) (workers []*model.Worker, err error) {
	var query = `select id, tag_ids from worker`
	rows, err := s.db.QueryContext(ctx, query)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	workers = make([]*model.Worker, 0)
	for rows.Next() {
		var (
			id     string
			tagIDs []string
		)
		err = rows.Scan(&id, pq.Array(&tagIDs))
		if err != nil {
			return nil, err
		}
		workers = append(workers, &model.Worker{
			ID:     id,
			TagIDs: tagIDs,
		})
	}
	return workers, nil
}

func (s *Storage) ListWorkerTagIDs(ctx context.Context, workerID string) (tagIDs []string, err error) {
	var query = `select tag_ids from worker where id = $1`
	row := s.db.QueryRowContext(ctx, query, workerID)
	err = row.Scan(pq.Array(&tagIDs))
	return
}

func (s *Storage) ListWorkerIDsByTagIDs(ctx context.Context, tagIDs ...string) (workerIDs []string, err error) {
	var query = `select id from worker where tag_ids = any($1)`
	rows, err := s.db.QueryContext(ctx, query, pq.Array(tagIDs))
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	workerIDs = make([]string, 0)
	for rows.Next() {
		var (
			id string
		)
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		workerIDs = append(workerIDs, id)
	}
	return workerIDs, nil
}
