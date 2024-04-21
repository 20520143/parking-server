package service

import (
	"context"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
)

type TimeFrameService struct {
	repo repo.PGInterface
}

func NewTimeFrameService(repo repo.PGInterface) TimeFrameServiceInterface {
	return &TimeFrameService{repo: repo}
}

type TimeFrameServiceInterface interface {
	GetAllTimeFrame(ctx context.Context, req model.GetListTimeFrameParam) (res model.ListTimeFrame, err error)
}

func (s *TimeFrameService) GetAllTimeFrame(ctx context.Context, req model.GetListTimeFrameParam) (model.ListTimeFrame, error) {
	res, err := s.repo.GetAllTimeFrame(ctx, req, nil)
	if err != nil {
		return model.ListTimeFrame{}, err
	}
	return *res, nil
}
