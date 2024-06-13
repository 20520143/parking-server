package handlers

import (
	"parking-server/pkg/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MigrationHandler struct {
	db *gorm.DB
}

func NewMigrationHandler(db *gorm.DB) *MigrationHandler {
	return &MigrationHandler{db: db}
}

func (h *MigrationHandler) Migrate(ctx *gin.Context) {
	_ = h.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	_ = h.db.Exec("CREATE EXTENSION IF NOT EXISTS \"postgis\"")
	_ = h.db.Exec("CREATE EXTENSION IF NOT EXISTS \"unaccent\"")

	models := []interface{}{
		model.Block{},
		model.Company{},
		model.Favorite{},
		model.LongTermTicket{},
		model.ParkingLot{},
		model.ParkingSlot{},
		model.RefreshToken{},
		model.Setting{},
		model.Ticket{},
		model.TicketExtend{},
		model.TimeFrame{},
		model.User{},
		model.Vehicle{},
		model.Employee{},
	}
	for _, m := range models {
		err := h.db.AutoMigrate(m)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
	}
}
