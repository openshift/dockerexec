GO_PACKAGES = $(shell find . -not \( -wholename ./vendor -prune -o -wholename ./.git -prune \) -name '*.go' -print0 | xargs -0n1 dirname | sort -u)

direct-build:
	go build -v $(GO_PACKAGES)

direct-install:
	go install -v $(GO_PACKAGES)
