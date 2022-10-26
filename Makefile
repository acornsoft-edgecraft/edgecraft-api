.PHONY: clean build run

APP_NAME = edgecraft-apiserver
CMD_DIR = $(PWD)/cmd
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/scripts/migrations
DATABASE_URL = postgres://edgecraft:edgecraft@192.168.77.42:31000/edgecraft?sslmode=disable
# DATABASE_URL = postgres://postgres:password@cgapp-postgres/postgres?sslmode=disable

clean:
	rm -rf ./build

build: 
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME) --kubeconfig=strategy=file,path=$${HOME}/.kube/config

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

docker.run: docker.network docker.postgres swag docker.edgecraft docker.redis migrate.up
# docker.run: docker.network docker.postgres docker.edgecraft docker.redis migrate.up

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.edgecraft.build:
	docker build -t edgecraft-api .

docker.edgecraft: docker.edgecraft.build
	docker run --rm -d \
		--name edgecraft-api \
		--network dev-network \
		-p 8100:8100 \
		edgecraft-api

docker.postgres:
	docker run --rm -d \
		--name edgecraft-postgres \
		--network dev-network \
		-e POSTGRES_USER=edgecraft \
		-e POSTGRES_PASSWORD=edgecraft \
		-e POSTGRES_DB=edgecraft \
		-v ${HOME}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.stop: docker.stop.edgecraft docker.stop.postgres

docker.stop.edgecraft:
	docker stop cgapp-edgecraft

docker.stop.postgres:
	docker stop cgapp-postgres

swag:
	swag init

# swag:
# 	swag init --parseDependency --parseInternal -d $(CMD_DIR)/$(APP_NAME) -o $(PWD)/api