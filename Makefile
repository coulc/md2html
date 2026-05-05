.PHONY: run build

run: 
	@go run ./cmd/md2html/main.go


build: 
	@go build ./cmd/md2html/main.go
