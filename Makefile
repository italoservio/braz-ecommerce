setup:
	git config --local core.hooksPath .commitlint/hooks
	chmod +x `pwd`/.commitlint/commitlint

start:
	docker compose up --detach --renew-anon-volumes --build --wait
	
stop:
	docker compose down
