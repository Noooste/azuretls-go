package azuretls

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// GetCookies return cookies from a given domain
func (s *Session) GetCookies(domain string) ([]Cookie, error) {
	response, err := s.caller("/session/cookies", []byte(`{"domain":"`+domain+`"}`))

	if err != nil {
		return nil, err
	}
	var cookiesList []Cookie

	if err := json.Unmarshal([]byte(response), &cookiesList); err == nil {
		return cookiesList, nil
	} else {
		return nil, err
	}
}

// SetCookies set cookies into the current session
func (s *Session) SetCookies(cookies []Cookie) bool {
	cookiesString, err := json.Marshal(cookies)

	if err != nil {
		panic(err)
	}

	response, err := s.caller("/session/cookies/set", cookiesString)

	if err != nil {
		panic(err)
	}

	var success Status
	if err := json.Unmarshal([]byte(response), &success); err == nil {
		return success.Success
	} else {
		return false
	}
}

// Close the current session
func (s *Session) Close() bool {
	if result, err := s.caller("/session/close", []byte{}); err != nil {
		return false
	} else {
		var success Status
		if err := json.Unmarshal([]byte(result), &success); err == nil {
			return success.Success
		} else {
			return false
		}
	}
}

// KeepAlive asks the API to keep the session up
func (s *Session) KeepAlive() bool {
	if result, err := s.caller("/session/keep-alive", []byte{}); err != nil {
		return false
	} else {
		var success Status
		if err := json.Unmarshal([]byte(result), &success); err == nil {
			return success.Success
		} else {
			return false
		}
	}
}

// ApplyJA3 apply your ja3 and its specifications to the current session
func (s *Session) ApplyJA3(ja3 string, specifications map[string]interface{}) bool {
	if information, err := json.Marshal(ja3Information{
		Ja3:            ja3,
		Specifications: specifications,
		Navigator:      s.Navigator,
	}); err == nil {
		if result, err := s.caller("/session/tls/ja3", information); err != nil {
			return false
		} else {
			var success Status
			if err := json.Unmarshal([]byte(result), &success); err == nil {
				return success.Success
			} else {
				return false
			}
		}
	} else {
		panic(err)
	}
}

const (
	HeaderTableSize      = "HEADER_TABLE_SIZE"
	EnablePush           = "ENABLE_PUSH"
	MaxConcurrentStreams = "MAX_CONCURRENT_STREAMS"
	InitialWindowsSize   = "INITIAL_WINDOW_SIZE"
	MaxFrameSize         = "MAX_FRAME_SIZE"
	MaxHeaderListSize    = "MAX_HEADER_LIST_SIZE"
)

// ApplyHTTP2Settings apply given HTTP2Settings to the current session (only on http2 requests)
func (s *Session) ApplyHTTP2Settings(settings []HTTP2Settings) bool {
	if information, err := json.Marshal(settings); err == nil {
		if result, err := s.caller("/session/http2/settings", []byte(`{"settings":`+string(information)+`}`)); err != nil {
			return false
		} else {
			var success Status
			if err := json.Unmarshal([]byte(result), &success); err == nil {
				return success.Success
			} else {
				return false
			}
		}
	} else {
		panic(err)
	}
}

// ApplyWindowsUpdate apply given value to the session's Windows update (only on http2 requests)
func (s *Session) ApplyWindowsUpdate(value int) bool {
	if value > 0 && value < 1<<31 {
		if result, err := s.caller("/session/http2/windows-update", []byte(`{"value":`+strconv.Itoa(value)+`}`)); err != nil {
			return false
		} else {
			var success Status
			if err := json.Unmarshal([]byte(result), &success); err == nil {
				return success.Success
			} else {
				return false
			}
		}
	} else {
		panic("WINDOWS_UPDATE error : The legal range for the increment to the flow-control window is 1 to 2^31-1 (2,147,483,647) octets")
	}
}

// ApplyStreamPriorities apply given StreamInformation information to the current session (only on http2 requests)
func (s *Session) ApplyStreamPriorities(streams []StreamInformation) bool {
	if information, err := json.Marshal(streams); err == nil {
		if result, err := s.caller("/session/http2/stream-priorities", []byte(`{"streams":`+string(information)+`}`)); err != nil {
			return false
		} else {
			var success Status
			if err := json.Unmarshal([]byte(result), &success); err == nil {
				return success.Success
			} else {
				return false
			}
		}
	} else {
		panic(err)
	}
}

// Do the given Request.
func (s *Session) Do(request Request) (*Response, error) {
	if request.Header == nil {
		request.Header = s.Header
	}
	if request.PHeader == nil {
		request.PHeader = s.PHeader
	}
	if request.HeaderOrder == nil {
		request.HeaderOrder = s.HeaderOrder
	}
	if request.Proxy == "" {
		request.Proxy = s.Proxy
	}
	if request.Navigator == "" {
		request.Navigator = s.Navigator
	}
	if request.Timeout == 0 {
		request.Timeout = s.Timeout
	}

	if information, err := json.Marshal(request); err == nil {
		if result, err := s.caller("/session/request", information); err != nil {
			return nil, err
		} else {
			var response Response
			if err := json.Unmarshal([]byte(result), &response); err == nil {
				return &response, nil
			} else {
				var apiErr apiError
				if err := json.Unmarshal([]byte(result), &apiErr); err == nil {
					return nil, errors.New(apiErr.Error)
				} else {
					return nil, err
				}
			}
		}
	} else {
		panic(err)
	}
}

func (h Header) Get(key string) string {
	key = strings.ToLower(key)
	for hKey, hValue := range h {
		if key == strings.ToLower(hKey) {
			return hValue
		}
	}
	return ""
}
