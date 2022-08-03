package azuretls

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

//InitApiRemote init api with ip and authorization key
func InitApiRemote(ip string, key string) *Server {
	if !strings.Contains(ip, "http://") && !strings.Contains(ip, "https://") {
		ip = "https://" + ip
	}

	if ip[len(ip)-1] == '/' {
		ip = ip[:len(ip)-1]
	}

	server := &Server{
		client: &http.Client{
			Timeout: 30 * time.Second, Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		header: http.Header{},
	}

	server.header.Set("authorization", key)
	server.header.Set("content-type", "application/json")

	server.endpoint = ip

	if server.Ping() {
		return server
	} else {
		return nil
	}
}

//InitSession create a new session on the API and return its struct.
func (serv *Server) InitSession() (*Session, error) {
	session := &Session{
		server:      serv,
		Header:      map[string]string{},
		PHeader:     []string{},
		HeaderOrder: []string{},
		Timeout:     30,
		Cookies:     map[string]string{},
		Navigator:   "chrome",
	}

	var realPath = serv.endpoint + "/session/new"

	request, err := http.NewRequest("POST", realPath, bytes.NewBuffer([]byte{}))

	if err != nil {
		return nil, err
	}

	request.Header = serv.header

	response, err := serv.client.Do(request)

	if err != nil {
		return nil, err
	}

	body := getBodyString(response)

	var sessionResponse sessionResponse
	log.Print(body)

	err = json.Unmarshal([]byte(body), &sessionResponse)

	if err != nil {
		return nil, err
	}
	log.Print(sessionResponse)

	if sessionResponse.Success == true {
		session.id = sessionResponse.Sid
	} else {
		return nil, errors.New("couldn't create any session")
	}

	return session, nil

}
