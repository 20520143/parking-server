package service

import (
	"context"
	"github.com/google/uuid"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
)

type ParkingLotService struct {
	repo repo.PGInterface
}

func NewParkingLotService(repo repo.PGInterface) ParkingLotInterface {
	return &ParkingLotService{repo: repo}
}

type ParkingLotInterface interface {
	CreateParkingLot(ctx context.Context, req model.ParkingLotReq) (*model.ParkingLot, error)
	GetListParkingLot(ctx context.Context, req model.ListParkingLotReq) (model.ListParkingLotRes, error)
	GetOneParkingLot(ctx context.Context, id uuid.UUID) (model.ParkingLot, error)
	UpdateParkingLot(ctx context.Context, req model.ParkingLotReq) (model.ParkingLot, error)
	DeleteParkingLot(ctx context.Context, id uuid.UUID) error
	GetListParkingLotCompany(ctx context.Context, req model.GetListParkingLotReq) (model.ListParkingLotRes, error)
	UpdateParkingLotV2(ctx context.Context, req model.UpdateParkingLotReq) (model.ParkingLot, error)
}

func (s *ParkingLotService) CreateParkingLot(ctx context.Context, req model.ParkingLotReq) (*model.ParkingLot, error) {
	ParkingLot := &model.ParkingLot{
		Name:        valid.String(req.Name),
		Description: valid.String(req.Description),
		Address:     valid.String(req.Address),
		StartTime:   valid.DayTime(req.StartTime),
		EndTime:     valid.DayTime(req.EndTime),
		Lat:         valid.Float64(req.Lat),
		Long:        valid.Float64(req.Long),
		CompanyID:   valid.UUID(req.CompanyID),
	}

	if err := s.repo.CreateParkingLot(ctx, ParkingLot); err != nil {
		return nil, err
	}
	return ParkingLot, nil
}

func (s *ParkingLotService) GetListParkingLot(ctx context.Context, req model.ListParkingLotReq) (model.ListParkingLotRes, error) {
	return s.repo.GetListParkingLot(ctx, req)
}

func (s *ParkingLotService) GetOneParkingLot(ctx context.Context, id uuid.UUID) (model.ParkingLot, error) {
	return s.repo.GetOneParkingLot(ctx, id)
}

func (s *ParkingLotService) UpdateParkingLot(ctx context.Context, req model.ParkingLotReq) (model.ParkingLot, error) {
	ParkingLot, err := s.repo.GetOneParkingLot(ctx, valid.UUID(req.ID))
	if err != nil {
		return ParkingLot, err
	}

	utils.Sync(req, &ParkingLot)
	if err := s.repo.UpdateParkingLot(ctx, &ParkingLot); err != nil {
		return ParkingLot, err
	}

	return ParkingLot, nil
}

func (s *ParkingLotService) DeleteParkingLot(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteParkingLot(ctx, id)
}

func (s *ParkingLotService) GetListParkingLotCompany(ctx context.Context, req model.GetListParkingLotReq) (res model.ListParkingLotRes, err error) {
	return s.repo.GetListParkingLotCompany(ctx, req)
}

func (s *ParkingLotService) UpdateParkingLotV2(ctx context.Context, req model.UpdateParkingLotReq) (model.ParkingLot, error) {
	ParkingLot, err := s.repo.GetOneParkingLot(ctx, *req.ID)
	if err != nil {
		return ParkingLot, err
	}

	parkingLot := model.ParkingLot{
		BaseModel:   model.BaseModel{ID: *req.ID},
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Lat:         req.Lat,
		Long:        req.Long,
		CompanyID:   ParkingLot.CompanyID,
	}

	var newBlocks []model.Block
	var newTimeFrames []model.TimeFrame

	for _, block := range req.Blocks {
		block.ParkingLotID = *req.ID

		if block.ID == uuid.Nil {
			newBlocks = append(newBlocks, block)
		}
		parkingLot.Blocks = append(parkingLot.Blocks, block)
	}

	for _, timeFrame := range req.TimeFrames {
		timeFrame.ParkingLotId = *req.ID

		if timeFrame.ID == uuid.Nil {
			newTimeFrames = append(newTimeFrames, timeFrame)
		}
		parkingLot.TimeFrames = append(parkingLot.TimeFrames, timeFrame)
	}

	if err := s.repo.UpdateParkingLotV2(ctx, parkingLot, newTimeFrames, newBlocks); err != nil {
		return ParkingLot, err
	}

	resp, err := s.repo.GetOneParkingLot(ctx, *req.ID)
	if err != nil {
		return ParkingLot, err
	}

	return resp, nil

}
