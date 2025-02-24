//go:build tools
// +build tools

package tools

import (
	// lock tool gomodifytags version
	// lock tool protoc-gen-go version
	// lock tool protoc-gen-go-grpc version

	_ "github.com/fatih/gomodifytags"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
