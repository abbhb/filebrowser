package settings

// Oauth2  settings of the app.
// 符合qq的oauth2
type Oauth2 struct {
	Name         string `json:"name"`         // name 为oauth登录按钮的提示文字
	Disable      bool   `json:"disable"`      // disable true为禁用
	Tokenurl     string `json:"tokenurl"`     // code换token
	Meurl        string `json:"meurl"`        // token拿openid
	Userinfourl  string `json:"userinfourl"`  // 获取用户信息包含用户名信息的接口地址
	Clientid     string `json:"clientid"`     // 客户端id
	Clientsecret string `json:"clientsecret"` // 密钥
	Redirecturi  string `json:"redirecturi"`  // 重定向地址
	Scope        string `json:"scope"`        // 权限，逗号分隔
	State        string `json:"state"`        // 防攻击字段
	Authorizeurl string `json:"authorizeurl"` // 一键登录地址
}

func GetAuthorizeUrl(oauth2 *Oauth2) string {
	return oauth2.Authorizeurl + "?response_type=code&client_id=" + oauth2.Clientid + "&redirect_uri=" + oauth2.Redirecturi
}
