BUILDS=./builds

build:
	tinygo build -o ${BUILDS}/auth.wasm -scheduler=none -target=wasi ./auth/main.go
	tinygo build -o ${BUILDS}/headers.wasm -scheduler=none -target=wasi ./headers/main.go
	curl http://localhost:8081/tyk/reload -H 'x-tyk-authorization: foo'

reload:
	curl http://localhost:8081/tyk/reload -H 'x-tyk-authorization: foo'
