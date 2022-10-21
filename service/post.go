package service

import (
	"context"

	pb "github.com/Asliddin3/post-servise/genproto/post"
	l "github.com/Asliddin3/post-servise/pkg/logger"
	grpcclient "github.com/Asliddin3/post-servise/service/grpc_client"
	"github.com/Asliddin3/post-servise/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	storage storage.IStorage
	client  *grpcclient.ServiceManager
	logger  l.Logger
}

func NewPostService(client *grpcclient.ServiceManager, db *sqlx.DB, log l.Logger) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		client:  client,
		logger:  log,
	}
}
func (r *PostService) GetPostCustomerId(ctx context.Context, req *pb.CustomerId) (*pb.ListPostCustomer, error) {
	customerPosts, err := r.storage.Post().GetPostCustomerId(req)
	if err != nil {
		r.logger.Error("error getting customer post", l.Any("error getting customer posts", err))
		return &pb.ListPostCustomer{}, status.Error(codes.Internal, "something went wrong")
	}
	return customerPosts, nil
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	Post, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error("error while creating Post", l.Any("error creating Post", err))
		return &pb.PostResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return Post, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.PostId) (*pb.Empty, error) {
	post, err := s.storage.Post().DeletePost(req)
	if err != nil {
		s.logger.Error("error while deleting Post", l.Any("error deleting post", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrong ")
	}
	return post, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.PostResponse) (*pb.PostResponse, error) {
	post, err := s.storage.Post().UpdatePost(req)
	if err != nil {
		s.logger.Error("error while updating post", l.Any("error updating", err))
		return &pb.PostResponse{}, status.Error(codes.Internal, "somthing went wrong please check argument")
	}
	return post, nil
}

func (s *PostService) GetPost(ctc context.Context, req *pb.PostId) (*pb.PostResponse, error) {
	post, err := s.storage.Post().GetPost(req)
	if err != nil {
		s.logger.Error("error while getting post", l.Any("error getting post", err))
		return &pb.PostResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return post, nil
}
