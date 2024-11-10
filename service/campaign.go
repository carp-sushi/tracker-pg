package service

import "time"

// CampaignService manages campaigns
type CampaignService interface {
	CampaignReader
	CampaignWriter
}

// CampaignReader reads campaigns
type CampaignReader interface {
	GetCampaign(CampaignID) (Campaign, error)
	GetCampaigns(Account, PageParams) CampaignPage
}

// CampaignWriter writes campaigns
type CampaignWriter interface {
	CreateCampaign(Account, string) (Campaign, error)
	UpdateCampaign(id CampaignID, name string, expiresAt time.Time) (Campaign, error)
}
