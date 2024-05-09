package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filebrowser/filebrowser/v2/users"
	"io"
	"net/http"
	"sort"
)

var syncUserGetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	//response := map[string]interface{}{}
	users, err := d.store.Users.Gets(d.server.Root)
	response := []map[string]interface{}{}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	for _, u := range users {
		u.Password = ""
		response = append(response, map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"scope":    u.Scope,
		})
	}
	w.Header().Set("Content-Disposition", "attachment; filename=users.json")
	w.Header().Set("Content-Type", "application/json")
	marsh, err := json.Marshal(response)
	w.Write([]byte(marsh))
	return 0, nil
})

// 用户结构体
type UserSync struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Scope    string `json:"scope"`
}

var syncUserPostHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// 解析表单数据
	err := r.ParseMultipartForm(10 << 20) // 10 MB 最大内存
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to parse form data")

	}

	// 获取上传的文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, errors.New("Failed to get file from form data")
	}
	defer file.Close()

	// 检查文件类型是否为 JSON
	if handler.Header.Get("Content-Type") != "application/json" {
		return http.StatusInternalServerError, errors.New("Expected content type application/json")

	}

	// 打开临时文件
	//file.Seek(0, 0)
	// 反序列化 JSON 数据
	fileBytes, err := io.ReadAll(file)

	var usersyncs []UserSync
	err = json.Unmarshal(fileBytes, &usersyncs)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to decode JSON data")

	}
	// 处理用户数据
	for _, user := range usersyncs {
		if user.Username == "admin" {
			continue
		}
		newuser := &users.User{
			Username: user.Username,
		}
		pwd, err := users.HashPwd("IUJHSW984e89H(HSHeh4H(hyu4gjn)(JH(r5rhgsgre5yerg")
		if err != nil {
		}
		newuser.Password = pwd
		d.settings.Defaults.Apply(newuser)
		newuser.Scope = user.Scope
		err = d.store.Users.Save(newuser)
		if err != nil {
			fmt.Println(err)
		}
	}
	return 0, nil
})
