package postgres

import (
	pb "github.com/Asliddin3/post-servise/genproto/post"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) GetPost(req *pb.PostId) (*pb.PostResponse, error) {
	postResp := pb.PostResponse{}
	err := r.db.QueryRow(
		`select id,customer_id,name,description,created_at,updated_at from post  where id=$1`, req.Id,
	).Scan(&postResp.Id, &postResp.CustomerId, &postResp.Name, &postResp.Description, &postResp.CreatedAt, &postResp.UpdatedAt)
	if err != nil {
		return &pb.PostResponse{}, err
	}
	rows, err := r.db.Query(
		`select post_id,name,link,type from media where post_id=$1`, req.Id,
	)
	if err != nil {
		return &pb.PostResponse{}, err
	}
	for rows.Next() {
		medResp := pb.MediasResponse{}
		err = rows.Scan(&medResp.PostId, &medResp.Name, &medResp.Link, &medResp.Type)
		if err != nil {
			return &pb.PostResponse{}, err
		}
		postResp.Media = append(postResp.Media, &medResp)
	}

	return &postResp, nil
}

func (r *postRepo) GetPostCustomerId(req *pb.CustomerId) (*pb.ListPostCustomer, error) {
	posts := []*pb.PostResponse{}
	rows, err := r.db.Query(`
	 select id,customer_id,name,description,created_at,updated_at from post  where customer_id=$1
	`, req.Id)
	if err != nil {
		return &pb.ListPostCustomer{}, err
	}
	for rows.Next() {
		postResp := pb.PostResponse{}
		err := rows.Scan(&postResp.Id, &postResp.CustomerId, &postResp.Name, &postResp.Description, &postResp.CreatedAt, &postResp.UpdatedAt)
		if err != nil {
			return &pb.ListPostCustomer{}, err
		}
		posts = append(posts, &postResp)
	}
	return &pb.ListPostCustomer{
		Posts: posts,
	}, nil
}

func (r *postRepo) CreatePost(req *pb.PostRequest) (*pb.PostResponse, error) {
	postResp := pb.PostResponse{}
	err := r.db.QueryRow(`
	insert into post(customer_id,name,description)
	values($1,$2,$3) returning id,customer_id,name,description,created_at,updated_at
	`, req.CustomerId, req.Name, req.Description).Scan(&postResp.Id, &postResp.CustomerId, &postResp.Name, &postResp.Description,
		&postResp.CreatedAt, &postResp.UpdatedAt)
	if err != nil {
		return &pb.PostResponse{}, err
	}
	for _, media := range req.Media {
		mediaResp := pb.MediasResponse{}
		err = r.db.QueryRow(`
			insert into media(post_id,name,link,type)
			values($1,$2,$3,$4)
			returning post_id,name,link,type
		`, postResp.Id, media.Name, media.Link, media.Type,
		).Scan(&mediaResp.PostId, &mediaResp.Name, &mediaResp.Link, &mediaResp.Type)
		if err != nil {
			return &pb.PostResponse{}, err
		}
		postResp.Media = append(postResp.Media, &mediaResp)
	}
	return &postResp, nil
}

func (r *postRepo) UpdatePost(req *pb.PostResponse) (*pb.PostResponse, error) {
	postResp := pb.PostResponse{}
	err := r.db.QueryRow(`
	update post set name=$1,description=$2,updated_at=current_timestamp where id=$3
	returning id,name,description,created_at,updated_at
	`, req.Name, req.Description, req.Id).Scan(
		&postResp.Id, &postResp.Name, &postResp.Description, &postResp.CreatedAt, &postResp.UpdatedAt,
	)
	if err != nil {
		return &pb.PostResponse{}, err
	}
	return &postResp, nil
}
func (r *postRepo) DeletePost(req *pb.PostId) (*pb.Empty, error) {
	_, err := r.db.Exec(`
	update post set deleted_at=current_timestamp where id=$1
	`, req.Id)
	if err != nil {
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}
