package zendesk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotFoundError(t *testing.T) {
	act := func(t *testing.T, err error) bool {
		zendesk, _ := NewGateway(newMockDoer(t), "subdomain", "host")

		return zendesk.IsErrNotFound(err)
	}

	testCases := []struct{
		name string
		arrange func() error
		assert func(*testing.T, bool)
	}{
		{
			name: "is ErrNotFound",
			arrange: func() error {
				return ErrNotFound
			},
			assert: func(t *testing.T, actualBool bool) {
				assert.True(t, actualBool)
			},
		},
		{
			name: "is not ErrNotFound",
			arrange: func() error {
				return ErrUnexpected
			},
			assert: func(t *testing.T, actualBool bool) {
				assert.False(t, actualBool)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualErr := tc.arrange()
			actualBool := act(t, actualErr)
			tc.assert(t, actualBool)
		})
	}
}

func TestIsUnexpectedError(t *testing.T) {
	act := func(t *testing.T, err error) bool {
		zendesk, _ := NewGateway(newMockDoer(t), "subdomain", "host")

		return zendesk.isErrUnexpected(err)
	}

	testCases := []struct{
		name string
		arrange func() error
		assert func(*testing.T, bool)
	}{
		{
			name: "is ErrUnexpected",
			arrange: func() error {
				return ErrUnexpected
			},
			assert: func(t *testing.T, actualBool bool) {
				assert.True(t, actualBool)
			},
		},
		{
			name: "is not ErrUnexpected",
			arrange: func() error {
				return ErrNotFound
			},
			assert: func(t *testing.T, actualBool bool) {
				assert.False(t, actualBool)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualErr := tc.arrange()
			actualBool := act(t, actualErr)
			tc.assert(t, actualBool)
		})
	}
}