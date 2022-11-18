package grpcClient

import (
	"fmt"

	"github.com/Asliddin3/post-servise/config"
	customerPB "github.com/Asliddin3/post-servise/genproto/customer"
	reviewPB "github.com/Asliddin3/post-servise/genproto/review"
	"google.golang.org/grpc"
)

//GrpcClientI ...
type ServiceManager struct {
	conf            config.Config
	reviewService   reviewPB.ReviewServiceClient
	customerService customerPB.CustomerServiceClient
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
	connCustomer, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cnfg.CustomerSerivceHost, cnfg.CustomerSerivcePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial product service: host: %s and port: %d",
			cnfg.CustomerSerivceHost, cnfg.CustomerSerivcePort)
	}

	serviceManager := &ServiceManager{
		conf:            cnfg,
		reviewService:   reviewPB.NewReviewServiceClient(connReview),
		customerService: customerPB.NewCustomerServiceClient(connCustomer),
	}

	return serviceManager, nil
}

func (s *ServiceManager) CustomerService() customerPB.CustomerServiceClient {
	return s.customerService
}

func (s *ServiceManager) ReviewService() reviewPB.ReviewServiceClient {
	return s.reviewService
}
