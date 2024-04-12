package httpgin

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"adflow/internal/ads"
	"adflow/internal/app"
	auth "adflow/internal/app/auth"
)

// Метод для получения объявления
func getAd(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.GetAd(c, int64(adID))

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для получения пользователя
func getUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {

		UserID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, err := a.GetUser(c, int64(UserID))

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// Метод для создания пользователя
func createUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		user, err := a.CreateUser(c, reqBody.FirstName, reqBody.SecondName, reqBody.Nickname, reqBody.Password,
			reqBody.Email, reqBody.Phone)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// Метод для создания пользователя
func loginUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		token, err := a.LoginUser(c, reqBody.Nickname, reqBody.Password)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, LoginSuccessResponse(token))
	}
}

// Метод для получения фильтрованных
func listAds(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		filter := ads.Filter{
			Published:    nil,
			AuthorID:     nil,
			TitlePrefix:  nil,
			CreationTime: nil,
		}
		var err error

		published := c.Query("published")
		if published != "" {
			tmp, err := strconv.ParseInt(published, 10, 64)
			filter.Published = (tmp != 0)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
		}

		authorID := c.Query("author")
		if authorID != "" {
			filter.AuthorID, err = strconv.ParseInt(authorID, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
		}

		filter.TitlePrefix = c.Query("title")

		creationTime := c.Query("creation")
		if creationTime != "" {
			filter.CreationTime, err = time.Parse(time.Now().UTC().String(), creationTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
		}

		ad, err := a.ListAds(c, filter)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, ListAdsSuccessResponse(ad))
	}
}

// Метод для создания объявления (ad)
func createAd(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		code, msg := auth.ValidateToken(bearerToken, reqBody.UserID)

		if code != 200 {
			c.JSON(code, gin.H{
				"message": msg,
			})
			return
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		code, msg := auth.ValidateToken(bearerToken, reqBody.UserID)

		if code != 200 {
			c.JSON(code, gin.H{
				"message": msg,
			})
			return
		}

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		code, msg := auth.ValidateToken(bearerToken, reqBody.UserID)

		if code != 200 {
			c.JSON(code, gin.H{
				"message": msg,
			})
			return
		}

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для удаления объявления
func deleteAd(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody deleteRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		code, msg := auth.ValidateToken(bearerToken, reqBody.UserID)

		if code != 200 {
			c.JSON(code, gin.H{
				"message": msg,
			})
			return
		}

		err = a.DeleteAd(c, int64(adID), reqBody.UserID)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, DeleteSuccessResponse())
	}
}

// Метод для удаления пользователя
func deleteUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody deleteRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		code, msg := auth.ValidateToken(bearerToken, reqBody.UserID)

		if code != 200 {
			c.JSON(code, gin.H{
				"message": msg,
			})
			return
		}

		err = a.DeleteUser(c, int64(id), reqBody.UserID)

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, DeleteSuccessResponse())
	}
}

// Метод для обновления пользователя
func updateUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		code, msg := auth.ValidateToken(bearerToken, int64(userID))

		if code != 200 {
			c.JSON(code, gin.H{
				"message": msg,
			})
			return
		}

		user, err := a.UpdateUser(c, reqBody.FirstName, reqBody.SecondName, reqBody.Email, reqBody.Phone, int64(userID))

		if err != nil {
			var status int
			if errors.Is(err, ads.ErrAccessDenied) {
				status = http.StatusForbidden
			} else if errors.Is(err, ads.ErrBadRequest) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}
