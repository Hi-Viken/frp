package configmanager

import (
	"encoding/json"
	"net/http"
	"time"

	httppkg "github.com/fatedier/frp/pkg/util/http"
)

type RestartFunc func()

type Controller struct {
	manager *Manager
	restart RestartFunc
}

func NewController(manager *Manager, restart RestartFunc) *Controller {
	return &Controller{manager: manager, restart: restart}
}

func (c *Controller) GetConfig(ctx *httppkg.Context) (any, error) {
	cfg := c.manager.Get()
	return APIResponse{Success: true, Data: cfg}, nil
}

func (c *Controller) UpdateConfig(ctx *httppkg.Context) (any, error) {
	var cfg FrpsConfig
	if err := json.NewDecoder(ctx.Req.Body).Decode(&cfg); err != nil {
		return APIResponse{Success: false, Error: "invalid JSON: " + err.Error()}, nil
	}
	if err := c.manager.Update(cfg); err != nil {
		return APIResponse{Success: false, Error: "save failed: " + err.Error()}, nil
	}
	return APIResponse{Success: true, Data: cfg}, nil
}

func (c *Controller) ReloadConfig(ctx *httppkg.Context) (any, error) {
	if err := c.manager.Load(); err != nil {
		return APIResponse{Success: false, Error: "reload failed: " + err.Error()}, nil
	}
	cfg := c.manager.Get()
	return APIResponse{Success: true, Data: cfg}, nil
}

func (c *Controller) Restart(ctx *httppkg.Context) (any, error) {
	if c.restart == nil {
		return APIResponse{Success: false, Error: "restart not available"}, nil
	}
	go func() {
		time.Sleep(300 * time.Millisecond)
		c.restart()
	}()
	return APIResponse{Success: true, Data: map[string]string{"message": "restarting"}}, nil
}

func MakeHandler(handler func(ctx *httppkg.Context) (any, error)) http.HandlerFunc {
	return httppkg.MakeHTTPHandlerFunc(handler)
}
