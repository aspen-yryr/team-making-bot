// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: proto/match/match.proto

package match

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MatchSvcClient is the client API for MatchSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchSvcClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	Create(ctx context.Context, in *CreateMatchRequest, opts ...grpc.CallOption) (*CreateMatchResponse, error)
	Find(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error)
	AppendMembers(ctx context.Context, in *AppendMemberRequest, opts ...grpc.CallOption) (*Match, error)
	Shuffle(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleResponse, error)
}

type matchSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchSvcClient(cc grpc.ClientConnInterface) MatchSvcClient {
	return &matchSvcClient{cc}
}

func (c *matchSvcClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/match.MatchSvc/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchSvcClient) Create(ctx context.Context, in *CreateMatchRequest, opts ...grpc.CallOption) (*CreateMatchResponse, error) {
	out := new(CreateMatchResponse)
	err := c.cc.Invoke(ctx, "/match.MatchSvc/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchSvcClient) Find(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error) {
	out := new(FindResponse)
	err := c.cc.Invoke(ctx, "/match.MatchSvc/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchSvcClient) AppendMembers(ctx context.Context, in *AppendMemberRequest, opts ...grpc.CallOption) (*Match, error) {
	out := new(Match)
	err := c.cc.Invoke(ctx, "/match.MatchSvc/AppendMembers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchSvcClient) Shuffle(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleResponse, error) {
	out := new(ShuffleResponse)
	err := c.cc.Invoke(ctx, "/match.MatchSvc/Shuffle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MatchSvcServer is the server API for MatchSvc service.
// All implementations must embed UnimplementedMatchSvcServer
// for forward compatibility
type MatchSvcServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	Create(context.Context, *CreateMatchRequest) (*CreateMatchResponse, error)
	Find(context.Context, *FindRequest) (*FindResponse, error)
	AppendMembers(context.Context, *AppendMemberRequest) (*Match, error)
	Shuffle(context.Context, *ShuffleRequest) (*ShuffleResponse, error)
	mustEmbedUnimplementedMatchSvcServer()
}

// UnimplementedMatchSvcServer must be embedded to have forward compatible implementations.
type UnimplementedMatchSvcServer struct {
}

func (UnimplementedMatchSvcServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedMatchSvcServer) Create(context.Context, *CreateMatchRequest) (*CreateMatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedMatchSvcServer) Find(context.Context, *FindRequest) (*FindResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedMatchSvcServer) AppendMembers(context.Context, *AppendMemberRequest) (*Match, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendMembers not implemented")
}
func (UnimplementedMatchSvcServer) Shuffle(context.Context, *ShuffleRequest) (*ShuffleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shuffle not implemented")
}
func (UnimplementedMatchSvcServer) mustEmbedUnimplementedMatchSvcServer() {}

// UnsafeMatchSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchSvcServer will
// result in compilation errors.
type UnsafeMatchSvcServer interface {
	mustEmbedUnimplementedMatchSvcServer()
}

func RegisterMatchSvcServer(s grpc.ServiceRegistrar, srv MatchSvcServer) {
	s.RegisterService(&MatchSvc_ServiceDesc, srv)
}

func _MatchSvc_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchSvcServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/match.MatchSvc/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchSvcServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchSvc_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchSvcServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/match.MatchSvc/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchSvcServer).Create(ctx, req.(*CreateMatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchSvc_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchSvcServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/match.MatchSvc/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchSvcServer).Find(ctx, req.(*FindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchSvc_AppendMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchSvcServer).AppendMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/match.MatchSvc/AppendMembers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchSvcServer).AppendMembers(ctx, req.(*AppendMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchSvc_Shuffle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShuffleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchSvcServer).Shuffle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/match.MatchSvc/Shuffle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchSvcServer).Shuffle(ctx, req.(*ShuffleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MatchSvc_ServiceDesc is the grpc.ServiceDesc for MatchSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MatchSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "match.MatchSvc",
	HandlerType: (*MatchSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _MatchSvc_CreateUser_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _MatchSvc_Create_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _MatchSvc_Find_Handler,
		},
		{
			MethodName: "AppendMembers",
			Handler:    _MatchSvc_AppendMembers_Handler,
		},
		{
			MethodName: "Shuffle",
			Handler:    _MatchSvc_Shuffle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/match/match.proto",
}