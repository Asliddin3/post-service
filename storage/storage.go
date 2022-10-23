package storage

import (
	"github.com/Asliddin3/post-servise/storage/postgres"
	"github.com/Asliddin3/post-servise/storage/repo"
	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Post() repo.PostStorageI
	Review() repo.ReviewStorageI
}

type storagePg struct {
	db         *sqlx.DB
	postRepo   repo.PostStorageI
	reviewRepo repo.ReviewStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		postRepo: postgres.NewPostRepo(db),
	}
}
func (s storagePg) Post() repo.PostStorageI {
	return s.postRepo
}

func (s storagePg) Review() repo.ReviewStorageI {
	return s.reviewRepo
}
