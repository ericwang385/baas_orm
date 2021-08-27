package auth

import (
	"net/http"
)

type MockAuth struct{}

func (m *MockAuth) Authentication(req *http.Request) (UserInfo, error) {
	c, err := req.Cookie("uid")
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{Uid: c.Value}, err
}
