package traefik_hydrate_headers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"
)

type Hydrate struct {
	next     http.Handler
	name     string
	client   http.Client
	config   *Config
	template *template.Template
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Remote.Url) == 0 {
		return nil, fmt.Errorf("remote.url cannot be empty")
	}

	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	h := &Hydrate{
		config: config,
		next:   next,
		name:   name,
		client: http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 30 * time.Second,
		},
		template: template.New("").Funcs(template.FuncMap{
			"unmarshalJson": func(s string) interface{} {
				var data interface{}

				json.Unmarshal([]byte(s), &data)

				return data
			},
		}),
	}

	return h, nil
}

func (h *Hydrate) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if len(h.config.FetchOn.Cookies) > 0 {
		var cookieFound bool = false

		for _, cookie := range h.config.FetchOn.Cookies {
			_, err := req.Cookie(cookie)
			if err == nil {
				cookieFound = true
				break
			}
		}

		if !cookieFound {
			h.NextIfRequired(rw, req, nil)
			return
		}
	}

	remoteReq, err := http.NewRequest(h.config.Remote.Method, h.config.Remote.Url, nil)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, header := range h.config.ForwardHeaders {
		remoteReq.Header.Add(header, req.Header.Get(header))
	}

	remoteRes, err := h.client.Do(remoteReq)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(h.config.AppendOn.StatusCodes) > 0 && !contains(h.config.AppendOn.StatusCodes, remoteRes.StatusCode) {
		h.NextIfRequired(rw, req, remoteRes)
		return
	}

	body, err := io.ReadAll(remoteRes.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer remoteRes.Body.Close()

	var bodyString string

	if strings.Contains(remoteRes.Header.Get("Content-Type"), "application/json") {
		compactBody, err := compactJson(body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		bodyString = string(compactBody)
	} else {
		bodyString = strings.ReplaceAll(string(body), "\n", "\\n")
	}

	data := make(map[string]interface{})
	data["Request"] = req
	data["RemoteBody"] = bodyString
	data["RemoteResponse"] = remoteRes

	for key, value := range h.config.Headers {
		tpl, err := h.template.Parse(strings.TrimSpace(value))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer := &bytes.Buffer{}

		err = tpl.Execute(writer, data)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		req.Header.Add(key, strings.TrimSpace(writer.String()))
	}

	h.NextIfRequired(rw, req, remoteRes)
}

func (h *Hydrate) NextIfRequired(rw http.ResponseWriter, req *http.Request, remoteRes *http.Response) {
	if remoteRes != nil && len(h.config.NextOn.StatusCodes) > 0 && !contains(h.config.NextOn.StatusCodes, remoteRes.StatusCode) {
		return
	}

	h.next.ServeHTTP(rw, req)
}
