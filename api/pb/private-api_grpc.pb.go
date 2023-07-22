// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: private-api.proto

package pigeomail_api_pb

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const ()

// PrivateAPIClient is the client API for PrivateAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrivateAPIClient interface {
}

type privateAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewPrivateAPIClient(cc grpc.ClientConnInterface) PrivateAPIClient {
	return &privateAPIClient{cc}
}

// PrivateAPIServer is the server API for PrivateAPI service.
// All implementations should embed UnimplementedPrivateAPIServer
// for forward compatibility
type PrivateAPIServer interface {
}

// UnimplementedPrivateAPIServer should be embedded to have forward compatible implementations.
type UnimplementedPrivateAPIServer struct {
}

// UnsafePrivateAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrivateAPIServer will
// result in compilation errors.
type UnsafePrivateAPIServer interface {
	mustEmbedUnimplementedPrivateAPIServer()
}

func RegisterPrivateAPIServer(s grpc.ServiceRegistrar, srv PrivateAPIServer) {
	s.RegisterService(&PrivateAPI_ServiceDesc, srv)
}

// PrivateAPI_ServiceDesc is the grpc.ServiceDesc for PrivateAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrivateAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pigeomail.PrivateAPI",
	HandlerType: (*PrivateAPIServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "private-api.proto",
}
