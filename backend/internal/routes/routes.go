package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/handlers"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/repositories"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")

	plotRepository := repositories.NewPlotRepository(db)
	buyerRepository := repositories.NewBuyerRepository(db)
	saleRepository := repositories.NewSaleRepository(db)
	maintenanceRepository := repositories.NewMaintenanceRepository(db)

	plotService := services.NewPlotService(plotRepository, saleRepository, maintenanceRepository)
	buyerService := services.NewBuyerService(buyerRepository)
	saleService := services.NewSaleService(saleRepository)
	maintenanceService := services.NewMaintenanceService(maintenanceRepository)
	dashboardService := services.NewDashboardService(saleRepository, maintenanceRepository)

	plotHandler := handlers.NewPlotHandler(plotService)
	buyerHandler := handlers.NewBuyerHandler(buyerService)
	saleHandler := handlers.NewSaleHandler(saleService)
	maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	plotHandler.RegisterRoutes(api)
	buyerHandler.RegisterRoutes(api)
	saleHandler.RegisterRoutes(api)
	maintenanceHandler.RegisterRoutes(api)
	dashboardHandler.RegisterRoutes(api)
}
