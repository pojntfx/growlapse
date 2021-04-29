all: build

agent:
	go build -o out/growlapse-agent/growlapse-agent cmd/growlapse-agent/main.go

frontend:
	rm -f web/app.wasm
	GOOS=js GOARCH=wasm go build -o web/app.wasm cmd/growlapse-frontend/main.go
	go build -o /tmp/growlapse-frontend-build cmd/growlapse-frontend/main.go
	rm -rf out/growlapse-frontend
	/tmp/growlapse-frontend-build -build
	cp -r web/* out/growlapse-frontend/web

build: agent frontend

release-agent:
	[ "$$(uname -m)" = 'armv6l' ] && go build -o out/release/growlapse-agent/growlapse-agent.linux-$$(uname -m) cmd/growlapse-agent/main.go || CGO_ENABLED=1 go build -ldflags="-extldflags=-static" -tags netgo -o out/release/growlapse-agent/growlapse-agent.linux-$$(uname -m) cmd/growlapse-agent/main.go\

release-frontend: frontend
	rm -rf out/release/growlapse-frontend
	mkdir -p out/release/growlapse-frontend
	cd out/growlapse-frontend && tar -czvf ../release/growlapse-frontend/growlapse-frontend.tar.gz .

release-frontend-github-pages: frontend
	rm -rf out/release/growlapse-frontend-github-pages
	mkdir -p out/release/growlapse-frontend-github-pages
	/tmp/growlapse-frontend-build -build -path growlapse -out out/release/growlapse-frontend-github-pages
	cp -r web/* out/release/growlapse-frontend-github-pages/web

release: release-agent release-frontend release-frontend-github-pages

install: release-agent
	sudo install out/release/growlapse-agent/growlapse-agent.linux-$$(uname -m) /usr/local/bin/growlapse-agent
	
dev:
	while [ -z "$$FRONTEND_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd web/*.css)" ]; do\
		$(MAKE);\
		kill -9 $$FRONTEND_PID 2>/dev/null 1>&2;\
		wait $$FRONTEND_PID;\
		/tmp/growlapse-frontend-build -serve & export FRONTEND_PID="$$!";\
	done

clean:
	rm -rf out

depend:
	# Generate bindings
	go generate ./...