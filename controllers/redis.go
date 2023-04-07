package controllers

import "github.com/beego/beego/v2/server/web"

type RedisController struct {
	web.Controller
}

func RegisterRedisControllers() {
	web.Router("/api/redis/instances", &RedisInstanceController{}, "get:Get")

	redisOpController := &RedisOpController{}
	web.Router("/api/redis/instances/:instance/keys", redisOpController, "get:Keys")
	web.Router("/api/redis/instances/:instance/keys", redisOpController, "put:Set")
	web.Router("/api/redis/instances/:instance/keys/:key", redisOpController, "get:Get")
	web.Router("/api/redis/instances/:instance/keys/:key", redisOpController, "delete:Delete")
}
