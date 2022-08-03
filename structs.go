package azuretls

import "net/http"

type Server struct {
	client   *http.Client
	header   http.Header
	endpoint string
}

type Session struct {
	id     uint64
	server *Server

	Header      map[string]string
	PHeader     []string
	HeaderOrder []string

	Navigator string
	Cookies   map[string]string

	Timeout int
	Proxy   string
}

type Request struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Data   string `json:"data"`

	PHeader     []string          `json:"pheader"`
	Header      map[string]string `json:"header"`
	HeaderOrder []string          `json:"header-order"`

	Navigator string `json:"navigator"`

	Proxy         string `json:"proxy"`
	AllowRedirect bool   `json:"allow-redirect"`
	Timeout       int    `json:"timeout"`

	ServerPush bool `json:"server-push"`
	Verify     bool `json:"verify"`
}

type Response struct {
	StatusCode int `json:"status_code"`

	Cookies map[string]interface{} `json:"cookies"`
	Url     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
	Text    string                 `json:"body"`

	ServerPush []Response `json:"server-push"`
}

type sessionResponse struct {
	Success bool   `json:"success"`
	Sid     uint64 `json:"session-id"`
}

type Status struct {
	Success bool `json:"success"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path    string `json:"path"`
	Domain  string `json:"domain"`
	Expires string `json:"expires"`

	MaxAge   int  `json:"max-age"`
	Secure   bool `json:"secure"`
	HttpOnly bool `json:"http-only"`
}

type apiError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type ja3Information struct {
	Ja3            string                 `json:"ja3"`
	Specifications map[string]interface{} `json:"specifications"`
	Navigator      string                 `json:"navigator"`
}

type HTTP2Settings struct {
	name  string
	value int
}

type StreamInformation struct {
	StreamId  uint32 `json:"stream-id"`
	StreamDep uint32 `json:"stream-dep"`
	Exclusive bool   `json:"exclusive"`
	Weight    uint8  `json:"weight"`
}
