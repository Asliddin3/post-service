package repo

import (
	pb "github.com/Asliddin3/post-servise/genproto/review"
)

type ReviewStorageI interface {
	GetPostReview(*pb.PostId) (*pb.PostReview, error)
	DeleteReview(*pb.PostId) (*pb.Empty, error)
}
