all: ## Remove all & Generate all
	make clean && make gen
gen: ## Generate all file proto in folder gen
	[ -d "gen" ] || mkdir gen && protoc --go_out=./gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative  --go-grpc_out=require_unimplemented_servers=false:./gen protobuf/*.proto

clean: ## Remove all file proto in folder gen
	rm -rf gen/

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'