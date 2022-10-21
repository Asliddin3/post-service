package grpcClient

import (
	"fmt"

	"github.com/Asliddin3/post-servise/config"
	reviewPB "github.com/Asliddin3/post-servise/genproto/review"
	"google.golang.org/grpc"
)

//GrpcClientI ...
type ServiceManager struct {
	conf           config.Config
	reviewServisce reviewPB.ReviewServiceClient
}

func New(cnfg config.Config) (*ServiceManager, error) {
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cnfg.ReviewServiceHost, cnfg.ReviewServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial product service: host: %s and port: %d",
			cnfg.ReviewServiceHost, cnfg.ReviewServicePort)
	}

	serviceManager := &ServiceManager{
		conf:           cnfg,
		reviewServisce: reviewPB.NewReviewServiceClient(connReview),
	}

	return serviceManager, nil
}

func (s *ServiceManager) ReviewService() reviewPB.ReviewServiceClient {
	return s.reviewServisce
}
