package zendesk

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrNotFound = errors.New("not found error")
var ErrUnexpected = errors.New("unexpected error")

func (g *gateway) IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func (g *gateway) isErrUnexpected(err error) bool {
	return errors.Is(err, ErrUnexpected)
}

func convertToError(statusCode int) error {
	if statusCode == http.StatusNotFound {
		return ErrNotFound
	}

	return errors.Join(ErrUnexpected, fmt.Errorf("unexpected status code %d", statusCode))
}
