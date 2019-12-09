package userssvc

import (
	"context"
	"fmt"

	pb "github.com/cagodoy/tenpo-services/protos/go"
	users "github.com/cagodoy/tenpo-users-api"
	"github.com/cagodoy/tenpo-users-api/database"
	"github.com/cagodoy/tenpo-users-api/service"
	"golang.org/x/crypto/bcrypt"
)

var _ pb.UserServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	usersSvc users.Service
}

// New ...
func New(store database.Store) *Service {
	return &Service{
		usersSvc: service.NewUsers(store),
	}
}

// Get Gets a user by ID.
func (us *Service) Get(ctx context.Context, gr *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	id := gr.GetUserId()
	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Get][Request] id = %v", id))

	user, err := us.usersSvc.GetByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Get][Error] %v", err))
		return &pb.UserGetResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	a := user.ToProto()

	res := &pb.UserGetResponse{
		Meta:  nil,
		Data:  a,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Get][Response] %v", res))
	return res, nil
}

// GetByEmail get a user by Email
func (us *Service) GetByEmail(ctx context.Context, gr *pb.UserGetByEmailRequest) (*pb.UserGetByEmailResponse, error) {
	email := gr.GetEmail()
	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][GetByEmail][Request] email = %v", email))

	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][GetByEmail][Error] %v", "email user params empty"))
		return &pb.UserGetByEmailResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "email user params empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][GetByEmail][Error] %v", "user not found"))
		return &pb.UserGetByEmailResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    404,
				Message: "user not found",
			},
		}, nil
	}

	res := &pb.UserGetByEmailResponse{
		Meta:  nil,
		Data:  user.ToProto(),
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][GetByEmail][Response] %v", res))
	return res, nil
}

// Create creates a new user into database.
func (us *Service) Create(ctx context.Context, gr *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Request] data = %v", gr.GetData()))

	email := gr.GetData().GetEmail()
	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Error] %v", "email user param is empty"))
		return &pb.UserCreateResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "email user param is empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		name := gr.GetData().GetName()
		if name == "" {
			fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Error] %v", "name user param is empty"))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    400,
					Message: "name user param is empty",
				},
			}, nil
		}

		password := gr.GetData().GetPassword()
		if password == "" {
			fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Error] %v", "password user params is empty"))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    400,
					Message: "password user params is empty",
				},
			}, nil
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Error] %v", err))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    500,
					Message: "could not generate hashed password",
				},
			}, nil
		}

		user := &users.User{
			Email:    email,
			Name:     name,
			Password: string(hashedPassword),
		}

		if err := us.usersSvc.Create(user); err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Error] %v", err))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    500,
					Message: err.Error(),
				},
			}, nil
		}

		res := &pb.UserCreateResponse{
			Meta:  nil,
			Data:  user.ToProto(),
			Error: nil,
		}
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Response] %v", res))
		return res, nil
	}

	res := &pb.UserCreateResponse{
		Meta:  nil,
		Data:  user.ToProto(),
		Error: nil,
	}
	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][Create][Response] %v", res))
	return res, nil
}

// VerifyPassword ...
func (us *Service) VerifyPassword(ctx context.Context, gr *pb.UserVerifyPasswordRequest) (*pb.UserVerifyPasswordResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][VerifyPassword][Request] email = %v password = %v", gr.GetEmail(), gr.GetPassword()))
	email := gr.GetEmail()
	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][VerifyPassword][Error] %v", "email user param is empty"))
		return &pb.UserVerifyPasswordResponse{
			Meta:  nil,
			Valid: false,
			Error: &pb.Error{
				Code:    400,
				Message: "email user param is empty",
			},
		}, nil
	}

	password := gr.GetPassword()
	if password == "" {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][VerifyPassword][Error] %v", "password user param is empty"))
		return &pb.UserVerifyPasswordResponse{
			Meta:  nil,
			Valid: false,
			Error: &pb.Error{
				Code:    400,
				Message: "password user param is empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][VerifyPassword][Error] %v", "user not found"))
		return &pb.UserVerifyPasswordResponse{
			Meta:  nil,
			Valid: false,
			Error: &pb.Error{
				Code:    404,
				Message: "user not found",
			},
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][VerifyPassword][Error] %v", "invalid password"))
		return &pb.UserVerifyPasswordResponse{
			Meta:  nil,
			Valid: false,
			Error: &pb.Error{
				Code:    403,
				Message: "invalid password",
			},
		}, nil
	}

	res := &pb.UserVerifyPasswordResponse{
		Meta:  nil,
		Valid: true,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][TenpoUsersService][VerifyPassword][Response] %v", res))
	return res, nil
}
