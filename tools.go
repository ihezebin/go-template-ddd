//go:build tools
// +build tools

package tools

import (
	// lock tool gomodifytags version
	_ "github.com/fatih/gomodifytags"
	// lock tool protoc-gen-go version
	_ "github.com/golang/protobuf/protoc-gen-go"
	// lock tool swag version
	_ "github.com/swaggo/swag/cmd/swag"
	// lock tool protoc-gen-go-grpc version
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
