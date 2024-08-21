package service

import (
	"api-geteway/config"
	"api-geteway/generated/mainservice"
	"api-geteway/generated/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManager interface {
	UserService() user.UserServiceClient
	ProductService() mainservice.MainServiceClient
}

type serviceManagerImpl struct {
	userClient    user.UserServiceClient
	productClient mainservice.MainServiceClient
}

func (s *serviceManagerImpl) UserService() user.UserServiceClient {
	return s.userClient
}

func (s *serviceManagerImpl) ProductService() mainservice.MainServiceClient {
	return s.productClient
}

func NewServiceManager(cfg config.Config) (ServiceManager, error) {
	connUser, err := grpc.NewClient(
		cfg.GRPC_USER_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	connProduct, err := grpc.NewClient(
		cfg.GRPC_PRODUCT_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceManagerImpl{
		userClient:    user.NewUserServiceClient(connUser),
		productClient: mainservice.NewMainServiceClient(connProduct),
	}, nil
}
