## local install
.PHONY: local-install
local-install:
	go install golang.org/x/tools/cmd/goimports
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

## fmt
.PHONY: fmt
fmt:
	goimports -w -local "game-api" cmd/ pkg/
	gofmt -s -w cmd/ pkg/

## lint
.PHONY: lint
lint:
	golangci-lint run -v cmd/... pkg/...

## token
token := "77a792a8-0523-4960-8f65-46e7fcc7fa20"

## exec create user
.PHONY: create_user
create_user:
	curl -X POST "http://localhost:8080/user/create" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{\"name\":\"tanaka\"}"

## exec get user
.PHONY: get_user
get_user:
	curl -X GET "http://localhost:8080/user/get" -H  "accept: application/json" -H  "x-token: $(token)"

# exec update user by name
.PHONY: update_user
update_user:
	curl -X PATCH "http://localhost:8080/user/update" -H  "accept: application/json" -H  "Content-Type: application/json" -H "x-token: $(token)" -d "{\"name\":\"updated_sample\"}"

## exec get collection list(user)
.PHONY: get_collection_list
get_collection_list:
	curl -X GET "http://localhost:8080/collection/list" -H  "accept: application/json" -H  "x-token: $(token)"

## exec get ranking list(user)
.PHONY: get_ranking_list
get_ranking_list:
	curl -X GET "http://localhost:8080/ranking/list?start=1" -H  "accept: application/json" -H  "x-token: $(token)"

## exec finish game
.PHONY: finish_game
finish_game:
	curl -X POST "http://localhost:8080/game/finish" -H  "accept: application/json" -H  "Content-Type: application/json" -H  "x-token: $(token)" -d "{\"score\":41110}"

## exec draw gacha
.PHONY: draw_gacha
draw_gacha:
	curl -X POST "http://localhost:8080/gacha/draw" -H  "accept: application/json" -H  "Content-Type: application/json" -H  "x-token: $(token)" -d "{\"times\":3}"
