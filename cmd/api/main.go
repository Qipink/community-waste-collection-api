package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"community-waste-collection-api/internal/config"
	"community-waste-collection-api/internal/handler"
	"community-waste-collection-api/internal/repository"
	"community-waste-collection-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db := config.InitDB()

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	maxUploadMB := 5
	if v := os.Getenv("MAX_UPLOAD_MB"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxUploadMB = n
		}
	}
	maxUploadBytes := int64(maxUploadMB) * 1024 * 1024
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("gagal membuat folder upload: %v", err)
	}

	householdRepo := repository.NewHouseholdRepository(db)
	pickupRepo := repository.NewPickupRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	reportRepo := repository.NewReportRepository(db)

	householdService := service.NewHouseholdService(householdRepo)
	pickupService := service.NewPickupService(pickupRepo, paymentRepo, householdRepo)
	paymentService := service.NewPaymentService(paymentRepo, householdRepo, pickupRepo, uploadDir)
	reportService := service.NewReportService(reportRepo)

	householdHandler := handler.NewHouseholdHandler(householdService)
	pickupHandler := handler.NewPickupHandler(pickupService)
	paymentHandler := handler.NewPaymentHandler(paymentService, maxUploadBytes)
	reportHandler := handler.NewReportHandler(reportService)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/uploads", uploadDir)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok"})
	})

	api := r.Group("/api")
	api.POST("/households", householdHandler.Create)
	api.GET("/households", householdHandler.GetAll)
	api.GET("/households/:id", householdHandler.GetByID)
	api.DELETE("/households/:id", householdHandler.Delete)

	api.POST("/pickups", pickupHandler.Create)
	api.GET("/pickups", pickupHandler.List)
	api.PUT("/pickups/:id/schedule", pickupHandler.Schedule)
	api.PUT("/pickups/:id/complete", pickupHandler.Complete)
	api.PUT("/pickups/:id/cancel", pickupHandler.Cancel)

	api.POST("/payments", paymentHandler.Create)
	api.GET("/payments", paymentHandler.List)
	api.PUT("/payments/:id/confirm", paymentHandler.Confirm)

	api.GET("/reports/waste-summary", reportHandler.WasteSummary)
	api.GET("/reports/payment-summary", reportHandler.PaymentSummary)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Printf("server berjalan di port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gagal menjalankan server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("menerima sinyal shutdown, menutup server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown gagal: %v", err)
	}
	log.Println("server berhenti dengan baik")
}
