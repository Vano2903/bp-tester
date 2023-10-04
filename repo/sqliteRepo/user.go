package sqliteRepo

import (
	"context"
	"fmt"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepoerSqlite struct {
	filePath string
	db       *gorm.DB
}

func NewUserRepo(filePath string) (repo.UserRepoer, error) {
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.Table("users").AutoMigrate(&model.User{})
	return &UserRepoerSqlite{
		filePath: filePath,
		db:       db.Table("users"),
	}, nil
}

func (r *UserRepoerSqlite) FindByID(ctx context.Context, id uint) (*model.User, error) {
	u := new(model.User)
	result := r.db.WithContext(ctx).Find(u, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, repo.ErrNotFound
		}
		return nil, result.Error
	}
	return u, nil
}

func (r *UserRepoerSqlite) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	u := new(model.User)
	result := r.db.WithContext(ctx).Where("username = ?", username).Find(u)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, repo.ErrNotFound
		}
		return nil, result.Error
	}
	return u, nil
}

func (r *UserRepoerSqlite) InsertOne(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		fmt.Println("error inserting user", result.Error.Error())
		if result.Error.Error() == "UNIQUE constraint failed: users.username" {
			return repo.ErrUsernameTaken
		}
		return result.Error
	}
	return nil
}

// will only update password as the username cannot be changed
func (r *UserRepoerSqlite) UpdateOne(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Model(user).Update("password", user.Password).Error
}

func (r *UserRepoerSqlite) DeleteByID(ctx context.Context, code string) error {
	// result, err := r.collection.DeleteOne(ctx, bson.M{
	// 	"code": code,
	// })
	// if err != nil {
	// 	return false, err
	// }
	// return result.DeletedCount > 0, nil
	return nil
}
