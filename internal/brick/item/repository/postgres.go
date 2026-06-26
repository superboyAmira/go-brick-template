package repository

import (
	"context"

	"github.com/go-brick-template/go-brick-template/internal/brick/item/model"
	"github.com/go-brick-template/go-brick-template/module/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgres(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context) ([]model.Item, error) {
	q := psql.Select("id", "title", "created_at", "updated_at").
		From("items").
		OrderBy("created_at DESC")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Read(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Item
	for rows.Next() {
		var it model.Item
		if err := rows.Scan(&it.ID, &it.Title, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	q := psql.Select("id", "title", "created_at", "updated_at").
		From("items").
		Where(sq.Eq{"id": id})
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	var it model.Item
	err = r.db.Read(ctx).QueryRow(ctx, sql, args...).Scan(&it.ID, &it.Title, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *PostgresRepository) Create(ctx context.Context, title string) (*model.Item, error) {
	q := psql.Insert("items").
		Columns("title").
		Values(title).
		Suffix("RETURNING id, title, created_at, updated_at")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	var it model.Item
	err = r.db.Write(ctx).QueryRow(ctx, sql, args...).Scan(&it.ID, &it.Title, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := psql.Delete("items").Where(sq.Eq{"id": id})
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}
	ct, err := r.db.Write(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
