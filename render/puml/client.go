package puml

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Host string
}

func (c *Client) Render(uml io.Reader) (io.Reader, error) {
	noRedirect := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	b, err := ioutil.ReadAll(uml)
	if err != nil {
		return nil, err
	}

	resp, err := noRedirect.PostForm(fmt.Sprintf("%s/form", c.Host), url.Values{
		"text": []string{string(b)},
	})

	if err != nil {
		return nil, fmt.Errorf("create uml: %w", err)
	}

	id, err := extractID(resp.Header.Get("Location"))
	if err != nil {
		return nil, fmt.Errorf("extracting id: %w", err)
	}

	log.Println("uml id:", id)
	resp, err = http.Get(fmt.Sprintf("%s/png/%s", c.Host, id))
	if err != nil {
		return nil, fmt.Errorf("get png image: %w", err)
	}

	return resp.Body, nil
}

func extractID(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Path, "/")

	return parts[2], nil
}
