package main

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"paas-dashboard-go/checker"
	"paas-dashboard-go/controllers"
	"paas-dashboard-go/log"
	"path/filepath"
)

func main() {
	log.InitLogger()
	checker.Start()

	root := filepath.Join(".", "static")
	web.SetStaticPath("/", root)

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
	web.Router("/kafka", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/kafka/**", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/kubernetes", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/kubernetes/**", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/mongo", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/mongo/**", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/pulsar", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/pulsar/**", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/redis", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/redis/**", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/rocketmq", &controllers.MainController{}, "get:RedirectToIndex")
	web.Router("/rocketmq/**", &controllers.MainController{}, "get:RedirectToIndex")

	// Run the app
	web.BConfig.Listen.HTTPPort = 11111
	web.BConfig.CopyRequestBody = true
	web.Run()
}
