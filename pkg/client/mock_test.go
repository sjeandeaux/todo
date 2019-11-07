package client_test

import (
	"context"

	. "github.com/onsi/gomega"
	pb "github.com/sjeandeaux/todo/todo-grpc/v1"
	grpc "google.golang.org/grpc"
)

type mockToDoServiceClient struct {
	expectedRequest interface{}

	response interface{}
	err      error
}

var _ pb.ToDoServiceClient = &mockToDoServiceClient{}

func (s *mockToDoServiceClient) Create(ctx context.Context, r *pb.CreateRequest, opts ...grpc.CallOption) (*pb.CreateResponse, error) {
	Ω(r).Should(Equal(s.expectedRequest))
	if s.response == nil {
		return nil, s.err
	}
	return s.response.(*pb.CreateResponse), s.err
}

func (s *mockToDoServiceClient) Read(ctx context.Context, r *pb.ReadRequest, opts ...grpc.CallOption) (*pb.ReadResponse, error) {
	Ω(r).Should(Equal(s.expectedRequest))
	if s.response == nil {
		return nil, s.err
	}
	return s.response.(*pb.ReadResponse), s.err
}

func (s *mockToDoServiceClient) Update(ctx context.Context, r *pb.UpdateRequest, opts ...grpc.CallOption) (*pb.UpdateResponse, error) {
	Ω(r).Should(Equal(s.expectedRequest))
	if s.response == nil {
		return nil, s.err
	}
	return s.response.(*pb.UpdateResponse), s.err
}

func (s *mockToDoServiceClient) Delete(ctx context.Context, r *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteResponse, error) {
	Ω(r).Should(Equal(s.expectedRequest))
	if s.response == nil {
		return nil, s.err
	}
	return s.response.(*pb.DeleteResponse), s.err
}

func (s *mockToDoServiceClient) Search(ctx context.Context, r *pb.SearchRequest, opts ...grpc.CallOption) (*pb.SearchResponse, error) {
	Ω(r).Should(Equal(s.expectedRequest))
	if s.response == nil {
		return nil, s.err
	}
	return s.response.(*pb.SearchResponse), s.err
}
