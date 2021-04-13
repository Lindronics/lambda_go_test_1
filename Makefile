.PHONY: build

test:
	cd hello-world; go test -v -cover .

generate_model:
	cd hello-world/data; swagger generate model --spec=spec.json

build: generate_model
	sam build
