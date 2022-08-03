package azuretls

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (s *Session) caller(path string, information []byte) (string, error) {
	var realPath = s.server.endpoint + path + "?sid=" + strconv.FormatUint(s.id, 10)

	request, err := http.NewRequest("POST", realPath, bytes.NewBuffer(information))

	if err != nil {
		return "", err
	}

	request.Header = s.server.header

	response, err := s.server.client.Do(request)

	if err != nil {
		return "", err
	}

	return getBodyString(response), nil
}

func (serv *Server) Ping() bool {
	request, err := http.NewRequest("POST", serv.endpoint, bytes.NewBuffer([]byte{}))
	request.Header = serv.header

	if err != nil {
		return false
	}

	response, err := serv.client.Do(request)

	if err != nil {
		return false
	}

	body := getBodyString(response)

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	err = json.Unmarshal([]byte(body), &result)

	if err != nil {
		return false
	}

	return result["status"] == "ok"
}

func getBodyString(response *http.Response) string {
	defer response.Body.Close()

	encoding := response.Header.Get("content-encoding")

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return ""
	}

	return DecompressBody(bodyBytes, encoding)
}
