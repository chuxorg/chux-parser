build: 
	go build -o chux-parser-api  -ldflags  "-X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD` -X main.Version=`git tag --sort=-version:refname | head -n 1`" app/main.go
run:
	go run -ldflags  "-X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD` -X main.Version=`git tag --sort=-version:refname | head -n 1`" app/main.go

kill:
	npx kill-port 8097

.PHONY: test
test:
	go test ./...

.PHONY: release-version
release-version:
	./scripts/release_version.sh

.PHONY: changelog
changelog:
	./scripts/changelog.sh