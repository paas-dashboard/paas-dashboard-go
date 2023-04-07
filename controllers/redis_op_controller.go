package controllers

import (
	"context"
	"encoding/json"
	"paas-dashboard-go/module"
)

type RedisOpController struct {
	RedisController
}

func (r *RedisOpController) Get() {
	instance := r.GetString(":instance")
	ins, loaded := redisInstanceMap[instance]
	if !loaded {
		r.CustomAbort(500, "instance not found")
		return
	}
	ctx := context.TODO()
	rdb, err := ins.NewClient(ctx)
	if err != nil {
		r.CustomAbort(500, err.Error())
		return
	}
	key := r.GetString(":key")
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		r.CustomAbort(500, err.Error())
		return
	}
	r.Data["json"] = val
	_ = r.ServeJSON()
}

func (r *RedisOpController) Keys() {
	instance := r.GetString(":instance")
	ins, loaded := redisInstanceMap[instance]
	if !loaded {
		r.CustomAbort(500, "instance not found")
		return
	}
	ctx := context.TODO()
	rdb, err := ins.NewClient(ctx)
	if err != nil {
		r.CustomAbort(500, err.Error())
		return
	}
	val, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		r.CustomAbort(500, err.Error())
		return
	}
	r.Data["json"] = val
	_ = r.ServeJSON()
}

func (r *RedisOpController) Set() {
	instance := r.GetString(":instance")
	ins, loaded := redisInstanceMap[instance]
	if !loaded {
		r.CustomAbort(500, "instance not found")
		return
	}
	ctx := context.TODO()
	rdb, err := ins.NewClient(ctx)
	if err != nil {
		r.CustomAbort(500, err.Error())
		return
	}
	i := &module.SetItem{}
	if err = json.Unmarshal(r.Ctx.Input.RequestBody, i); err != nil {
		r.CustomAbort(500, err.Error())
		return
	}
	if _, err = rdb.Set(ctx, i.Key, i.Value, module.EXPIRATION).Result(); err != nil {
		r.CustomAbort(500, err.Error())
		return
	}

	r.Data["json"] = map[string]string{"message": "set key successful"}
	_ = r.ServeJSON()
}

func (r *RedisOpController) Delete() {
	instance := r.GetString(":instance")
	ins, loaded := redisInstanceMap[instance]
	if !loaded {
		r.CustomAbort(500, "instance not found")
		return
	}
	ctx := context.TODO()
	rdb, err := ins.NewClient(ctx)
	if err != nil {
		r.CustomAbort(500, err.Error())
		return
	}

	key := r.GetString(":key")
	if _, err = rdb.Del(ctx, key).Result(); err != nil {
		r.CustomAbort(500, err.Error())
		return
	}

	r.Data["json"] = map[string]string{"message": "delete key successful"}
	_ = r.ServeJSON()
}
