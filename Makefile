GOMARKDOC_VERSION = v0.4.1
LOCALGOBIN = $$PWD/.bin
ENV_VARS = GOBIN="$(LOCALGOBIN)" PATH="$(LOCALGOBIN):$$PATH"

GOMARKDOC:
	$(ENV_VARS) go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@$(GOMARKDOC_VERSION)

generate:
	$(ENV_VARS) go generate ./...

test:
	go test -count=1 -coverprofile=cover.out ./...
