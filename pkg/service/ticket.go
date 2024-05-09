package service

import (
	"context"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
)

type TicketService struct {
	repo repo.PGInterface
}

func NewTicketService(repo repo.PGInterface) TicketServiceInterface {
	return &TicketService{repo: repo}
}

type TicketServiceInterface interface {
	GetAllTicketCompany(ctx context.Context, req model.GetListTicketReq) ([]model.GetListTicketRes, error)
}

func (s *TicketService) GetAllTicketCompany(ctx context.Context, req model.GetListTicketReq) ([]model.GetListTicketRes, error) {
	return s.repo.GetAllTicketCompany(ctx, req)
}
