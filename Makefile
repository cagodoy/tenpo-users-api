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
	@go run cmd/main.go

build b:
	@echo "[build] Building service..."
	@cd cmd && go build -o $(BIN)

linux l:
	@echo "[build-linux] Building service..."
	@cd cmd && GOOS=linux GOARCH=amd64 go build -o $(BIN)

docker d:
	@echo "[copy] Copy parent bin..."
	@cp ../../bin/goose ../../bin/wait-db bin
	@echo "[docker] Building image..."
	@docker build -t $(SVC):$(VERSION) .
	@echo "[remove] Removing parent bin..."
	@rm bin/goose bin/wait-db

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

.PHONY: clean c run r build b linux l docker d add-migration am migrations m