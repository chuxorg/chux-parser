build: 
	go build ./...
.PHONY: test
test:
	go test ./...

.PHONY: release-version
release-version:
	./scripts/release_version.sh

.PHONY: changelog
changelog:
	./scripts/changelog.sh