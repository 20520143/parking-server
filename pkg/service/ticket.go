package service

import (
	"context"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/valid"
	"time"
)

type TicketService struct {
	repo repo.PGInterface
}

func NewTicketService(repo repo.PGInterface) TicketServiceInterface {
	return &TicketService{repo: repo}
}

type TicketServiceInterface interface {
	CreateTicket(ctx context.Context, req *model.TicketReq) (*model.Ticket, error)
	ProcedureWithTicket(ctx context.Context, req *model.ProcedureReq) (bool, error)
	ExtendTicket(ctx context.Context, req *model.ExtendTicketReq) (*model.TicketExtend, error)
	GetAllTicket(ctx context.Context, req model.GetListTicketParam) ([]model.Ticket, error)
	GetOneTicketWithExtend(ctx context.Context, id string) (model.TicketResponse, error)
	CancelTicket(ctx context.Context, id string) error
	GetAllTicketCompany(ctx context.Context, req model.GetListTicketReq) ([]model.GetListTicketRes, error)
}

func (s *TicketService) GetAllTicketCompany(ctx context.Context, req model.GetListTicketReq) ([]model.GetListTicketRes, error) {
	return s.repo.GetAllTicketCompany(ctx, req)
}

func (s *TicketService) CreateTicket(ctx context.Context, req *model.TicketReq) (*model.Ticket, error) {
	ticket := &model.Ticket{
		BaseModel: model.BaseModel{
			CreatorID: req.UserId,
			UpdaterID: req.UserId,
		},
		UserId:        req.UserId,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		VehicleId:     req.VehicleId,
		ParkingLotId:  req.ParkingLotId,
		ParkingSlotId: req.ParkingSlotId,
		TimeFrameId:   req.TimeFrameId,
		State:         "new",
		Total:         valid.Float64(req.Total),
	}

	if req.IsLongTerm {
		longTermTicket := &model.LongTermTicket{
			BaseModel: model.BaseModel{
				CreatorID: req.UserId,
				UpdaterID: req.UserId,
			},
			StartTime:     req.StartTime,
			EndTime:       req.EndTime,
			VehicleId:     req.VehicleId,
			ParkingLotId:  req.ParkingLotId,
			ParkingSlotId: req.ParkingSlotId,
			TimeFrameId:   req.TimeFrameId,
		}
		s.repo.CreateLongTermTicket(ctx, longTermTicket, nil)

		//create ticket normal
		switch req.Type {
		case "DAILY":
		case "CYCLE":
		case "CUSTOM":
		}
	}
	err := s.repo.CreateTicket(ctx, ticket, nil)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
func (s *TicketService) ExtendTicket(ctx context.Context, req *model.ExtendTicketReq) (*model.TicketExtend, error) {
	ticket, err := s.repo.GetOneTicket(ctx, valid.UUID(req.TicketOriginId).String(), nil)
	if err != nil {
		return nil, err
	}
	extendTicket := &model.Ticket{
		BaseModel: model.BaseModel{
			CreatorID: ticket.CreatorID,
			UpdaterID: ticket.UpdaterID,
		},
		UserId:        ticket.UserId,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		VehicleId:     ticket.VehicleId,
		ParkingLotId:  ticket.ParkingLotId,
		ParkingSlotId: ticket.ParkingSlotId,
		TimeFrameId:   req.TimeFrameId,
		State:         "extend",
		Total:         valid.Float64(req.Total),
	}
	ticket.IsExtend = true
	if err := s.repo.UpdateTicket(ctx, &ticket, nil); err != nil {
		return nil, err
	}
	if err := s.repo.CreateTicket(ctx, extendTicket, nil); err != nil {
		return nil, err
	}
	//create extend ticket table
	ticketEx := &model.TicketExtend{
		TicketExtendId: extendTicket.ID,
		TicketId:       ticket.ID,
	}
	if err := s.repo.CreateTicketExtend(ctx, ticketEx, nil); err != nil {
		return nil, err
	}
	return ticketEx, nil
}
func (s *TicketService) GetAllTicket(ctx context.Context, req model.GetListTicketParam) ([]model.Ticket, error) {
	res, err := s.repo.GetAllTicket(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *TicketService) GetOneTicketWithExtend(ctx context.Context, id string) (model.TicketResponse, error) {
	ticket, err := s.repo.GetOneTicketWithExtend(ctx, id, nil)
	if err != nil {
		return model.TicketResponse{}, err
	}
	ticketExtend, err := s.repo.GetListExtendTicketByOrigin(ctx, ticket.ID.String(), nil)
	if err != nil {
		return model.TicketResponse{}, err
	}
	ticketRes := model.TicketResponse{
		Ticket:       ticket,
		TicketExtend: ticketExtend,
	}
	return ticketRes, nil
}
func (s *TicketService) CancelTicket(ctx context.Context, id string) error {
	ticket, err := s.repo.GetOneTicket(ctx, id, nil)
	if err != nil {
		return err
	}
	ticket.State = "cancel"
	if err := s.repo.UpdateTicket(ctx, &ticket, nil); err != nil {
		return err
	}
	return nil
}
func (s *TicketService) ProcedureWithTicket(ctx context.Context, req *model.ProcedureReq) (bool, error) {
	ticket, err := s.repo.GetOneTicket(ctx, req.TicketId, nil)
	if err != nil {
		return false, err
	}
	switch req.Type {
	case "check_in":
		ticket.State = "ongoing"
		ticket.EntryTime = valid.DayTimePointer(time.Now())
	case "check_out":
		ticket.State = "completed"
		ticket.ExitTime = valid.DayTimePointer(time.Now())
	}
	if err := s.repo.UpdateTicket(ctx, &ticket, nil); err != nil {
		return false, err
	}
	return true, nil
}
