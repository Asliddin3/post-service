package repo

import (
	pb "github.com/Asliddin3/post-servise/genproto/post"
)

type PostStorageI interface {
	// CheckField(*pb.CheckFieldRequest) (*pb.CheckFieldResponse,error)
	CreatePost(*pb.PostRequest) (*pb.PostResponse, error)
	DeletePost(*pb.PostId) (*pb.Empty, error)
	UpdatePost(*pb.PostResponse) (*pb.PostResponse, error)
	GetPost(*pb.PostId) (*pb.PostResponse, error)
}
