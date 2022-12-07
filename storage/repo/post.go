package repo

import (
	pb "github.com/Asliddin3/post-servise/genproto/post"
)

type PostStorageI interface {
	CreatePost(*pb.PostRequest) (*pb.PostResponse, error)
	DeletePost(*pb.PostId) (*pb.Empty, error)
	UpdatePost(*pb.PostUpdate) (*pb.PostResponse, error)
	GetPost(*pb.PostId) (*pb.PostResponseCustomer, error)
	GetPostCustomerId(*pb.CustomerId) (*pb.ListPostCustomer, error)
	GetListPosts(*pb.Empty) (*pb.ListAllPostResponse, error)
	DeletePostByCustomerId(*pb.CustomerId) (*pb.DeletedReview, error)
	ListPost(limit int64, page int64) (*pb.ListPostResp, error)
	SearchOrderedPagePost(*pb.SearchRequest) (*pb.SearchResponse, error)
	CreateCustomer(*pb.CustomerResponse)(*pb.CustomerResponse,error)
}
