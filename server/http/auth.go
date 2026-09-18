package http

import (
	"net/http"
	"time"

	httppkg "github.com/fatedier/frp/pkg/util/http"
	"github.com/fatedier/frp/pkg/util/util"
)

type loginRequest struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type statusResponse struct {
	AuthRequired bool   `json:"authRequired"`
	LoggedIn     bool   `json:"loggedIn"`
	User         string `json:"user,omitempty"`
}

func (c *Controller) Login(ctx *httppkg.Context) (any, error) {
	cfgUser := c.serverCfg.WebServer.User
	cfgPassword := c.serverCfg.WebServer.Password

	if cfgUser == "" && cfgPassword == "" {
		return loginResponse{}, nil
	}

	var req loginRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, httppkg.NewError(http.StatusBadRequest, "invalid request body")
	}

	if !util.ConstantTimeEqString(req.User, cfgUser) || !util.ConstantTimeEqString(req.Password, cfgPassword) {
		time.Sleep(200 * time.Millisecond)
		return nil, httppkg.NewError(http.StatusUnauthorized, "invalid credentials")
	}

	expiry := 24 * time.Hour
	if req.RememberMe {
		expiry = 30 * 24 * time.Hour
	}

	token, expiresAt := c.sessionMgr.CreateWithExpiry(expiry)

	http.SetCookie(ctx.Resp, &http.Cookie{
		Name:     "frp-auth-token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})

	return loginResponse{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (c *Controller) Logout(ctx *httppkg.Context) (any, error) {
	cookie, err := ctx.Req.Cookie("frp-auth-token")
	if err == nil && cookie.Value != "" {
		c.sessionMgr.Delete(cookie.Value)
	}

	http.SetCookie(ctx.Resp, &http.Cookie{
		Name:     "frp-auth-token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return httppkg.GeneralResponse{Code: 200, Msg: "success"}, nil
}

func (c *Controller) LoginStatus(ctx *httppkg.Context) (any, error) {
	cfgUser := c.serverCfg.WebServer.User
	cfgPassword := c.serverCfg.WebServer.Password

	if cfgUser == "" && cfgPassword == "" {
		return statusResponse{AuthRequired: false, LoggedIn: false}, nil
	}

	cookie, err := ctx.Req.Cookie("frp-auth-token")
	if err != nil {
		return nil, httppkg.NewError(http.StatusUnauthorized, "not logged in")
	}

	if !c.sessionMgr.Validate(cookie.Value) {
		return nil, httppkg.NewError(http.StatusUnauthorized, "session expired")
	}

	return statusResponse{
		AuthRequired: true,
		LoggedIn:     true,
		User:         cfgUser,
	}, nil
}
