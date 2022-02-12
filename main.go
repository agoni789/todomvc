package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todomvc-app-template-golang/db"
	"todomvc-app-template-golang/handler"
)

func main() {
	engine := gin.Default()
	initRouter(engine)
	server := initServer(engine)
	db.InitDB()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initRouter(engine *gin.Engine) {
	r := engine.Group("api")
	{
		r.POST("add", handler.Add)
		r.POST("del", handler.Del)
		r.POST("update", handler.Update)
		r.POST("find", handler.Find)
	}

}

func initServer(engine *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           "localhost:8080",
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
