package request

type (
	Cookie struct {
	}

	HttpRequest interface {
		GetAuthType() string
		GetCookie() []Cookie
		GetHeader(string) string
		GetMethod() string
		GetPathInfo() string
		GetQueryString() string
		GetParameter(string) string
		GetSessionId() string
		GetSession() HttpSession
		GetOrCreateSession(bool) HttpSession
		Login(string, string)
		Logout()
	}

	HttpSession interface {
		GetID() string
		GetCreationTime() int64
		GetLastAccessTime() int64
		GetMaxInactiveInterval() int
		SetMaxInactiveInterval(int)
		GetAttribute(string) interface{}
		RemoveAttribute(string)
		SetAttribute(string, interface{})
		Invalidate()
	}
)
