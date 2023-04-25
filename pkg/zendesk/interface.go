package zendesk

import "net/http"

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}
