package auth

import "net/http"

type UserInfo struct {
	Uid      string
	Username string
	Role     string
}

type Auth interface {
	Authentication(req *http.Request) (UserInfo, error)
}
