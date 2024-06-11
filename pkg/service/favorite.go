package service

import (
	"context"
	"github.com/google/uuid"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/valid"
)

type FavoriteService struct {
	repo repo.PGInterface
}

func NewFavoriteService(repo repo.PGInterface) FavoriteServiceInterface {
	return &FavoriteService{repo: repo}
}

type FavoriteServiceInterface interface {
	GetAllFavoriteParkingByUser(ctx context.Context, userId string) (res []model.Favorite, err error)
	Create(ctx context.Context, req model.FavoriteRequest) (res *model.Favorite, err error)
	GetOne(ctx context.Context, req model.FavoriteRequestV2) (res model.Favorite, err error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
}

func (s *FavoriteService) Create(ctx context.Context, req model.FavoriteRequest) (res *model.Favorite, err error) {
	favorite := &model.Favorite{
		UserId:       req.UserId,
		ParkingLotId: req.ParkingLotId,
		BaseModel: model.BaseModel{
			CreatorID: valid.UUIDPointer(req.UserId),
			UpdaterID: valid.UUIDPointer(req.UserId),
		},
	}
	if err := s.repo.CreateFavorite(ctx, favorite, nil); err != nil {
		return nil, err
	}
	return favorite, nil
}
func (s *FavoriteService) GetAllFavoriteParkingByUser(ctx context.Context, userId string) (res []model.Favorite, err error) {
	rs, err := s.repo.GetAllFavoriteParkingByUser(ctx, userId, nil)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *FavoriteService) DeleteOne(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteOneFavorite(ctx, id, nil); err != nil {
		return err
	}
	return nil
}
func (s *FavoriteService) GetOne(ctx context.Context, req model.FavoriteRequestV2) (model.Favorite, error) {
	res, err := s.repo.GetOne(ctx, req, nil)
	if err != nil {
		return res, err
	}
	return res, nil
}
