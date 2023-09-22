package sqliteRepo

import (
	"context"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ExecitionRepoerSqlite struct {
	filePath string
	db       *gorm.DB
}

func NewExecutionRepo(filePath string) (repo.ExecutionRepoer, error) {
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.Table("executions").AutoMigrate(&model.Execution{})
	return &ExecitionRepoerSqlite{
		filePath: filePath,
		db:       db.Table("executions"),
	}, nil
}

func (r *ExecitionRepoerSqlite) FindByID(ctx context.Context, id uint) (*model.Execution, error) {
	e := new(model.Execution)
	result := r.db.WithContext(ctx).Find(e, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, repo.ErrNotFound
		}
		return nil, result.Error
	}
	return e, nil
}

func (r *ExecitionRepoerSqlite) FindByAttemptID(ctx context.Context, attemptID uint) ([]*model.Execution, error) {
	var executions []*model.Execution
	result := r.db.WithContext(ctx).Where("attempt_id = ?", attemptID).Find(&executions)
	return executions, result.Error
}

func (r *ExecitionRepoerSqlite) InsertOne(ctx context.Context, execution *model.Execution) error {
	e := new(model.Execution)
	result := r.db.WithContext(ctx).Create(e)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ExecitionRepoerSqlite) UpdateOne(ctx context.Context, execution *model.Execution) error {
	return r.db.WithContext(ctx).Save(execution).Error
}

func (r *ExecitionRepoerSqlite) DeleteByID(ctx context.Context, code string) error {
	// result, err := r.collection.DeleteOne(ctx, bson.M{
	// 	"code": code,
	// })
	// if err != nil {
	// 	return false, err
	// }
	// return result.DeletedCount > 0, nil
	return nil
}
