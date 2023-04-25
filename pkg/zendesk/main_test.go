package zendesk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGateway(t *testing.T) {
	subdomain := "subdomain"
	host := "host"

	act := func(doer doer, subdomain, host string) (*gateway, error) {
		return NewGateway(doer, subdomain, host)
	}

	testCases := []struct{
		name string
		arrange func(*testing.T) (doer, string, string)
		assert func(*testing.T, *gateway, error)
	}{
		{
			name: "happy path",
			arrange: func(t *testing.T) (doer, string, string) {
				return newMockDoer(t), subdomain, host
			},
			assert: func(t *testing.T, actualGateway *gateway, actualErr error) {
				assert.NoError(t, actualErr)
				assert.Equal(t, subdomain, actualGateway.subdomain)
				assert.Equal(t, host, actualGateway.host)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doer, subdomain, host := tc.arrange(t)
			actualGateway, actualErr := act(doer, subdomain, host)
			tc.assert(t, actualGateway, actualErr)
		})
	}
}