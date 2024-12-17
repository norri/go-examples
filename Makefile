test:
	find . -name 'go.mod' -execdir go test ./... \;

lint:
	# install golangci-lint https://golangci-lint.run/welcome/install/
	find . -name 'go.mod' -execdir golangci-lint run \;
