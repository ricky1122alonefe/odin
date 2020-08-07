package pikaRepo

import (
	"context"
	"encoding/json"
	"odin/app/common"
	"odin/app/entity"
	"sync"

	redisdao "github.com/tal-tech/xredis"
)

type UserRepo struct {
}

var once sync.Once

var repo *UserRepo

func NewUserRepo() *UserRepo {
	if repo != nil {
		return repo
	}
	once.Do(func() {
		repo = new(UserRepo)
	})
	return repo
}

const PIKA_USER_PREFIX = "micro_pika_user_%v"

func (this *UserRepo) AddUser(ctx context.Context, user *entity.User) (id int, err error) {
	common.CurrentUserID += 1
	user.Id = common.CurrentUserID
	bytes, _ := json.Marshal(user)
	xesRedis := redisdao.NewSimpleXesRedis(ctx, "pika")
	_, err = xesRedis.Set(PIKA_USER_PREFIX, []interface{}{user.Id}, string(bytes), 0)
	return
}

func (this *UserRepo) GetUserInfo(ctx context.Context, id int) (user *entity.User, err error) {
	xesRedis := redisdao.NewSimpleXesRedis(ctx, "pika")
	ret, err := xesRedis.Get(PIKA_USER_PREFIX, []interface{}{id})
	if err != nil {
		return
	}
	user = new(entity.User)
	json.Unmarshal([]byte(ret), user)
	return
}

func (this *UserRepo) UpdateUser(ctx context.Context, user *entity.User) (err error) {
	bytes, _ := json.Marshal(user)
	xesRedis := redisdao.NewSimpleXesRedis(ctx, "pika")
	_, err = xesRedis.Set(PIKA_USER_PREFIX, []interface{}{user.Id}, string(bytes), 0)
	return
}
