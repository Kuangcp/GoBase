package request

type MythHttpRequest struct {
	headers    map[string]string
	parameters map[string][]string
}

func (m MythHttpRequest) GetAuthType() string {
	panic("implement me")
}

func (m MythHttpRequest) GetCookie() []Cookie {
	panic("implement me")
}

func (m MythHttpRequest) GetHeader(s string) string {
	panic("implement me")
}

func (m MythHttpRequest) GetMethod() string {
	panic("implement me")
}

func (m MythHttpRequest) GetPathInfo() string {
	panic("implement me")
}

func (m MythHttpRequest) GetParameter(s string) string {
	panic("implement me")
}

func (m MythHttpRequest) GetQueryString() string {
	panic("implement me")
}

func (m MythHttpRequest) GetSessionId() string {
	panic("implement me")
}

func (m MythHttpRequest) GetSession() HttpSession {
	panic("implement me")
}

func (m MythHttpRequest) GetOrCreateSession(b bool) HttpSession {
	panic("implement me")
}

func (m MythHttpRequest) Login(s string, s2 string) {
	panic("implement me")
}

func (m MythHttpRequest) Logout() {
	panic("implement me")
}
