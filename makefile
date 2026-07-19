 SSO_BINARY            = ssoApp


## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:  build_sso
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"


## build_auth: builds the auth binary as a linux executable
build_sso:
	@echo "Building auth binary..."
	cd ../sso-service && bash -c source /etc/profile && env GOOS=linux CGO_ENABLED=0 go build -o  ${SSO_BINARY} ./
	@echo "Done!"

