// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: plugin.proto

package proto

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

// NegotiatorClient is the client API for Negotiator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NegotiatorClient interface {
	Negotiate(ctx context.Context, opts ...grpc.CallOption) (Negotiator_NegotiateClient, error)
}

type negotiatorClient struct {
	cc grpc.ClientConnInterface
}

func NewNegotiatorClient(cc grpc.ClientConnInterface) NegotiatorClient {
	return &negotiatorClient{cc}
}

func (c *negotiatorClient) Negotiate(ctx context.Context, opts ...grpc.CallOption) (Negotiator_NegotiateClient, error) {
	stream, err := c.cc.NewStream(ctx, &Negotiator_ServiceDesc.Streams[0], "/plugin.Negotiator/Negotiate", opts...)
	if err != nil {
		return nil, err
	}
	x := &negotiatorNegotiateClient{stream}
	return x, nil
}

type Negotiator_NegotiateClient interface {
	Send(*NegotiateRequest) error
	Recv() (*NegotiateResponse, error)
	grpc.ClientStream
}

type negotiatorNegotiateClient struct {
	grpc.ClientStream
}

func (x *negotiatorNegotiateClient) Send(m *NegotiateRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *negotiatorNegotiateClient) Recv() (*NegotiateResponse, error) {
	m := new(NegotiateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NegotiatorServer is the server API for Negotiator service.
// All implementations must embed UnimplementedNegotiatorServer
// for forward compatibility
type NegotiatorServer interface {
	Negotiate(Negotiator_NegotiateServer) error
	mustEmbedUnimplementedNegotiatorServer()
}

// UnimplementedNegotiatorServer must be embedded to have forward compatible implementations.
type UnimplementedNegotiatorServer struct {
}

func (UnimplementedNegotiatorServer) Negotiate(Negotiator_NegotiateServer) error {
	return status.Errorf(codes.Unimplemented, "method Negotiate not implemented")
}
func (UnimplementedNegotiatorServer) mustEmbedUnimplementedNegotiatorServer() {}

// UnsafeNegotiatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NegotiatorServer will
// result in compilation errors.
type UnsafeNegotiatorServer interface {
	mustEmbedUnimplementedNegotiatorServer()
}

func RegisterNegotiatorServer(s grpc.ServiceRegistrar, srv NegotiatorServer) {
	s.RegisterService(&Negotiator_ServiceDesc, srv)
}

func _Negotiator_Negotiate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NegotiatorServer).Negotiate(&negotiatorNegotiateServer{stream})
}

type Negotiator_NegotiateServer interface {
	Send(*NegotiateResponse) error
	Recv() (*NegotiateRequest, error)
	grpc.ServerStream
}

type negotiatorNegotiateServer struct {
	grpc.ServerStream
}

func (x *negotiatorNegotiateServer) Send(m *NegotiateResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *negotiatorNegotiateServer) Recv() (*NegotiateRequest, error) {
	m := new(NegotiateRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Negotiator_ServiceDesc is the grpc.ServiceDesc for Negotiator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Negotiator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.Negotiator",
	HandlerType: (*NegotiatorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Negotiate",
			Handler:       _Negotiator_Negotiate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "plugin.proto",
}

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginClient interface {
	Load(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*LoadResponse, error)
	Terminate(ctx context.Context, in *TerminateRequest, opts ...grpc.CallOption) (*TerminateResponse, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) Load(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*LoadResponse, error) {
	out := new(LoadResponse)
	err := c.cc.Invoke(ctx, "/plugin.Plugin/Load", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Terminate(ctx context.Context, in *TerminateRequest, opts ...grpc.CallOption) (*TerminateResponse, error) {
	out := new(TerminateResponse)
	err := c.cc.Invoke(ctx, "/plugin.Plugin/Terminate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
// All implementations must embed UnimplementedPluginServer
// for forward compatibility
type PluginServer interface {
	Load(context.Context, *LoadRequest) (*LoadResponse, error)
	Terminate(context.Context, *TerminateRequest) (*TerminateResponse, error)
	mustEmbedUnimplementedPluginServer()
}

// UnimplementedPluginServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (UnimplementedPluginServer) Load(context.Context, *LoadRequest) (*LoadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Load not implemented")
}
func (UnimplementedPluginServer) Terminate(context.Context, *TerminateRequest) (*TerminateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Terminate not implemented")
}
func (UnimplementedPluginServer) mustEmbedUnimplementedPluginServer() {}

// UnsafePluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServer will
// result in compilation errors.
type UnsafePluginServer interface {
	mustEmbedUnimplementedPluginServer()
}

func RegisterPluginServer(s grpc.ServiceRegistrar, srv PluginServer) {
	s.RegisterService(&Plugin_ServiceDesc, srv)
}

func _Plugin_Load_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Load(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Plugin/Load",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Load(ctx, req.(*LoadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Terminate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TerminateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Terminate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Plugin/Terminate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Terminate(ctx, req.(*TerminateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Plugin_ServiceDesc is the grpc.ServiceDesc for Plugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Plugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Load",
			Handler:    _Plugin_Load_Handler,
		},
		{
			MethodName: "Terminate",
			Handler:    _Plugin_Terminate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}

// ConnectorClient is the client API for Connector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectorClient interface {
	Ls(ctx context.Context, opts ...grpc.CallOption) (Connector_LsClient, error)
	Cat(ctx context.Context, in *ReadFileRequest, opts ...grpc.CallOption) (*File, error)
}

type connectorClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectorClient(cc grpc.ClientConnInterface) ConnectorClient {
	return &connectorClient{cc}
}

func (c *connectorClient) Ls(ctx context.Context, opts ...grpc.CallOption) (Connector_LsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Connector_ServiceDesc.Streams[0], "/plugin.Connector/Ls", opts...)
	if err != nil {
		return nil, err
	}
	x := &connectorLsClient{stream}
	return x, nil
}

type Connector_LsClient interface {
	Send(*ListChildrenRequest) error
	Recv() (*DirectoryEntry, error)
	grpc.ClientStream
}

type connectorLsClient struct {
	grpc.ClientStream
}

func (x *connectorLsClient) Send(m *ListChildrenRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *connectorLsClient) Recv() (*DirectoryEntry, error) {
	m := new(DirectoryEntry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *connectorClient) Cat(ctx context.Context, in *ReadFileRequest, opts ...grpc.CallOption) (*File, error) {
	out := new(File)
	err := c.cc.Invoke(ctx, "/plugin.Connector/Cat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectorServer is the server API for Connector service.
// All implementations must embed UnimplementedConnectorServer
// for forward compatibility
type ConnectorServer interface {
	Ls(Connector_LsServer) error
	Cat(context.Context, *ReadFileRequest) (*File, error)
	mustEmbedUnimplementedConnectorServer()
}

// UnimplementedConnectorServer must be embedded to have forward compatible implementations.
type UnimplementedConnectorServer struct {
}

func (UnimplementedConnectorServer) Ls(Connector_LsServer) error {
	return status.Errorf(codes.Unimplemented, "method Ls not implemented")
}
func (UnimplementedConnectorServer) Cat(context.Context, *ReadFileRequest) (*File, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cat not implemented")
}
func (UnimplementedConnectorServer) mustEmbedUnimplementedConnectorServer() {}

// UnsafeConnectorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectorServer will
// result in compilation errors.
type UnsafeConnectorServer interface {
	mustEmbedUnimplementedConnectorServer()
}

func RegisterConnectorServer(s grpc.ServiceRegistrar, srv ConnectorServer) {
	s.RegisterService(&Connector_ServiceDesc, srv)
}

func _Connector_Ls_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ConnectorServer).Ls(&connectorLsServer{stream})
}

type Connector_LsServer interface {
	Send(*DirectoryEntry) error
	Recv() (*ListChildrenRequest, error)
	grpc.ServerStream
}

type connectorLsServer struct {
	grpc.ServerStream
}

func (x *connectorLsServer) Send(m *DirectoryEntry) error {
	return x.ServerStream.SendMsg(m)
}

func (x *connectorLsServer) Recv() (*ListChildrenRequest, error) {
	m := new(ListChildrenRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Connector_Cat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServer).Cat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Connector/Cat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServer).Cat(ctx, req.(*ReadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Connector_ServiceDesc is the grpc.ServiceDesc for Connector service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Connector_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.Connector",
	HandlerType: (*ConnectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Cat",
			Handler:    _Connector_Cat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Ls",
			Handler:       _Connector_Ls_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "plugin.proto",
}

// AnalyzerClient is the client API for Analyzer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalyzerClient interface {
	Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error)
}

type analyzerClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyzerClient(cc grpc.ClientConnInterface) AnalyzerClient {
	return &analyzerClient{cc}
}

func (c *analyzerClient) Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error) {
	out := new(AnalyzeResponse)
	err := c.cc.Invoke(ctx, "/plugin.Analyzer/Analyze", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalyzerServer is the server API for Analyzer service.
// All implementations must embed UnimplementedAnalyzerServer
// for forward compatibility
type AnalyzerServer interface {
	Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error)
	mustEmbedUnimplementedAnalyzerServer()
}

// UnimplementedAnalyzerServer must be embedded to have forward compatible implementations.
type UnimplementedAnalyzerServer struct {
}

func (UnimplementedAnalyzerServer) Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Analyze not implemented")
}
func (UnimplementedAnalyzerServer) mustEmbedUnimplementedAnalyzerServer() {}

// UnsafeAnalyzerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalyzerServer will
// result in compilation errors.
type UnsafeAnalyzerServer interface {
	mustEmbedUnimplementedAnalyzerServer()
}

func RegisterAnalyzerServer(s grpc.ServiceRegistrar, srv AnalyzerServer) {
	s.RegisterService(&Analyzer_ServiceDesc, srv)
}

func _Analyzer_Analyze_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyzerServer).Analyze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Analyzer/Analyze",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyzerServer).Analyze(ctx, req.(*AnalyzeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Analyzer_ServiceDesc is the grpc.ServiceDesc for Analyzer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Analyzer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.Analyzer",
	HandlerType: (*AnalyzerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Analyze",
			Handler:    _Analyzer_Analyze_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}

// TransformerClient is the client API for Transformer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransformerClient interface {
	Transform(ctx context.Context, in *TransformRequest, opts ...grpc.CallOption) (*TransformResponse, error)
}

type transformerClient struct {
	cc grpc.ClientConnInterface
}

func NewTransformerClient(cc grpc.ClientConnInterface) TransformerClient {
	return &transformerClient{cc}
}

func (c *transformerClient) Transform(ctx context.Context, in *TransformRequest, opts ...grpc.CallOption) (*TransformResponse, error) {
	out := new(TransformResponse)
	err := c.cc.Invoke(ctx, "/plugin.Transformer/Transform", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransformerServer is the server API for Transformer service.
// All implementations must embed UnimplementedTransformerServer
// for forward compatibility
type TransformerServer interface {
	Transform(context.Context, *TransformRequest) (*TransformResponse, error)
	mustEmbedUnimplementedTransformerServer()
}

// UnimplementedTransformerServer must be embedded to have forward compatible implementations.
type UnimplementedTransformerServer struct {
}

func (UnimplementedTransformerServer) Transform(context.Context, *TransformRequest) (*TransformResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transform not implemented")
}
func (UnimplementedTransformerServer) mustEmbedUnimplementedTransformerServer() {}

// UnsafeTransformerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransformerServer will
// result in compilation errors.
type UnsafeTransformerServer interface {
	mustEmbedUnimplementedTransformerServer()
}

func RegisterTransformerServer(s grpc.ServiceRegistrar, srv TransformerServer) {
	s.RegisterService(&Transformer_ServiceDesc, srv)
}

func _Transformer_Transform_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransformerServer).Transform(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Transformer/Transform",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransformerServer).Transform(ctx, req.(*TransformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Transformer_ServiceDesc is the grpc.ServiceDesc for Transformer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Transformer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.Transformer",
	HandlerType: (*TransformerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transform",
			Handler:    _Transformer_Transform_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}