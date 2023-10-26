package db

import (
	"context"
	"time"

	"github.com/fmiskovic/go-starter/data"
	"github.com/uptrace/bun"
)

// UserDao implements Dao interface
type UserDao struct {
	db *bun.DB
}

func NewUserDao() *UserDao {
	return &UserDao{Bun}
}

// Get returns user queried by id
func (dao *UserDao) Get(id int64) (*data.User, error) {
	var user = &data.User{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout))
	defer cancel()

	err := dao.db.NewSelect().Model(user).Where("? = ?", bun.Ident("id"), id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Save persists new user or updates existing one
func (dao *UserDao) Save(u *data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout))
	defer cancel()

	if u.ID > 0 {
		// update
		_, err := dao.db.NewUpdate().Model(u).Where("id = ?", u.ID).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	// save
	_, err := dao.db.NewInsert().Model(u).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Delete user by id
func (dao *UserDao) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout))
	defer cancel()

	dao.db.NewDelete().Where("id = ?", id).Exec(ctx)
	return nil
}

// GetPage of users
func (dao *UserDao) GetPage(p Pageable) (Page[data.User], error) {
	// TODO
	return Page[data.User]{}, nil
}
