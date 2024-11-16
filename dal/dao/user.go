package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/dal/model"
	"github.com/kackerx/go-mall/logic/do"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (u *UserDao) CreateUser(ctx context.Context, user *do.UserBaseInfo, passwordHash string) (userID int64, err error) {
	userModel := new(model.User)

	if err = util.Copy(userModel, user); err != nil {
		err = errcode.Wrap("CreateUser copy err", err)
		return
	}

	userModel.Password = passwordHash
	if err = DBMaster().WithContext(ctx).Create(userModel).Error; err != nil {
		err = errcode.Wrap("CreateUser db create err", err)
		return
	}

	return userModel.ID, nil
}

func (u *UserDao) FindUserByUserName(ctx context.Context, userName string) (user *do.UserBaseInfo, err error) {
	userPO := new(model.User)
	if err = DBMaster().WithContext(ctx).
		Where("user_name = ?", userName).
		First(&userPO).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	user = new(do.UserBaseInfo)
	util.Copy(user, userPO)
	return user, nil
}
