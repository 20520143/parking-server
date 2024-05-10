package route

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/service"
	"parking-server/conf"
	"parking-server/pkg/handlers"
	"parking-server/pkg/repo"
	service2 "parking-server/pkg/service"
)

type extraSetting struct {
	DbDebugEnable bool `env:"DB_DEBUG_ENABLE" envDefault:"true"`
}

type Service struct {
	*service.BaseApp
	setting *extraSetting
}

func NewService() *Service {
	s := &Service{
		service.NewApp("Parking", "v1.0"),
		&extraSetting{},
	}
	// repo
	_ = env.Parse(s.setting)
	s.Config.DB.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s connect_timeout=5",
		conf.GetConfig().DBHost,
		conf.GetConfig().DBPort,
		conf.GetConfig().DBUser,
		conf.GetConfig().DBName,
		conf.GetConfig().DBPass,
	)
	db := s.GetDB()
	if s.setting.DbDebugEnable {
		db = db.Debug()
	}
	repoPG := repo.NewPGRepo(db)

	//service
	authService := service2.NewAuthService(repoPG)
	userService := service2.NewUserService(repoPG)
	lotService := service2.NewParkingLotService(repoPG)
	blockService := service2.NewBlockService(repoPG)
	vehicleService := service2.NewVehicleService(repoPG)
	timeFrameService := service2.NewTimeFrameService(repoPG)
	ticketService := service2.NewTicketService(repoPG)
	companyService := service2.NewCompanyService(repoPG)

	//handler
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	lotHandler := handlers.NewParkingLotHandler(lotService)
	blockHandler := handlers.NewBlockHandler(blockService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	timeFrameHandler := handlers.NewTimeFrameHandler(timeFrameService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	companyHanler := handlers.NewCompanyHandler(companyService)

	route := s.Router
	route.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		}
	}(),
	)

	v1Api := s.Router.Group("/api/v1")
	v2Api := s.Router.Group("/api/v2")
	merchantApi := s.Router.Group("/api/merchant")
	swaggerApi := s.Router.Group("/")

	// swagger
	swaggerApi.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	//auth
	v1Api.POST("/user/login", ginext.WrapHandler(authHandler.Login))
	v1Api.POST("/user/reset-password", ginext.WrapHandler(authHandler.ResetPassword))
	v1Api.POST("/user/send-otp", ginext.WrapHandler(authHandler.SendOtp))
	v1Api.POST("/user/verify-otp", ginext.WrapHandler(authHandler.VerifyOtp))

	// user
	v1Api.POST("/user/check-phone", ginext.WrapHandler(userHandler.CheckDuplicatePhone))

	// parking lot
	v1Api.POST("/parking-lot/create", ginext.WrapHandler(lotHandler.CreateParkingLot))
	v1Api.GET("/parking-lot/get-one/:id", ginext.WrapHandler(lotHandler.GetOneParkingLot))
	v1Api.GET("/parking-lot/get-list", ginext.WrapHandler(lotHandler.GetListParkingLot))
	v1Api.PUT("/parking-lot/update/:id", ginext.WrapHandler(lotHandler.UpdateParkingLot))
	v1Api.DELETE("/parking-lot/delete/:id", ginext.WrapHandler(lotHandler.DeleteParkingLot))

	v2Api.PUT("/parking-lot/update", ginext.WrapHandler(lotHandler.UpdateParkingLotV2))
	// block
	v1Api.POST("/block/create", ginext.WrapHandler(blockHandler.CreateBlock))
	v1Api.GET("/block/get-one/:id", ginext.WrapHandler(blockHandler.GetOneBlock))
	v1Api.GET("/block/get-list", ginext.WrapHandler(blockHandler.GetListBlock))
	v1Api.PUT("/block/update/:id", ginext.WrapHandler(blockHandler.UpdateBlock))
	v1Api.DELETE("/block/delete/:id", ginext.WrapHandler(blockHandler.DeleteBlock))

	//time frame
	v1Api.POST("/time-frame/create", ginext.WrapHandler(timeFrameHandler.CreateTimeFrame))

	// parking lot
	v1Api.POST("/vehicle/create", ginext.WrapHandler(vehicleHandler.CreateVehicle))
	v1Api.GET("/vehicle/get-one/:id", ginext.WrapHandler(vehicleHandler.GetOneVehicle))
	v1Api.GET("/vehicle/get-list", ginext.WrapHandler(vehicleHandler.GetListVehicle))
	v1Api.PUT("/vehicle/update/:id", ginext.WrapHandler(vehicleHandler.UpdateVehicle))
	v1Api.DELETE("/vehicle/delete/:id", ginext.WrapHandler(vehicleHandler.DeleteVehicle))

	// company
	merchantApi.POST("/company/create", cors.Default(), ginext.WrapHandler(companyHanler.CreateCompany))
	merchantApi.PUT("/company/update/:id", cors.Default(), ginext.WrapHandler(companyHanler.UpdateCompany))
	merchantApi.POST("/company/login", cors.Default(), ginext.WrapHandler(companyHanler.Login))
	merchantApi.GET("/company/get-one/:id", cors.Default(), ginext.WrapHandler(companyHanler.GetOneCompany))
	merchantApi.PUT("/company/update-password/:id", cors.Default(), ginext.WrapHandler(companyHanler.UpdateCompanyPassword))

	merchantApi.GET("/parking-lot/get-list", ginext.WrapHandler(lotHandler.GetListParkingLotCompany))
	merchantApi.GET("/parking-lot/get-one/:id", ginext.WrapHandler(lotHandler.GetOneParkingLot))

	merchantApi.GET("/block/get-list", ginext.WrapHandler(blockHandler.GetListBlock))

	merchantApi.GET("/time-frame/get-list", ginext.WrapHandler(timeFrameHandler.GetAllTimeFrame))
	merchantApi.GET("/ticket/get-all", ginext.WrapHandler(ticketHandler.GetAllTicketCompany))
	// Migrate
	migrateHandler := handlers.NewMigrationHandler(db)
	s.Router.POST("/internal/migrate", migrateHandler.Migrate)
	return s
}
