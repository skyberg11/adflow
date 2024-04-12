package httpgin

import (
	"github.com/gin-gonic/gin"

	"adflow/internal/ads"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type deleteRequest struct {
	UserID int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type createUserRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Nickname   string `json:"nickname"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type updateUserRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	UserID     int64  `json:"user_id"`
}

type userResponse struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	UserID     int64  `json:"user_id"`
}

func DeleteSuccessResponse() *gin.H {
	return &gin.H{
		"data":  "delete success",
		"error": nil,
	}
}

func ListAdsSuccessResponse(ad []*ads.Ad) *gin.H {
	var copy []adResponse

	for _, v := range ad {
		copy = append(copy, adResponse{v.ID, v.Title, v.Text, v.AuthorID, v.Published, v.CreationTime, v.UpdateTime})
	}

	return &gin.H{
		"data":  copy,
		"error": nil,
	}
}

func UserSuccessResponse(user *ads.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			Email:      user.Email,
			Phone:      user.Phone,
			UserID:     user.ID,
		},
		"error": nil,
	}
}

func LoginSuccessResponse(token string) *gin.H {
	return &gin.H{
		"token": token,
		"error": nil,
	}
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationTime: ad.CreationTime,
			UpdateTime:   ad.UpdateTime,
		},
		"error": nil,
	}
}

func ErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
