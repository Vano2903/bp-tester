package sqliteRepo

import (
	"context"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AttemptRepoerSqlite struct {
	filePath string
	db       *gorm.DB
}

func NewAttemptRepo(filePath string) (repo.AttemptRepoer, error) {
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.Table("attempts").AutoMigrate(&model.Attempt{})
	return &AttemptRepoerSqlite{
		filePath: filePath,
		db:       db.Table("attempts"),
	}, nil
}

func (r *AttemptRepoerSqlite) FindByID(ctx context.Context, id uint) (*model.Attempt, error) {
	a := new(model.Attempt)
	result := r.db.WithContext(ctx).Find(a, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, repo.ErrNotFound
		}
		return nil, result.Error
	}
	return a, nil
}

func (r *AttemptRepoerSqlite) FindByCode(ctx context.Context, code string) (*model.Attempt, error) {
	a := new(model.Attempt)
	result := r.db.WithContext(ctx).Where("code = ?", code).First(a)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, repo.ErrNotFound
		}
		return nil, result.Error
	}
	return a, nil
}

func (r *AttemptRepoerSqlite) InsertOne(ctx context.Context, attempt *model.Attempt) error {
	// a.Attempt.CreatedAt = time.Now()
	result := r.db.WithContext(ctx).Create(attempt)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AttemptRepoerSqlite) UpdateOne(ctx context.Context, attempt *model.Attempt) error {
	result := r.db.WithContext(ctx).Save(attempt)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AttemptRepoerSqlite) DeleteByID(ctx context.Context, code string) error {
	// result, err := r.collection.DeleteOne(ctx, bson.M{
	// 	"code": code,
	// })
	// if err != nil {
	// 	return false, err
	// }
	// return result.DeletedCount > 0, nil
	return nil
}
