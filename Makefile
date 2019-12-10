#
# SO variables
#
# GITHUB_USER
# GITHUB_TOKEN
#

#
# Internal variables
#
VERSION=0.0.1
SVC=tenpo-users-api
BIN_PATH=$(PWD)/bin
BIN=$(BIN_PATH)/$(SVC)
GITHUB_REGISTRY_URL=docker.pkg.github.com/$(GITHUB_USER)/$(SVC)

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

run r:
	@echo "[running] Running service..."
	@go run cmd/server/main.go

build b:
	@echo "[build] Building service..."
	@cd cmd/server && go build -o $(BIN)

linux l:
	@echo "[build-linux] Building service..."
	@cd cmd/server && GOOS=linux GOARCH=amd64 go build -o $(BIN)

wait-db wd:
	@cd cmd/wait-db && GOOS=linux GOARCH=amd64 go build -o $(BIN_PATH)/wait-db 

docker d:
	@echo "[docker] Building image..."
	@docker build -t $(SVC):$(VERSION) .

add-migration am: 
	@echo "[add-migration] Adding migration"
	@goose -dir "./database/migrations" create $(name) sql

migrations m:
	@echo "[migrations] Runing migrations..."
	@cd database/migrations && goose postgres $(DSN) up

docker-login dl:
	@echo "[docker] Login to docker..."
	@docker login docker.pkg.github.com -u $(GITHUB_USER) -p $(GITHUB_TOKEN)

push p: linux docker docker-login
	@echo "[docker] pushing $(GITHUB_REGISTRY_URL)/$(SVC):$(VERSION)"
	@docker tag $(SVC):$(VERSION) $(GITHUB_REGISTRY_URL)/$(SVC):$(VERSION)
	@docker push $(GITHUB_REGISTRY_URL)/$(SVC):$(VERSION)

.PHONY: clean c run r build b linux l wait-db wd docker d add-migration am migrations m