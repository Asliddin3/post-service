package kafka

import (
	"fmt"

	"github.com/Asliddin3/post-servise/config"
	pb "github.com/Asliddin3/post-servise/genproto/post"
	"github.com/Asliddin3/post-servise/pkg/logger"

	"github.com/Asliddin3/post-servise/storage"
)

type KafkaHandler struct {
	config  config.Config
	storage storage.IStorage
	log     logger.Logger
}

func NewKafkaHandlerFunc(config config.Config, storage storage.IStorage, log logger.Logger) *KafkaHandler {
	return &KafkaHandler{
		storage: storage,
		config:  config,
		log:     log,
	}
}

func (h *KafkaHandler) Handle(value []byte) error {
	post := &pb.CustomerResponse{}

	err := post.Unmarshal(value)
	if err != nil {
		return err
	}
	resp, err := h.storage.Post().CreateCustomer(post)
	fmt.Println(resp)
	if err != nil {
		return err
	}
	return nil
}
