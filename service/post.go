package service

import (
	"context"
	"fmt"

	"github.com/Asliddin3/post-servise/genproto/customer"
	pb "github.com/Asliddin3/post-servise/genproto/post"
	"github.com/Asliddin3/post-servise/genproto/review"
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
func (r *PostService) GetListPosts(ctx context.Context, req *pb.Empty) (*pb.ListAllPostResponse, error) {
	postsInfo, err := r.storage.Post().GetListPosts(req)
	if err != nil {
		r.logger.Error("error getting all posts", l.Any("error getting  all posts", err))
		return &pb.ListAllPostResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return postsInfo, nil
}

func (r *PostService) GetPostCustomerId(ctx context.Context, req *pb.CustomerId) (*pb.ListPostCustomer, error) {
	customerPosts, err := r.storage.Post().GetPostCustomerId(req)
	if err != nil {
		r.logger.Error("error getting customer post", l.Any("error getting customer posts", err))
		return &pb.ListPostCustomer{}, status.Error(codes.Internal, "something went wrong")
	}
	for i, post := range customerPosts.Posts {
		// postReview, err := r.storage.Review().GetPostReview(&review.PostId{Id: post.Id})
		postReview, err := r.client.ReviewService().GetPostReview(context.Background(), &review.PostId{Id: post.Id})

		if err != nil {
			r.logger.Error("error getting post by customer", l.Any("error getting post reivew", err))
			return &pb.ListPostCustomer{}, status.Error(codes.Internal, "something went wrong")
		}
		postWithReview := pb.PostReviewResponse{
			Id:          post.Id,
			CustomerId:  post.CustomerId,
			Name:        post.Name,
			Media:       post.Media,
			Description: post.Description,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		}
		postWithReview.Count = postReview.Count
		postWithReview.Review = postReview.Review
		customerPosts.Posts[i] = &postWithReview
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

func (s *PostService) DeletePostByCustomerId(ctx context.Context, req *pb.CustomerId) (*pb.DeletedReview, error) {
	arrReview, err := s.storage.Post().DeletePostByCustomerId(req)
	if err != nil {
		s.logger.Error("error while deleting customer post", l.Any("error deleting customer post", err))
		return &pb.DeletedReview{}, status.Error(codes.Internal, "something went wrong ")
	}
	fmt.Println(err)
	for _, id := range arrReview.ReviewsIds {
		fmt.Println(id)
		_, err = s.client.ReviewService().DeleteReview(context.Background(), &review.PostId{Id: id.Id})
		if err != nil {
			return &pb.DeletedReview{}, err
		}
	}
	return arrReview, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.PostId) (*pb.Empty, error) {
	post, err := s.storage.Post().DeletePost(req)
	if err != nil {
		s.logger.Error("error while deleting Post", l.Any("error deleting post", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrong ")
	}
	fmt.Println("error in delet post", err)
	_, err = s.client.ReviewService().DeleteReview(context.Background(), &review.PostId{Id: req.Id})
	if err != nil {
		return &pb.Empty{}, err
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

func (s *PostService) GetPost(ctc context.Context, req *pb.PostId) (*pb.PostResponseCustomer, error) {
	post, err := s.storage.Post().GetPost(req)
	if err != nil {
		s.logger.Error("error while getting post", l.Any("error getting post", err))
		return &pb.PostResponseCustomer{}, status.Error(codes.Internal, "something went wrong")
	}
	customerInfo, err := s.client.CustomerService().GetCustomerInfo(ctc, &customer.CustomerId{Id: post.CustomerId})
	if err != nil {
		s.logger.Error("error while getting customer", l.Any("error getting customer", err))
		return &pb.PostResponseCustomer{}, status.Error(codes.Internal, "error getting customer")
	}
	post.Firstname = customerInfo.FirstName
	post.Lastname = customerInfo.LastName
	post.Email = customerInfo.Email
	post.Phonenumber = customerInfo.PhoneNumber
	for _, address := range customerInfo.Adderesses {
		post.Adderesses = append(post.Adderesses, &pb.AddressResponse{
			Id:       address.Id,
			District: address.District,
			Street:   address.Street,
		})
	}
	return post, nil
}
