package app

import (
	"adflow/internal/ads"
	"context"

	auth "adflow/internal/app/auth"

	validator "github.com/skyberg11/args-validator"
)

type Repository interface {
	Create(ad *ads.Ad) error
	Get(id int64) (*ads.Ad, error)
	Update(id int64, title, text string) (*ads.Ad, error)
	UpdateStatus(id int64, published bool) (*ads.Ad, error)
	GetAllAds() ([]*ads.Ad, error)
	GetAds(filter ads.Filter) ([]*ads.Ad, error)
	DeleteAd(id int64) error
}

type Users interface {
	Create(ad *ads.User) error
	Get(id int64) (*ads.User, error)
	GetByNickname(nickname string) (*ads.User, error)
	Update(id int64, first_name, second_name, email, phone string) (*ads.User, error)
	Delete(id int64) error
}

type App interface {
	ListAds(ctx context.Context, filter ads.Filter) ([]*ads.Ad, error)
	CreateAd(ctx context.Context, title, text string, authorID int64) (*ads.Ad, error)
	GetAd(ctx context.Context, id int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, id, UserID int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, id, UserID int64, title, text string) (*ads.Ad, error)
	DeleteAd(ctx context.Context, id int64, UserID int64) error

	CreateUser(ctx context.Context, first_name, second_name, nickname, password, email, phone string) (*ads.User, error)
	LoginUser(ctx context.Context, nickname, password string) (string, error)
	UpdateUser(ctx context.Context, first_name, second_name, email, phone string, UserID int64) (*ads.User, error)
	GetUser(ctx context.Context, UserID int64) (*ads.User, error)
	DeleteUser(ctx context.Context, id int64, UserID int64) error
}

type localApp struct {
	repo  Repository
	users Users
}

func (a *localApp) DeleteAd(ctx context.Context, id int64, UserID int64) error {
	if _, err := a.users.Get(UserID); err != nil {
		return err
	}

	ad, err := a.repo.Get(id)

	if err != nil {
		return err
	}

	if ad.AuthorID != UserID {
		return ads.ErrAccessDenied
	}

	err = a.repo.DeleteAd(id)

	if err != nil {
		return err
	}

	return nil
}

func (a *localApp) ListAds(ctx context.Context, filter ads.Filter) ([]*ads.Ad, error) {
	ads, err := a.repo.GetAds(filter)

	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (a *localApp) CreateAd(ctx context.Context, title, text string, authorID int64) (*ads.Ad, error) {
	if _, err := a.users.Get(authorID); err != nil {
		return nil, err
	}

	ad := &ads.Ad{
		Title:     title,
		Text:      text,
		AuthorID:  authorID,
		Published: false,
	}

	if err := validator.Validate(*ad); err != nil {
		return nil, ads.ErrBadRequest
	}

	err := a.repo.Create(ad)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a *localApp) GetAd(ctx context.Context, id int64) (*ads.Ad, error) {
	return a.repo.Get(id)
}

func (a *localApp) ChangeAdStatus(ctx context.Context, id, UserID int64, published bool) (*ads.Ad, error) {
	if _, err := a.users.Get(UserID); err != nil {
		return nil, err
	}

	ad, err := a.repo.Get(id)

	if err != nil {
		return nil, err
	}

	if ad.AuthorID != UserID {
		return nil, ads.ErrAccessDenied
	}

	ad, err = a.repo.UpdateStatus(id, published)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (a *localApp) UpdateAd(ctx context.Context, id, UserID int64, title, text string) (*ads.Ad, error) {
	if _, err := a.users.Get(UserID); err != nil {
		return nil, err
	}

	ad, err := a.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if ad.AuthorID != UserID {
		return nil, ads.ErrAccessDenied
	}
	prevTitle, prevText := ad.Title, ad.Text

	ad, err = a.repo.Update(id, title, text)

	if err := validator.Validate(*ad); err != nil {
		if _, err = a.repo.Update(id, prevTitle, prevText); err != nil {
			panic("something went wrong")
		}
		return nil, ads.ErrBadRequest
	}

	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a *localApp) CreateUser(ctx context.Context, first_name, second_name, nickname, password, email, phone string) (*ads.User, error) {
	user := &ads.User{
		FirstName:  first_name,
		SecondName: second_name,
		Nickname:   nickname,
		Password:   password,
		Email:      email,
		Phone:      phone,
	}

	if err := validator.Validate(*user); err != nil {
		return nil, ads.ErrBadRequest
	}

	err := a.users.Create(user)
	if err != nil {
		return nil, ads.ErrBadRequest
	}

	return user, nil
}

func (a *localApp) LoginUser(ctx context.Context, nickname, password string) (string, error) {
	user, err := a.users.GetByNickname(nickname)
	if err != nil {
		return "error", ads.ErrBadRequest
	}

	if user.Password != password {
		return "error", ads.ErrAccessDenied
	}

	return auth.GenerateJWT(user.ID)
}

func (a *localApp) UpdateUser(ctx context.Context, first_name, second_name, email, phone string, UserID int64) (*ads.User, error) {
	user, err := a.users.Get(UserID)
	if err != nil {
		return nil, err
	}

	prevFirst, prevSecond, prevEmail, prevPhone := user.FirstName, user.SecondName, user.Email, user.Phone

	user, err = a.users.Update(UserID, first_name, second_name, email, phone)

	if err != nil {
		return nil, ads.ErrBadRequest
	}

	if err := validator.Validate(*user); err != nil {
		if _, err = a.users.Update(UserID, prevFirst, prevSecond, prevEmail, prevPhone); err != nil {
			panic("something went wrong")
		}
		return nil, ads.ErrBadRequest
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *localApp) GetUser(ctx context.Context, UserID int64) (*ads.User, error) {
	return a.users.Get(UserID)
}

func (a *localApp) DeleteUser(ctx context.Context, id int64, UserID int64) error {
	if id != UserID {
		return ads.ErrAccessDenied
	}

	return a.users.Delete(id)
}

func NewApp(repo Repository, users Users) App {
	return &localApp{
		repo:  repo,
		users: users,
	}
}
