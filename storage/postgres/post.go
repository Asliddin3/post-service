package postgres

import (
	"fmt"

	pb "github.com/Asliddin3/post-servise/genproto/post"

	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) GetPost(req *pb.PostId) (*pb.PostResponseCustomer, error) {
	postResp := pb.PostResponseCustomer{}
	fmt.Println(req.Id)
	err := r.db.QueryRow(
		`select id,customer_id,name,description,created_at,updated_at from post where id=$1 and deleted_at is null;`, req.Id,
	).Scan(&postResp.Id, &postResp.CustomerId, &postResp.Name, &postResp.Description, &postResp.CreatedAt, &postResp.UpdatedAt)
	fmt.Println(postResp)
	if err != nil {
		return &pb.PostResponseCustomer{}, err
	}
	rows, err := r.db.Query(
		`select post_id,name,link,type from media where post_id=$1`, req.Id,
	)
	if err != nil {
		return &pb.PostResponseCustomer{}, err
	}
	for rows.Next() {
		medResp := pb.MediasResponse{}
		err = rows.Scan(&medResp.PostId, &medResp.Name, &medResp.Link, &medResp.Type)
		if err != nil {
			return &pb.PostResponseCustomer{}, err
		}
		postResp.Media = append(postResp.Media, &medResp)
	}

	return &postResp, nil
}

func (r *postRepo) GetListPosts(req *pb.Empty) (*pb.ListAllPostResponse, error) {
	posts := &pb.ListAllPostResponse{}
	CleanMap := func(mapOfFunc map[int]string) {
		for k := range mapOfFunc {
			delete(mapOfFunc, k)
		}
	}
	rows, err := r.db.Query(`
	select id,deleted_at from post where deleted_at is not null
	`)
	if err != nil {
		return &pb.ListAllPostResponse{}, err
	}
	deletedPost := make(map[int]string)
	for rows.Next() {
		var id int
		var deleted_at string
		err = rows.Scan(&id, &deleted_at)
		if err != nil {
			return &pb.ListAllPostResponse{}, err
		}
		deletedPost[id] = deleted_at
	}
	rows, err = r.db.Query(`
	 select id,customer_id,name,description,created_at,updated_at from post
	`)
	if err != nil {
		return &pb.ListAllPostResponse{}, err
	}
	defer CleanMap(deletedPost)

	for rows.Next() {
		postResp := pb.PostReviewResponse{}
		err = rows.Scan(&postResp.Id, &postResp.CustomerId, &postResp.Name, &postResp.Description,
			&postResp.CreatedAt, &postResp.UpdatedAt)
		if err != nil {
			return &pb.ListAllPostResponse{}, err
		}
		medias, err := r.db.Query(`
		select name,link,type from media where post_id=$1
		`, postResp.Id)
		if err != nil {
			return &pb.ListAllPostResponse{}, err
		}
		fmt.Println(postResp.Id)
		for medias.Next() {
			mediaResp := pb.MediasResponse{}
			err = medias.Scan(&mediaResp.Name, &mediaResp.Link, &mediaResp.Type)
			if err != nil {
				return &pb.ListAllPostResponse{}, err
			}
			mediaResp.PostId = postResp.Id
			postResp.Media = append(postResp.Media, &mediaResp)
			fmt.Println(mediaResp)
		}
		if val, ok := deletedPost[int(postResp.Id)]; ok {
			postResp.DeletedAt = val
			posts.DeletedPost = append(posts.DeletedPost, &postResp)
		} else {
			posts.ActivePost = append(posts.ActivePost, &postResp)
		}
	}

	return posts, nil

}

func (r *postRepo) GetPostCustomerId(req *pb.CustomerId) (*pb.ListPostCustomer, error) {
	posts := []*pb.PostReviewResponse{}
	rows, err := r.db.Query(`
	 select id,customer_id,name,description,created_at,updated_at from post  where customer_id=$1
	`, req.Id)
	if err != nil {
		return &pb.ListPostCustomer{}, err
	}
	for rows.Next() {
		postResp := pb.PostReviewResponse{}
		err = rows.Scan(&postResp.Id, &postResp.CustomerId, &postResp.Name, &postResp.Description, &postResp.CreatedAt, &postResp.UpdatedAt)
		if err != nil {
			return &pb.ListPostCustomer{}, err
		}
		medias, err := r.db.Query(
			`select post_id,name,link,type from media where post_id=$1`, postResp.Id,
		)
		if err != nil {
			return &pb.ListPostCustomer{}, err
		}
		for medias.Next() {
			mediasResp := pb.MediasResponse{}
			err = medias.Scan(&mediasResp.PostId,
				&mediasResp.Name,
				&mediasResp.Link,
				&mediasResp.Type)
			if err != nil {
				return &pb.ListPostCustomer{}, err
			}
			postResp.Media = append(postResp.Media, &mediasResp)

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
	_, err = r.db.Exec(`
	delete from media where post_id=$1
	`, req.Id)
	if err != nil {
		return &pb.PostResponse{}, err
	}
	for _, media := range req.Media {
		mediaResp := pb.MediasResponse{}
		err = r.db.QueryRow(`
		insert into media (post_id,name,link,type) values($1,$2,$3,$4)
		returning post_id,name,link,type
		`, req.Id, media.Name, media.Link, media.Type).Scan(&mediaResp.PostId,
			&mediaResp.Name, &mediaResp.Link, &mediaResp.Type)
		if err != nil {
			return &pb.PostResponse{}, err
		}
		postResp.Media = append(postResp.Media, &mediaResp)
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
func (r *postRepo) DeletePostByCustomerId(req *pb.CustomerId) (*pb.DeletedReview, error) {
	rows, err := r.db.Query(`
	update post set deleted_at=current_timestamp
	where customer_id=$1 returning id
	`, req.Id)
	if err != nil {
		return &pb.DeletedReview{}, err
	}
	arr := pb.DeletedReview{}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return &pb.DeletedReview{}, err
		}
		arr.ReviewsIds = append(arr.ReviewsIds, &pb.ReviewsIds{Id: id})

	}
	return &arr, nil
}
