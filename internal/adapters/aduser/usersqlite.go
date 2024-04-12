package aduser

import (
	"adflow/internal/ads"
	"adflow/internal/app"
	"sync"

	"gorm.io/gorm"
)

type sqliteUsers struct {
	db  *gorm.DB
	cnt int64
	m   sync.Mutex
}

func (r *sqliteUsers) Create(user *ads.User) error {
	r.m.Lock()
	defer r.m.Unlock()

	var temp ads.User

	r.db.Where("Nickname = ?", user.Nickname).Find(&temp)

	if temp == (ads.User{}) {
		user.ID = r.cnt + 1
		r.db.Create(user)
		r.cnt += 1
	} else {
		return ads.ErrBadRequest
	}

	return nil
}

func (r *sqliteUsers) Get(id int64) (*ads.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var user ads.User

	r.db.Where("ID = ?", id).Find(&user)

	if user == (ads.User{}) {
		return nil, ads.ErrBadRequest
	}

	return &user, nil
}

func (r *sqliteUsers) GetByNickname(nickname string) (*ads.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var user ads.User

	r.db.Where("Nickname = ?", nickname).Find(&user)

	if user == (ads.User{}) {
		return nil, ads.ErrBadRequest
	}

	return &user, nil
}

func (r *sqliteUsers) Update(UserID int64, first_name, second_name, phone, email string) (*ads.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var user ads.User

	r.db.Where("ID = ?", UserID).Find(&user)

	if user == (ads.User{}) {
		return nil, ads.ErrBadRequest
	}

	user.FirstName = first_name
	user.SecondName = second_name
	user.Phone = phone
	user.Email = email

	r.db.Save(&user)

	return &user, nil
}

func (r *sqliteUsers) Delete(id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	var user ads.User

	result := r.db.Where("ID = ?", id).First(&user)

	if result.Error != nil {
		return ads.ErrBadRequest
	}

	r.db.Delete(&user)

	return nil
}

func NewSQLiteUsers(db *gorm.DB) app.Users {
	err := db.AutoMigrate(&ads.User{})
	if err != nil {
		panic(err)
	}

	return &sqliteUsers{
		db:  db,
		cnt: 0,
	}
}
