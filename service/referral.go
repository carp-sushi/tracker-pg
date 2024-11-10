package service

// ReferralService manages campaign referrals
type ReferralService interface {
	ReferralReader
	ReferralWriter
}

// ReferralReader reads campaign referrals
type ReferralReader interface {
	GetReferrals(CampaignID, PageParams) ReferralPage
	GetReferralsWithStatus(ReferralStatus, PageParams) ReferralPage
}

// ReferralWriter writes campaign referrals
type ReferralWriter interface {
	CreateReferral(CampaignID, Account) (Referral, error)
	UpdateReferral(ReferralID, ReferralStatus) (Referral, error)
}
