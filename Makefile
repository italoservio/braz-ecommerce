setup:
	git config --local core.hooksPath .commitlint/hooks
	chmod +x `pwd`/.commitlint/commitlint
	
local_docker_cmd:
	go install github.com/cosmtrek/air@latest \
    && air \
		--build.cmd "go build -o tmp/$(SERVICE_FOLDER) cmd/$(SERVICE_FOLDER)/main.go" \
		--build.bin "./tmp/$(SERVICE_FOLDER)"

start:
	docker compose up --detach --renew-anon-volumes --build --wait
	
stop:
	docker compose down
