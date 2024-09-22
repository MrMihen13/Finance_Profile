package transport

import (
	"context"
	"errors"
	_grpc "github.com/MrMihen13/finance-protos/gen/go/profile"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log/slog"
	"profile/internal/models"
	"profile/internal/pkg/validator"
)

type service interface {
	Create(email string) (*models.Profile, error)
	GetByID(id uuid.UUID) (*models.Profile, error)
	UpdateEmail(id uuid.UUID, newEmail string) (*models.Profile, error)
	DeleteByID(id uuid.UUID) error
}

type Server struct {
	_grpc.UnimplementedProfileServer
	log     *slog.Logger
	profile service
}

func NewServer(log *slog.Logger, srv service) *Server {
	return &Server{log: log, profile: srv}
}

func (s *Server) Register(gRPC *grpc.Server) {
	_grpc.RegisterProfileServer(gRPC, s)
}

func (s *Server) Create(_ context.Context, req *_grpc.RegisterRequest) (*_grpc.ProfileItem, error) {
	email := req.GetEmail()

	if !validator.IsEmailValid(email) {
		return nil, errors.New("invalid email")
	}

	profile, err := s.profile.Create(email)
	if err != nil {
		return nil, err
	}
	return &_grpc.ProfileItem{
		Id:    profile.ID.String(),
		Email: profile.Email,
	}, nil
}

func (s *Server) Get(_ context.Context, req *_grpc.GetRequest) (*_grpc.ProfileItem, error) {
	rawID := req.GetId()
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, err
	}
	profile, err := s.profile.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &_grpc.ProfileItem{
		Id:    profile.ID.String(),
		Email: profile.Email,
	}, nil
}

func (s *Server) Update(_ context.Context, req *_grpc.UpdateRequest) (*_grpc.ProfileItem, error) {
	newEmail := req.GetNewEmail()

	if !validator.IsEmailValid(newEmail) {
		return nil, errors.New("invalid email")
	}

	rawId := req.GetId()

	id, err := uuid.Parse(rawId)
	if err != nil {
		return nil, err
	}

	profile, err := s.profile.UpdateEmail(id, newEmail)
	if err != nil {
		return nil, err
	}

	return &_grpc.ProfileItem{
		Id:    profile.ID.String(),
		Email: profile.Email,
	}, nil
}

func (s *Server) Delete(_ context.Context, req *_grpc.DeleteRequest) (*_grpc.DeleteResponse, error) {
	rawID := req.GetId()

	id, err := uuid.Parse(rawID)
	if err != nil {
		return &_grpc.DeleteResponse{Status: _grpc.StatusType_FAILED}, err
	}

	if err := s.profile.DeleteByID(id); err != nil {
		return &_grpc.DeleteResponse{Status: _grpc.StatusType_FAILED}, err
	}

	return &_grpc.DeleteResponse{Status: _grpc.StatusType_SUCCESS}, nil
}
