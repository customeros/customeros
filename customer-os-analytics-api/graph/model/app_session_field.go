package model

import (
	"fmt"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
)

func (f *AppSessionField) GetColumnName() string {
	switch f.String() {
	case AppSessionFieldCountry.String():
		return entity.SessionColumnName_Country
	case AppSessionFieldCity.String():
		return entity.SessionColumnName_City
	case AppSessionFieldRegion.String():
		return entity.SessionColumnName_RegionName
	case AppSessionFieldReferrerSource.String():
		return entity.SessionColumnName_ReferrerSource
	case AppSessionFieldUtmCampaign.String():
		return entity.SessionColumnName_UtmCampaign
	case AppSessionFieldUtmContent.String():
		return entity.SessionColumnName_UtmContent
	case AppSessionFieldUtmNetwork.String():
		return entity.SessionColumnName_UtmNetwork
	case AppSessionFieldUtmMedium.String():
		return entity.SessionColumnName_UtmMedium
	case AppSessionFieldUtmSource.String():
		return entity.SessionColumnName_UtmSource
	case AppSessionFieldUtmTerm.String():
		return entity.SessionColumnName_UtmTerm
	default:
		panic(fmt.Errorf("unexpected field %s", f.String()))
	}
}
