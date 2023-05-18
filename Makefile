#SWAGGER_PATH := "$(shell go env GOPATH)/bin/swagger"
SWAGGER_PATH := "/home/mehedi/go/bin/swagger"
print_var:
	echo $(SWAGGER_PATH)
check_install: print_var
	which $(SWAGGER_PATH) || go get github.com/go-swagger/go-swagger/cmd/swagger@latest
swagger: check_install
	$(SWAGGER_PATH) generate spec -o ./swagger.yaml --scan-models

