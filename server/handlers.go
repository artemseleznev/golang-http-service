package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const maxUrlsCount = 20

type ReqBody struct {
	Urls []string `json:"urls"`
}

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

func (s *Server) Contents(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		errorResponse(w, http.StatusMethodNotAllowed, "")
		return
	}
	decoder := json.NewDecoder(req.Body)
	var reqBody ReqBody
	if err := decoder.Decode(&reqBody); err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not parse body: %s", err.Error()))
		return
	}
	if err := validateUrls(reqBody.Urls); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var urlContents sync.Map
	g, ctx := errgroup.WithContext(req.Context())
	for _, u := range reqBody.Urls {
		u := u
		g.Go(func() error {
			log.Printf("url: %s, start...", u)
			urlContent, err := s.getContent(ctx, u)
			if err != nil {
				log.Printf("url: %s, got error: %s", u, err.Error())
				return err
			}
			urlContents.Store(u, urlContent)
			log.Printf("url: %s, got content!", u)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		errorResponse(w, http.StatusBadGateway, err.Error())
		return
	}

	result := make(map[string]string, len(reqBody.Urls))
	urlContents.Range(func(key, val interface{}) bool {
		result[fmt.Sprint(key)] = fmt.Sprintf("%v", val)
		return true
	})
	okResponse(w, result)

}

func validateUrls(urls []string) error {
	if len(urls) > maxUrlsCount {
		return errors.New("too many urls in request")
	}
	for _, u := range urls {
		if _, err := url.Parse(u); err != nil {
			return fmt.Errorf("url `%s` is invalid: %w", u, err)
		}
	}
	return nil
}

func (s *Server) getContent(ctx context.Context, url string) (string, error) {
	// todo punycode conversion
	client := http.Client{
		Timeout: s.reqTimeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
