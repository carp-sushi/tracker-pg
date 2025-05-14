package repo

import (
	"fmt"

	"github.com/carp-sushi/tracker-pg/database/model"
	"github.com/carp-sushi/tracker-pg/database/query"
	"github.com/carp-sushi/tracker-pg/domain"

	"gorm.io/gorm"
)

// ReferralRepo manages referrals for campaigns.
type ReferralRepo struct {
	db *gorm.DB
}

// NewReferralRepo creates a new repository for managing referrals for campaigns.
func NewReferralRepo(db *gorm.DB) ReferralRepo {
	return ReferralRepo{db}
}

// GetReferrals gets a page of referrals for a campaign.
func (self ReferralRepo) GetReferrals(
	campaignID domain.CampaignID, pageParams domain.PageParams) domain.Page[domain.Referral] {

	var nextCursor uint64
	entities := query.SelectReferrals(self.db, campaignID.String(), pageParams.Cursor, pageParams.Limit)
	referrals := make([]domain.Referral, len(entities))
	for i, entity := range entities {
		referrals[i] = entity.ToDomain()
		nextCursor = max(nextCursor, uint64(entity.CreatedAt))
	}
	return domain.NewPage(nextCursor, pageParams.Limit, referrals)
}

// GetReferralsWithStatus gets a page of referrals with a given status.
func (self ReferralRepo) GetReferralsWithStatus(
	status domain.ReferralStatus, pageParams domain.PageParams) domain.Page[domain.Referral] {

	var nextCursor uint64
	cursor, limit := pageParams.Cursor, pageParams.Limit
	entities := query.SelectReferralsWithStatus(self.db, status, cursor, limit)
	referrals := make([]domain.Referral, len(entities))
	for i, entity := range entities {
		referrals[i] = entity.ToDomain()
		nextCursor = max(nextCursor, uint64(entity.CreatedAt))
	}
	return domain.NewPage(nextCursor, pageParams.Limit, referrals)
}

// CreateReferral creates a referral for a campaign.
func (self ReferralRepo) CreateReferral(
	campaignID domain.CampaignID, account domain.Account) (referral domain.Referral, err error) {

	var campaignEntity model.Campaign
	campaignEntity, err = query.SelectCampaign(self.db, campaignID.String())
	if err != nil {
		err = fmt.Errorf("campaign %s: %s", campaignID, err.Error())
		return
	}
	if campaignEntity.Type != model.CampaignTypeReferral {
		err = fmt.Errorf("error: non-referral campaign: %s", campaignEntity.Type.ToDomain())
		return
	}
	if campaignEntity.Account == account.String() {
		err = fmt.Errorf("self referral error: %s", account)
		return
	}
	var referralEntity model.Referral
	referralEntity, err = query.InsertReferral(self.db, campaignID.String(), account.String())
	if err == nil {
		referral = referralEntity.ToDomain()
	}
	return
}

// UpdateReferral updates the status of a referral for a campaign.
func (self ReferralRepo) UpdateReferral(
	referralID domain.ReferralID, status domain.ReferralStatus) (referral domain.Referral, err error) {

	var entity model.Referral
	entity, err = query.UpdateReferralStatus(self.db, referralID.String(), status)
	if err == nil {
		referral = entity.ToDomain()
	}
	return
}
