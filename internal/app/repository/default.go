// internal/app/repository/default.go
package repository

import (
	"context"

	"go-schedule/internal/model"
	"go-schedule/internal/pkg/database/cache"
	"go-schedule/internal/pkg/database/mysql"
	"go-schedule/internal/pkg/fetch"

	"go.uber.org/zap"
)

type repository struct {
	api   map[string]string
	cache *cache.Clients
	db    *mysql.Clients
	fetch *fetch.Fetch
	log   *zap.Logger
}

func NewRepository(
	api map[string]string,
	cache *cache.Clients,
	db *mysql.Clients,
	fetch *fetch.Fetch,
	log *zap.Logger,
) Repository {
	return &repository{
		api,
		cache,
		db,
		fetch,
		log,
	}
}

type Repository interface {
	FindByID(ctx context.Context, id int64) (model.UsersMobile, error)
}

func (r *repository) FindByID(ctx context.Context, uid int64) (model.UsersMobile, error) {
	db := r.db.Client("default.slave")

	result := model.UsersMobile{}
	query := "SELECT `uid`,`mid`,`region`,`encrypt`,`create_time` FROM `users_mobile` WHERE `uid` = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, query, uid)

	if err != nil {
		r.log.Error("FindByID", []zap.Field{
			zap.String("query", query),
			zap.Int64("uid", uid),
			zap.String("error", err.Error()),
		}...)
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&result.Uid,
			&result.Mid,
			&result.Region,
			&result.Encrypt,
			&result.CreateTime,
		)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
