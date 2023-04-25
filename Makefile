GOPATH := ${shell go env GOPATH}

mock:
	${GOPATH}/bin/mockery --name=doer --dir=./pkg/zendesk --inpackage --filename a_mock_doer.go
