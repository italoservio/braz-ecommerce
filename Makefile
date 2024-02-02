setup:
	git config --local core.hooksPath .commitlint/hooks
	chmod +x `pwd`/.commitlint/commitlint

start:
	docker compose up --detach --renew-anon-volumes --build --wait

stop:
	docker compose down

test:
	go test ./...

coverage:
	@mkdir -p .coverage
	@go test ./... -coverprofile=coverage.out > /dev/null || true
	@sed -i '/packages\/database\/setup\.go/d' coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@mv coverage.out .coverage
	@mv coverage.html .coverage
	@go tool cover -func .coverage/coverage.out
	@echo "Coverage HTML report can be found at: $(PWD)/.coverage/coverage.html"

local_docker_cmd:
	go install github.com/cosmtrek/air@latest \
    && air \
		--build.cmd "go build -o tmp/$(SERVICE_FOLDER) cmd/$(SERVICE_FOLDER)/main.go" \
		--build.bin "./tmp/$(SERVICE_FOLDER)"
