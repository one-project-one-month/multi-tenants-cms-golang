package cookies

import (
	"fmt"
	"net/url"
	"time"
)

type CookieBuilder struct {
	name     string
	value    string
	domain   string
	path     string
	maxAge   int
	expires  time.Time
	secure   bool
	httpOnly bool
	sameSite string
}

func NewCookieBuilder(name, value string) *CookieBuilder {
	return &CookieBuilder{
		name:     name,
		value:    value,
		path:     "/",
		httpOnly: true,
		secure:   true,
		sameSite: "Strict",
	}
}

func (cb *CookieBuilder) Domain(domain string) *CookieBuilder {
	cb.domain = domain
	return cb
}

func (cb *CookieBuilder) Path(path string) *CookieBuilder {
	cb.path = path
	return cb
}

func (cb *CookieBuilder) MaxAge(seconds int) *CookieBuilder {
	cb.maxAge = seconds
	return cb
}

func (cb *CookieBuilder) Expires(expires time.Time) *CookieBuilder {
	cb.expires = expires
	return cb
}

func (cb *CookieBuilder) Secure(secure bool) *CookieBuilder {
	cb.secure = secure
	return cb
}

func (cb *CookieBuilder) HttpOnly(httpOnly bool) *CookieBuilder {
	cb.httpOnly = httpOnly
	return cb
}

func (cb *CookieBuilder) SameSite(sameSite string) *CookieBuilder {
	cb.sameSite = sameSite
	return cb
}

func (cb *CookieBuilder) Build() string {
	cookie := fmt.Sprintf("%s=%s", cb.name, url.QueryEscape(cb.value))

	if cb.domain != "" {
		cookie += fmt.Sprintf("; Domain=%s", cb.domain)
	}

	if cb.path != "" {
		cookie += fmt.Sprintf("; Path=%s", cb.path)
	}

	if cb.maxAge > 0 {
		cookie += fmt.Sprintf("; Max-Age=%d", cb.maxAge)
	}

	if !cb.expires.IsZero() {
		cookie += fmt.Sprintf("; Expires=%s", cb.expires.UTC().Format(time.RFC1123))
	}

	if cb.secure {
		cookie += "; Secure"
	}

	if cb.httpOnly {
		cookie += "; HttpOnly"
	}

	if cb.sameSite != "" {
		cookie += fmt.Sprintf("; SameSite=%s", cb.sameSite)
	}

	return cookie
}
