package httpgin

import (
	"github.com/gin-gonic/gin"

	"adflow/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.DELETE("/ads/:ad_id", deleteAd(a))           // Метод для удаления объявления
	r.DELETE("/users/:user_id", deleteUser(a))     // Метод для удаления объявления
	r.GET("/ads/:ad_id", getAd(a))                 // Метод для получения объявления
	r.GET("/users/:user_id", getUser(a))           // Метод для получения пользователя
	r.GET("/ads", listAds(a))                      // Метод для получения отфильтр. об.
	r.POST("/users", createUser(a))                // Метод для создания пользователей
	r.POST("/users/login", loginUser(a))           // Метод для логирования пользователей
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.PUT("/users/:user_id", updateUser(a))        // Метод для обновления текста(Text) или заголовка(Title) объявления
}
