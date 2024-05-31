package service

import (
	"context"
	"github.com/google/uuid"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/utils"
)

type TimeFrameService struct {
	repo repo.PGInterface
}

func NewTimeFrameService(repo repo.PGInterface) TimeFrameServiceInterface {
	return &TimeFrameService{repo: repo}
}

type TimeFrameServiceInterface interface {
	GetAllTimeFrame(ctx context.Context, req model.GetListTimeFrameParam) (res model.ListTimeFrame, err error)
	CreateMultiTimeFrame(ctx context.Context, req model.ListTimeFrameReq) (err error)
	UpdateMultiTimeFrame(ctx context.Context, req model.ListTimeFrameReq) (err error)

	CreateTimeFrame(ctx context.Context, req model.TimeFrameReq) (*model.TimeFrame, error)
	GetOneTimeFrame(ctx context.Context, id uuid.UUID) (model.TimeFrame, error)
	UpdateTimeFrame(ctx context.Context, id uuid.UUID, req model.TimeFrameRequest) (model.TimeFrame, error)
	DeleteTimeFrame(ctx context.Context, id uuid.UUID) error
}

func (s *TimeFrameService) GetAllTimeFrame(ctx context.Context, req model.GetListTimeFrameParam) (model.ListTimeFrame, error) {
	res, err := s.repo.GetAllTimeFrame(ctx, req, nil)
	if err != nil {
		return model.ListTimeFrame{}, err
	}
	return *res, nil
}

func (s *TimeFrameService) CreateTimeFrame(ctx context.Context, req model.TimeFrameReq) (*model.TimeFrame, error) {
	time := &model.TimeFrame{Duration: req.Duration, Cost: req.Cost, ParkingLotId: req.ParkingLotId}

	if err := s.repo.CreateTimeframe(ctx, time); err != nil {
		return nil, err
	}
	return time, nil
}

func (s *TimeFrameService) GetOneTimeFrame(ctx context.Context, id uuid.UUID) (model.TimeFrame, error) {
	return s.repo.GetOneTimeframe(ctx, id)
}

func (s *TimeFrameService) UpdateTimeFrame(ctx context.Context, id uuid.UUID, req model.TimeFrameRequest) (model.TimeFrame, error) {
	time, err := s.repo.GetOneTimeframe(ctx, id)
	if err != nil {
		return time, err
	}

	utils.Sync(req, &time)
	if err := s.repo.UpdateTimeframe(ctx, &time); err != nil {
		return time, err
	}

	return time, nil
}

func (s *TimeFrameService) DeleteTimeFrame(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTimeframe(ctx, id)
}

func (s *TimeFrameService) CreateMultiTimeFrame(ctx context.Context, req model.ListTimeFrameReq) (err error) {
	listUser := []model.TimeFrame{}
	for _, item := range req.Data {
		tmp := model.TimeFrame{}
		utils.Sync(item, &tmp)
		listUser = append(listUser, tmp)
	}
	err = s.repo.CreateMultiTimeFrame(ctx, listUser, nil)
	return err
}

func (s *TimeFrameService) UpdateMultiTimeFrame(ctx context.Context, req model.ListTimeFrameReq) (err error) {
	//detele all time fram by parking lot
	err = s.repo.DeleteTimeFrameByParkingLotID(ctx, req.Data[0].ParkingLotId.String(), nil)
	if err != nil {
		return err
	}
	err = s.CreateMultiTimeFrame(ctx, req)
	return err
}
