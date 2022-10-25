package service

import (
	"context"

	"github.com/Asliddin3/post-servise/genproto/customer"
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



func (r *ReviewService) DeleteReview(ctx context.Context, req *pb.ReviewId) (*pb.Empty, error) {
	_, err := r.client.ReviewService().DeleteReview(context.Background(), req)
	if err != nil {
		r.logger.Error("error deleting review", logger.Any("error getting post review", err))
		return &pb.Empty{}, status.Error(codes.Internal, "error deleting post review")
	}
	return &pb.Empty{}, nil
}

func (r *ReviewService) GetPostReviews(ctx context.Context, req *pb.PostId) (*pb.ReviewsList, error) {
	reviews, err := r.client.ReviewService().GetPostReviews(context.Background(), req)
	for _, review := range reviews.Reviews {
		customerInfo, err := r.client.CustomerService().GetCustomerInfo(context.Background(), &customer.CustomerId{Id: review.CustomerId})
		if err != nil {
			r.logger.Error("error getting customer info", l.Any("error getting customer", err))
			return &pb.ReviewsList{}, status.Error(codes.Internal, "something went wrong")
		}
		review.FirstName = customerInfo.FirstName
		review.LastName = customerInfo.LastName
	}
	if err != nil {
		r.logger.Error("error getting post reviews", l.Any("error getting post reviews", err))
		return &pb.ReviewsList{}, status.Error(codes.Internal, "something went wrong")
	}
	return reviews, nil
}

func (r *ReviewService) GetPostOverall(ctx context.Context, req *pb.PostId) (*pb.PostReview, error) {
	postReview, err := r.client.ReviewService().GetPostOverall(context.Background(), req)
	if err != nil {
		r.logger.Error("error getting post review", logger.Any("error getting post review", err))
		return &pb.PostReview{}, status.Error(codes.Internal, "something went wrong")
	}
	return postReview, nil
}
