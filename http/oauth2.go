package http

import (
	"encoding/json"
	"fmt"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/go-resty/resty/v2"
	"log"
	"net"
	"net/http"
	"time"
)

type oauth2CallBackBody struct {
	Code string `json:"code"`
}

type AuthSuccess struct {
	Token     string `json:"access_token"`
	Refresh   string `json:"refresh_token"`
	ExpiresIn string `json:"expires_in"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
}
type OauthUserInfoChuanDI struct {
	ClientId string `json:"client_id"`
	OpenId   string `json:"open_id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
}
type OauthUserInfoChuanDII struct {
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Code     int    `json:"code"`
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Msg      string `json:"msg"`
}

// 添加oauth2登录回调处理接口.,调用第三方，根据返回id，username，password[创建私有加密]直接注册(静默)明面上直接登录
// 拓展接入oauth2

func oauth2CallBackHandler(tokenExpireTime time.Duration) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if r.Body == nil {
			return http.StatusForbidden, nil
		}
		info := &oauth2CallBackBody{}
		err := json.NewDecoder(r.Body).Decode(info)
		if err != nil {
			return http.StatusBadRequest, err
		}

		if info.Code == "" {
			return http.StatusBadRequest, nil
		}
		type Auth struct {
			Code         string `json:"code"`
			ClientId     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			GrantType    string `json:"grant_type"`
		}
		//auths := Auth{Code: , ClientId: "3a1cab58d4b743d681a042ba20955178", ClientSecret: "e21a3efafc7641a1a93542d14ffe8359", GrantType: "authorization_code"}
		//此处已经拿到了oauth2.0 code
		client := resty.New()
		authSuccessInfo := &AuthSuccess{}
		log.Println(r.Host)
		host, port, err := net.SplitHostPort(r.Host)
		log.Println(host + port + "-----------------")
		if err != nil {
			host = r.Host
		}
		_tokenUrl := fmt.Sprintf("%s?grant_type=authorization_code&client_id=%s&client_secret=%s&redirect_uri=%s&code=%s", d.Oauth2.Tokenurl, d.Oauth2.Clientid, d.Oauth2.Clientsecret, d.Oauth2.Redirecturi, info.Code)

		resp2, err := client.R().SetHeader("Content-Type", "application/json").
			SetResult(authSuccessInfo).
			Get(_tokenUrl)
		log.Println(resp2)
		if err != nil {
			log.Fatal(err)
		}
		if authSuccessInfo.Code != 0 {
			return http.StatusBadRequest, nil
		}
		type OauthUserInfo struct {
			OpenId         string `json:"openid"`
			Username       string `json:"username"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			Sex            string `json:"sex"`
			Status         int    `json:"status"`
			Avatar         string `json:"avatar"`
			Permission     string `json:"permission"`
			PermissionName string `json:"permission_name"`
			TGC            string `json:"tgc"`
			Email          string `json:"email"`
			Ticket         string `json:"ticket"`
			Code           string `json:"code"`
			StudentId      string `json:"student_id"`
			CreateTime     string `json:"create_time"`
			UpdateTime     string `json:"update_time"`
			RedirectUri    string `json:"redirect_uri"`
		}
		oauthUserInfoChuanDI := &OauthUserInfoChuanDI{}
		_meUrl := fmt.Sprintf("%s?access_token=%s", d.Oauth2.Meurl, authSuccessInfo.Token)
		resp3, err := client.R().SetResult(oauthUserInfoChuanDI).Get(_meUrl)
		log.Println("resp3")
		log.Println(resp3)
		if err != nil {
			log.Fatal(err)
			return http.StatusBadRequest, nil
		}
		oauthUserInfoChuanDII := &OauthUserInfoChuanDII{}
		_getUserInfoUrl := fmt.Sprintf("%s?oauth_consumer_key=%s&access_token=%s&openid=%s", d.Oauth2.Userinfourl, d.Oauth2.Clientid, authSuccessInfo.Token, oauthUserInfoChuanDI.OpenId)
		resp4, err := client.R().SetResult(oauthUserInfoChuanDII).Get(_getUserInfoUrl)
		log.Println(resp4)
		if err != nil {
			log.Fatal(err)
			return http.StatusBadRequest, nil
		}
		if oauthUserInfoChuanDII.Code != 0 {
			log.Fatal(err)
			return http.StatusBadRequest, nil
		}
		log.Println("oauthUserInfoChuanDII")
		log.Println(oauthUserInfoChuanDII)

		u, err := d.store.Users.Get(d.server.Root, oauthUserInfoChuanDII.Username)
		if err != nil {
			userddd := &users.User{
				Username: oauthUserInfoChuanDII.Username,
			}

			d.settings.Defaults.Apply(userddd)

			pwd, err := users.HashPwd("IUJHSW984e89H(HSHeh4H(hyu4gjn)(JH(r5rhgsgre5yerg")
			if err != nil {
				return http.StatusInternalServerError, err
			}

			userddd.Password = pwd

			userHome, err := d.settings.MakeUserDir(userddd.Username, userddd.Scope, d.server.Root)
			if err != nil {
				log.Printf("create user: failed to mkdir user home dir: [%s]", userHome)
				return http.StatusInternalServerError, err
			}
			userddd.Scope = userHome
			log.Printf("new user: %s, home dir: [%s].", userddd.Username, userHome)
			d.store.Users.Save(userddd)

			return printToken(w, r, d, userddd, tokenExpireTime)
		}
		return printToken(w, r, d, u, tokenExpireTime)
	}
}

var oauth2StateHandel = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if d.Oauth2.Disable {
		result := map[string]interface{}{
			"use": false,
		}
		return renderJSON(w, r, result)
	}
	result := map[string]interface{}{
		"use":  true,
		"url":  settings.GetAuthorizeUrl(&d.Oauth2),
		"name": d.Oauth2.Name,
	}
	return renderJSON(w, r, result)
}
