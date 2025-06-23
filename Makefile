# 初始化依赖（ent相关）
init:
	go get entgo.io/ent/entc/gen@v0.14.4
	go get entgo.io/ent/cmd/internal/printer@v0.14.4
	go get entgo.io/ent/cmd/ent@v0.14.4


# generate ent code
ent:
	@go run entgo.io/ent/cmd/ent generate \
		--feature privacy \
		--feature sql/modifier \
		--feature intercept,schema/snapshot \
		--feature entql \
		--feature sql/upsert \
		./internal/data/ent/schema

.PHONY: proto

proto:
	protoc \
		--proto_path=proto \
		--go_out=internal/pb \
		--go_opt=paths=source_relative \
		proto/inspection.proto proto/config.proto proto/nacos.proto


