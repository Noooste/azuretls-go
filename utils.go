package azuretls

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

var client *http.Client
var localEndpoint = ""
var localHeaders = http.Header{}

func (s *Session) caller(path string, information []byte) (string, error) {
	var realPath = localEndpoint + path + "?sid=" + strconv.FormatUint(s.id, 10)

	request, err := http.NewRequest("POST", realPath, bytes.NewBuffer(information))

	if err != nil {
		return "", err
	}

	request.Header = localHeaders

	response, err := client.Do(request)

	if err != nil {
		return "", err
	}

	return getBodyString(response), nil
}

func ping() bool {
	if localEndpoint == "" || client == nil || localHeaders == nil {
		return false
	}

	request, err := http.NewRequest("POST", localEndpoint, bytes.NewBuffer([]byte{}))
	request.Header = localHeaders

	if err != nil {
		return false
	}

	response, err := client.Do(request)

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
