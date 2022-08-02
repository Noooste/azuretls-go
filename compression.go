package azuretls

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/andybalholm/brotli"
	"io/ioutil"
)

func GUnzipData(data []byte) (resData []byte, err error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}
	defer gz.Close()
	respBody, err := ioutil.ReadAll(gz)
	return respBody, err
}

func EnflateData(data []byte) (resData []byte, err error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}
	defer zr.Close()
	enflated, err := ioutil.ReadAll(zr)
	return enflated, err
}

func UnBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := ioutil.ReadAll(br)
	return respBody, err
}

// DecompressBody unzips compressed data
func DecompressBody(Body []byte, encoding string) (parsedBody string) {
	if len(encoding) > 0 {
		if encoding == "gzip" {
			unz, err := GUnzipData(Body)
			if err != nil {
				return string(Body)
			}
			parsedBody = string(unz)

		} else if encoding == "deflate" {
			unz, err := EnflateData(Body)
			if err != nil {
				return string(Body)
			}
			parsedBody = string(unz)

		} else if encoding == "br" {
			unz, err := UnBrotliData(Body)
			if err != nil {
				return string(Body)
			}
			parsedBody = string(unz)

		} else {
			fmt.Print("Unknown Encoding" + encoding)
			parsedBody = string(Body)
		}

	} else {
		parsedBody = string(Body)
	}

	return parsedBody
}
