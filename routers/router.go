package routers

import (
	"fmt"
	"grpc-json-server/config"
	"grpc-json-server/controllers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

var cfg config.Configuration

func init() {
	err := cleanenv.ReadEnv(&cfg)
	fmt.Printf("%v", cfg)
	if err != nil {
		log.Fatalf("Unable to load configuration")
	}
}

func Start() {
	var r = gin.Default()
	r.POST("/sample", controllers.Sample)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/about", controllers.About)
	r.PUT("/update", controllers.Update)
	r.DELETE("/delete", controllers.Delete)
	r.Run()
}
