// internal/app/repository/default.go
package repository

import (
	"go-schedule/internal/model"
	"go-schedule/internal/pkg/database/cache"
	"go-schedule/internal/pkg/database/mysql"
	"go-schedule/internal/pkg/fetch"

	"go.uber.org/zap"
)

type repository struct {
	cache *cache.Clients
	db    *mysql.Clients
	fetch *fetch.Fetch
	log   *zap.Logger
	lb    map[string]string
}

func New(
	cache *cache.Clients,
	db *mysql.Clients,
	fetch *fetch.Fetch,
	log *zap.Logger,
	lb map[string]string,
) Repository {
	return &repository{

		cache,
		db,
		fetch,
		log,
		lb,
	}
}

type Repository interface {
	FindByID(id int64) (model.UsersMobile, error)
}

func (r *repository) FindByID(uid int64) (model.UsersMobile, error) {
	db := r.db.Client("default.slave")

	result := model.UsersMobile{}
	query := "SELECT `uid`,`mid`,`region`,`encrypt`,`create_time` FROM `users_mobile` WHERE `uid` = ? LIMIT 1"
	rows, err := db.Query(query, uid)

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
