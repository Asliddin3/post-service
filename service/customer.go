package service

import (
	"context"

	pb "github.com/Asliddin3/post-servise/genproto/customer"
	l "github.com/Asliddin3/post-servise/pkg/logger"
	grpcclient "github.com/Asliddin3/post-servise/service/grpc_client"
)

type CustomerService struct {
	client *grpcclient.ServiceManager
	logger l.Logger
}

func (r *CustomerService) GetPostCustomerId(ctx context.Context, req *pb.CustomerId) (*pb.CustomerResponse, error) {
	customerInfo, err := r.client.CustomerService().GetCustomerInfo(context.Background(), req)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	return customerInfo, nil
}
