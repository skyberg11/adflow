package app

import (
	"adflow/internal/ads"
	"adflow/internal/app"
	service "adflow/internal/ports/grpc/service"
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	validator "github.com/skyberg11/args-validator"
)

type AdService struct {
	repo  app.Repository
	users app.Users
	service.UnimplementedAdServiceServer
}

func (a *AdService) DeleteAd(ctx context.Context, req *service.DeleteAdRequest) (*empty.Empty, error) {
	if _, err := a.users.Get(req.AuthorId); err != nil {
		return nil, err
	}

	if req.AuthorId != req.AdId {
		return nil, ads.ErrAccessDenied
	}

	err := a.repo.DeleteAd(req.AdId)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (a *AdService) ListAds(ctx context.Context, filter *service.Filter) (*service.ListAdResponse, error) {
	f := ads.Filter{
		Published:    nil,
		AuthorID:     nil,
		TitlePrefix:  nil,
		CreationTime: nil,
	}
	var err error

	published := filter.Published
	if published != "" {
		f.Published, err = strconv.ParseBool(published)
		if err != nil {
			return nil, err
		}
	}

	authorID := filter.AuthorId
	if authorID != "" {
		f.AuthorID, err = strconv.ParseInt(authorID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	f.TitlePrefix = filter.Prefix

	creationTime := filter.CreationTime
	if creationTime != "" {
		f.CreationTime, err = time.Parse(time.Now().UTC().String(), creationTime)
		if err != nil {
			return nil, err
		}
	}

	ads, err := a.repo.GetAds(f)

	if err != nil {
		return nil, err
	}
	var List []*service.AdResponse

	for _, ad := range ads {
		List = append(List, &service.AdResponse{
			Id:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationTime.String(),
			UpdateDate:   ad.UpdateTime.String(),
		})
	}

	return &service.ListAdResponse{
		List: List,
	}, nil
}

func (a *AdService) CreateAd(ctx context.Context, req *service.CreateAdRequest) (*service.AdResponse, error) {
	if _, err := a.users.Get(req.UserId); err != nil {
		return nil, err
	}

	ad := &ads.Ad{
		Title:     req.Title,
		Text:      req.Text,
		AuthorID:  req.UserId,
		Published: false,
	}

	if err := validator.Validate(*ad); err != nil {
		return nil, ads.ErrBadRequest
	}

	err := a.repo.Create(ad)
	if err != nil {
		return nil, err
	}

	return &service.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationTime.String(),
		UpdateDate:   ad.UpdateTime.String(),
	}, nil
}

func (a *AdService) GetAd(ctx context.Context, req *service.GetAdRequest) (*service.AdResponse, error) {
	ad, err := a.repo.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return &service.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationTime.String(),
		UpdateDate:   ad.UpdateTime.String(),
	}, nil
}

func (a *AdService) ChangeAdStatus(ctx context.Context, req *service.ChangeAdStatusRequest) (*service.AdResponse, error) {
	if _, err := a.users.Get(req.UserId); err != nil {
		return nil, err
	}

	ad, err := a.repo.Get(req.AdId)

	if err != nil {
		return nil, err
	}

	if ad.AuthorID != req.UserId {
		return nil, ads.ErrAccessDenied
	}

	ad, err = a.repo.UpdateStatus(req.AdId, req.Published)
	if err != nil {
		return nil, err
	}
	return &service.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationTime.String(),
		UpdateDate:   ad.UpdateTime.String(),
	}, nil
}

func (a *AdService) UpdateAd(ctx context.Context, req *service.UpdateAdRequest) (*service.AdResponse, error) {
	if _, err := a.users.Get(req.UserId); err != nil {
		return nil, err
	}

	ad, err := a.repo.Get(req.AdId)
	if err != nil {
		return nil, err
	}

	if ad.AuthorID != req.UserId {
		return nil, ads.ErrAccessDenied
	}
	prevTitle, prevText := ad.Title, ad.Text

	ad, err = a.repo.Update(req.AdId, req.Title, req.Text)

	if err := validator.Validate(*ad); err != nil {
		if _, err = a.repo.Update(req.AdId, prevTitle, prevText); err != nil {
			panic("something went wrong")
		}
		return nil, ads.ErrBadRequest
	}

	if err != nil {
		return nil, err
	}

	return &service.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationTime.String(),
		UpdateDate:   ad.UpdateTime.String(),
	}, nil
}

func (a *AdService) CreateUser(ctx context.Context, req *service.CreateUserRequest) (*service.UserResponse, error) {
	user := &ads.User{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Nickname:   req.Nickname,
		Password:   req.Password,
		Email:      req.Email,
		Phone:      req.Phone,
	}

	if err := validator.Validate(*user); err != nil {
		return nil, ads.ErrBadRequest
	}

	err := a.users.Create(user)
	if err != nil {
		return nil, err
	}

	return &service.UserResponse{
		Id:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Email:      user.Email,
		Phone:      user.Phone,
	}, nil
}

func (a *AdService) GetUser(ctx context.Context, req *service.GetUserRequest) (*service.UserResponse, error) {
	user, err := a.users.Get(req.Id)

	if err != nil {
		return nil, err
	}

	return &service.UserResponse{
		Id:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Email:      user.Email,
		Phone:      user.Phone,
	}, nil
}

func (a *AdService) DeleteUser(ctx context.Context, req *service.DeleteUserRequest) (*empty.Empty, error) {
	if req.Id != req.UserId {
		return nil, ads.ErrAccessDenied
	}

	err := a.users.Delete(req.Id)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func NewAdService(repo app.Repository, users app.Users) service.AdServiceServer {
	return &AdService{
		repo:  repo,
		users: users,
	}
}
