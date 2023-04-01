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
	./release_version.sh

.PHONY: changelog
changelog:
	echo "# Changelog" > CHANGELOG.md
	git tag --sort=-version:refname | while read -r TAG; do \
	  echo -e "\n## $$TAG\n" >> CHANGELOG.md; \
	  if [ "$$PREVIOUS_TAG" != "" ]; then \
	    git log --no-merges --format="* %s (%h)" $$TAG..$$PREVIOUS_TAG >> CHANGELOG.md; \
	  else \
	    git log --no-merges --format="* %s (%h)" $$TAG >> CHANGELOG.md; \
	  fi; \
	  PREVIOUS_TAG=$$TAG; \
	done