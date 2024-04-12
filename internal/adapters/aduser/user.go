package aduser

import (
	"adflow/internal/ads"
	"adflow/internal/app"
	"errors"
	"sync"
)

type localUsers struct {
	users map[int64]*ads.User
	cnt   int64
	m     sync.Mutex
}

func (r *localUsers) Create(ad *ads.User) error {
	r.m.Lock()
	defer r.m.Unlock()

	for _, existingUser := range r.users {
		if existingUser.Nickname == ad.Nickname {
			return errors.New("пользователь с таким никнеймом уже существует")
		}
	}

	ad.ID = r.cnt
	r.users[r.cnt] = ad
	r.cnt += 1
	return nil
}

func (r *localUsers) Get(id int64) (*ads.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.users[id]
	if !ok {
		return nil, ads.ErrBadRequest
	}
	return r.users[id], nil
}

func (r *localUsers) GetByNickname(nickname string) (*ads.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	// _, ok := r.users[id]
	// if !ok {
	// 	return nil, ads.ErrBadRequest
	// }
	return nil, errors.New("unimplemented")
}

func (r *localUsers) Update(UserID int64, first_name, second_name, phone, email string) (*ads.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.users[UserID]; !ok {
		return nil, ads.ErrBadRequest
	}

	r.users[UserID].FirstName = first_name
	r.users[UserID].SecondName = second_name
	r.users[UserID].Phone = phone
	r.users[UserID].Email = email
	return r.users[UserID], nil
}

func (r *localUsers) Delete(id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, err := r.users[id]; !err {
		return ads.ErrBadRequest
	}

	delete(r.users, id)

	return nil
}

func New() app.Users {
	return &localUsers{
		users: make(map[int64]*ads.User),
		cnt:   0,
	}
}
