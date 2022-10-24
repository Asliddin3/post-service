package service

import (
	"context"

	pb "github.com/Asliddin3/post-servise/genproto/review"
	"github.com/Asliddin3/post-servise/pkg/logger"
	l "github.com/Asliddin3/post-servise/pkg/logger"
	grpcclient "github.com/Asliddin3/post-servise/service/grpc_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	client *grpcclient.ServiceManager
	logger l.Logger
}

// func (r *ReviewService) CreateReview(ctx context.Context, req *pb.Review) (*pb.Review, error) {
// 	reviewResp, err := r.client.ReviewService().CreateReview(context.Background(), req)
// 	if err != nil {
// 		r.logger.Error("error while creating review", logger.Any("creating review argument error", err))
// 		return &pb.Review{}, status.Error(codes.Internal, "Please check your argument")
// 	}
// 	return reviewResp, nil
// }

func (r *ReviewService) DeleteReview(ctx context.Context, req *pb.PostId) (*pb.Empty, error) {
	_, err := r.client.ReviewService().DeleteReview(context.Background(), req)
	if err != nil {
		r.logger.Error("error deleting review", logger.Any("error getting post review", err))
		return &pb.Empty{}, status.Error(codes.Internal, "error deleting post review")
	}
	return &pb.Empty{}, nil
}

func (r *ReviewService) GetPostReviews(ctx context.Context, req *pb.PostId) (*pb.ReviewsList, error) {
	reviews, err := r.client.ReviewService().GetPostReviews(context.Background(), req)
	if err != nil {
		r.logger.Error("error getting post reviews", l.Any("error getting post reviews", err))
		return &pb.ReviewsList{}, status.Error(codes.Internal, "something went wrong")
	}
	return reviews, nil
}

func (r *ReviewService) GetPostReview(ctx context.Context, req *pb.PostId) (*pb.PostReview, error) {
	postReview, err := r.client.ReviewService().GetPostReview(context.Background(), req)
	if err != nil {
		r.logger.Error("error getting post review", logger.Any("error getting post review", err))
		return &pb.PostReview{}, status.Error(codes.Internal, "something went wrong")
	}
	return postReview, nil
}
