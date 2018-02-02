package trace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ZipKin struct {
	endpoint string // format  ip:port
	api      string // default /api/v2/spans
}

var gzipkin = ZipKin{endpoint: "http://127.0.0.1:9411", api: "/api/v2/spans"}

func ZipKinEndpointSet(endpoint string) {
	gzipkin.endpoint = "http://" + endpoint
}

func readfully(conn io.ReadCloser) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}

	return result.Bytes(), nil
}

func httpRequest(method string, url string, req []byte) (rsp []byte, status int, err error) {

	body := bytes.NewBuffer(req)

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	request.Header.Add("Content-Type", "application/json")

	rspon, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, 0, err
	}

	defer rspon.Body.Close()

	rsp, err = readfully(rspon.Body)
	if err != nil {
		return nil, rspon.StatusCode, err
	}

	return rsp, rspon.StatusCode, nil
}

func PostSpan(spans interface{}) error {

	body, err := json.Marshal(spans)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	_, status, err := httpRequest("POST", gzipkin.endpoint+gzipkin.api, body)
	if err != nil {
		return err
	}

	if status != 202 {
		return errors.New(fmt.Sprintf("http rsponse status %d", status))
	}

	return nil
}
