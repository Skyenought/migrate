fmt:
	gofumpt -l -w -extra .
ast:
	go run main.go -mode ast
