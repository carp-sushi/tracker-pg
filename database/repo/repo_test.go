package repo_test

import (
	"testing"

	"github.com/carp-cobain/tracker-pg/database"
	"github.com/carp-cobain/tracker-pg/database/model"
	"github.com/carp-cobain/tracker-pg/database/repo"
	"github.com/carp-cobain/tracker-pg/domain"

	"gorm.io/gorm"
)

func createTestDB(t *testing.T) *gorm.DB {
	db, err := database.Connect("sqlite3", "file::memory:?cache=shared", 1)
	if err != nil {
		t.Fatalf("failed to connect to test database: %+v", err)
	}
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations on test database: %+v", err)
	}
	return db
}

func TestCampaignRepo(t *testing.T) {
	// Setup
	db := createTestDB(t)
	account := domain.MustValidateAccount("tp1mrzpjszjs6dc5e8fwy23trnz775rwqvhpzzzz2")
	campaignRepo := repo.NewCampaignRepo(db)

	// Create
	campaign, err := campaignRepo.CreateCampaign(account, "Campaign Unit Testing")
	if err != nil {
		t.Fatalf("failed to create campaign: %+v", err)
	}
	if campaign.Type != "referral" {
		t.Fatalf("expected campaign type: referral, got: %s", campaign.Type)
	}

	// Read
	if _, err := campaignRepo.GetCampaign(campaign.ID); err != nil {
		t.Fatalf("failed to get campaign: %+v", err)
	}
	params := domain.DefaultPageParams()
	if page := campaignRepo.GetCampaigns(account, params); len(page.Data) != 1 {
		t.Fatalf("got unexpected number of campaigns")
	}
}

func TestReferralRepo(t *testing.T) {
	// Setup
	db := createTestDB(t)
	campgaignRepo := repo.NewCampaignRepo(db)
	referralRepo := repo.NewReferralRepo(db)

	// Base campaign
	referer := domain.MustValidateAccount("tp1mrzpjszjs6dc5e8fwy23trnz775rwqvhpzzzz3")
	campaign, err := campgaignRepo.CreateCampaign(referer, "Referral Unit Testing")
	if err != nil {
		t.Fatalf("failed to create referral campaign: %+v", err)
	}

	// Create
	referee := domain.MustValidateAccount("tp1mrzpjszjs6dc5e8fwy23trnz775rwqvhpzzzz4")
	referral, err := referralRepo.CreateReferral(campaign.ID, referee)
	if err != nil {
		t.Fatalf("failed to create referral: %+v", err)
	}

	// Read
	params := domain.DefaultPageParams()
	if page := referralRepo.GetReferrals(campaign.ID, params); len(page.Data) != 1 {
		t.Fatalf("got unexpected number of referrals for campaign")
	}

	// Update (set status)
	verified := model.ReferralStatusVerified.ToDomain()
	updated, err := referralRepo.UpdateReferral(referral.ID, verified)
	if err != nil {
		t.Fatalf("failed to update referral status: %+v", err)
	}
	if updated.Status != verified {
		t.Fatalf("expected verified status, got: %s", updated.Status)
	}

	// Error check: ensure accounts can't add referrals for thier own campaigns.
	if _, err := referralRepo.CreateReferral(campaign.ID, referer); err == nil {
		t.Fatalf("expected self referral error")
	}
}
