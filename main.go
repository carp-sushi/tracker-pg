package main

import (
	"log"
	"os"

	"github.com/carp-sushi/tracker-pg/database"
	"github.com/carp-sushi/tracker-pg/database/repo"
	"github.com/carp-sushi/tracker-pg/processor"
	"github.com/carp-sushi/tracker-pg/web/handler"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	// Perform env checks
	if _, ok := os.LookupEnv("DISABLE_COLOR"); ok {
		gin.DisableConsoleColor()
	}

	// DB
	db, err := database.ConnectAndMigrate()
	if err != nil {
		log.Panicf("unable to connnect to db: %+v", err)
	}

	// Repos
	campaignRepo := repo.NewCampaignRepo(db)
	referralRepo := repo.NewReferralRepo(db)

	// Handlers
	campaignHandler := handler.NewCampaignHandler(campaignRepo)
	referralHandler := handler.NewReferralHandler(campaignRepo, referralRepo)

	// API
	r := gin.Default()
	v1 := r.Group("/tracker/api/v1")
	{
		v1.GET("/campaigns", campaignHandler.GetCampaigns)
		v1.POST("/campaigns", campaignHandler.CreateCampaign)
		v1.GET("/campaigns/:id", campaignHandler.GetCampaign)
		v1.DELETE("/campaigns/:id", campaignHandler.ExpireCampaign)
		v1.PATCH("/campaigns/:id", campaignHandler.UpdateCampaign)
		v1.GET("/campaigns/:id/referrals", referralHandler.GetReferrals)
		v1.POST("/campaigns/:id/referrals", referralHandler.CreateReferral)
	}

	// Processors
	referralVerifier := processor.NewReferralVerifier(referralRepo, 350, 0)
	referralPayer := processor.NewReferralPayer(referralRepo, 700, 0)

	// Processor scheduling
	c := cron.New()
	c.AddFunc("*/30 * * * *", referralVerifier.VerifyReferrals) // Run every 30th minute
	c.AddFunc("@hourly", referralPayer.PayVerifiedReferrals)    // Run once an hour, beginning of hour
	c.Start()

	// Run server
	if err := r.Run(); err != nil {
		log.Panicf("unable to start tracker server:  %+v", err)
	}
}
