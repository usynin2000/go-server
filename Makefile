.PHONY: run clear-port


run:
	go run .

clear-port:
	kill -9 $(lsof -t -i :3000)
