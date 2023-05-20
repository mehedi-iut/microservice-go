#SWAGGER_PATH := "$(shell go env GOPATH)/bin/swagger"
SWAGGER_PATH := "/home/mehedi/go/bin/swagger"

check_install:
	which $(SWAGGER_PATH) || go get github.com/go-swagger/go-swagger/cmd/swagger@latest
swagger: check_install
	$(SWAGGER_PATH) generate spec -o ./swagger.yaml --scan-models

