SWAGGER_PATH := "$(shell go env GOPATH)/bin/swagger"
check_install:
	which $(SWAGGER_PATH) || go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.4
swagger: check_install
	$(SWAGGER_PATH) generate spec -o ./swagger.yaml --scan-models