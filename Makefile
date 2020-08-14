all: lint_proto build_proto

lint_proto:
						buf check lint

build_proto:
						prototool generate

