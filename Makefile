generate-config:
	$(CURDIR)/scripts/generate-config
build: generate-config
	test -d $(CURDIR)/bin || mkdir $(CURDIR)/bin
	go build -C ./git-wrapper -o $(CURDIR)/bin/git-wrapper
build-watch:
	while inotifywait -r -e modify,move,create,delete ./git-wrapper; do make build; done
runtests:
	go test -C ./git-wrapper ./...

