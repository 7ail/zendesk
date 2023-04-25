package zendesk

type gateway struct {
	client    doer
	subdomain string
	host      string
}

func NewGateway(client doer, subdomain, host string) (*gateway, error) {
	return &gateway{
		client:    client,
		subdomain: subdomain,
		host:      host,
	}, nil
}
