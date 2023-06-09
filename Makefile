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

.PHONY: reset-tags
reset-tags:
	git tag -l | xargs git tag -d
	git fetch --tags
	git ls-remote --tags origin | awk '/refs\/tags\// {sub("refs/tags/", "", $2); print ":" $2}' | xargs -I {} git push origin {}

