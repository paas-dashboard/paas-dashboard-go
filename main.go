package main

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"paas-dashboard-go/controllers"
	"path/filepath"
)

func main() {
	root := filepath.Join(".", "static")
	web.SetStaticPath("/static", root)

	// CORS settings
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register controllers
	controllers.RegisterKafkaControllers()
	controllers.RegisterKubernetesControllers()
	controllers.RegisterPulsarControllers()
	controllers.RegisterRedisControllers()

	web.Router("/", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/**", &controllers.MainController{}, "get:RedirectToIndex")

	// Run the app
	web.BConfig.Listen.HTTPPort = 11111
	web.BConfig.CopyRequestBody = true
	web.Run()
}
