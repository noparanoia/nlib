package nlib

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func writeHeaders(request *http.Request, headers map[string]string) {
	if headers != nil && len(headers) != 0 {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}
}

// Get запрос, принимает в качестве аргумента map[string]string в котором содержатся заголовки,uri адрес
func Get(headers map[string]string, uri string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return []byte{}, err
	}

	writeHeaders(request, headers)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	res, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		if Body != nil {
			_ = Body.Close()
		}
	}(res.Body)

	return body, nil
}

// Post запрос, принимает в качестве аргумента map[string]string в котором содержатся заголовки,uri адрес, data []byte
func Post(headers map[string]string, uri string, data []byte) ([]byte, error) {
	var buf *bytes.Buffer
	if data != nil {
		buf = bytes.NewBuffer(data)
	}
	request, err := http.NewRequest(http.MethodPost, uri, buf)
	if err != nil {
		return []byte{}, err
	}

	writeHeaders(request, headers)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	res, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		if Body != nil {
			_ = Body.Close()
		}
	}(res.Body)

	return body, nil
}

// TryTCPNetAddress Проверяем доступность сервера через TCP
func TryTCPNetAddress(addr string) error {
	l, err := url.Parse(addr)
	if err != nil {
		return err
	}
	c, err := net.DialTimeout("tcp", l.Host, 15*time.Second)
	defer func(c net.Conn) {
		if c == nil {
			return
		}
		cErr := c.Close()
		if cErr != nil {
			return
		}
	}(c)
	if err != nil {
		return err
	}
	return nil
}
