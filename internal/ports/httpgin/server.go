package httpgin

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"adflow/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

func MyLogger(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				f(c, err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func handlePanic(c *gin.Context, err interface{}) {
	log.Printf("Panic occurred: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	})
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}

	s.app.Use(gin.LoggerWithFormatter(MyLogger))

	s.app.Use(Recovery(handlePanic))

	api := s.app.Group("/api/v1")
	AppRouter(api, a)
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
