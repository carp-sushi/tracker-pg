package service

import "github.com/carp-sushi/tracker-pg/domain"

type Account = domain.Account
type Campaign = domain.Campaign
type CampaignID = domain.CampaignID
type CampaignPage = domain.Page[domain.Campaign]
type PageParams = domain.PageParams
type Referral = domain.Referral
type ReferralID = domain.ReferralID
type ReferralPage = domain.Page[domain.Referral]
type ReferralStatus = domain.ReferralStatus
