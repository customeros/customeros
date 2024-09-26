package entity

import "time"

type ScrapInFlow string

const (
	ScrapInFlowPersonSearch   ScrapInFlow = "PERSON_SEARCH"
	ScrapInFlowPersonProfile  ScrapInFlow = "PERSON_PROFILE"
	ScrapInFlowCompanySearch  ScrapInFlow = "COMPANY_SEARCH"
	ScrapInFlowCompanyProfile ScrapInFlow = "COMPANY_PROFILE"
)

type EnrichDetailsScrapIn struct {
	ID            uint64      `gorm:"primary_key;autoIncrement:true" json:"id"`
	Flow          ScrapInFlow `gorm:"column:flow;type:varchar(255);NOT NULL" json:"flow"`
	Param1        string      `gorm:"column:param1;type:varchar(1000);" json:"param1"`
	Param2        string      `gorm:"column:param2;type:varchar(1000);" json:"param2"`
	Param3        string      `gorm:"column:param3;type:varchar(1000);" json:"param3"`
	Param4        string      `gorm:"column:param4;type:varchar(1000);" json:"param4"`
	AllParamsJson string      `gorm:"column:all_params_json;type:text;DEFAULT:'';NOT NULL" json:"allParams"`
	Data          string      `gorm:"column:data;type:text;DEFAULT:'';NOT NULL" json:"data"`
	CreatedAt     time.Time   `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time   `gorm:"column:updated_at;type:timestamp;;DEFAULT:current_timestamp" json:"updatedAt"`
	Success       bool        `gorm:"column:success;type:boolean;DEFAULT:false" json:"success"`
	PersonFound   bool        `gorm:"column:person_found;type:boolean;DEFAULT:false" json:"personFound"`
	CompanyFound  bool        `gorm:"column:company_found;type:boolean;DEFAULT:false" json:"companyFound"`
}

func (EnrichDetailsScrapIn) TableName() string {
	return "enrich_details_scrapin"
}

// ScrapInResponseBody is getting serialized in the Data field
type ScrapInResponseBody struct {
	Success       bool                   `json:"success"`
	Email         string                 `json:"email"`
	EmailType     string                 `json:"emailType"`
	CreditsLeft   int                    `json:"credits_left"`
	RateLimitLeft int                    `json:"rate_limit_left"`
	Person        *ScrapinPersonDetails  `json:"person,omitempty"`
	Company       *ScrapinCompanyDetails `json:"company,omitempty"`
}

type ScrapinPersonDetails struct {
	PublicIdentifier   string `json:"publicIdentifier"`
	LinkedInIdentifier string `json:"linkedInIdentifier"`
	LinkedInUrl        string `json:"linkedInUrl"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Headline           string `json:"headline"`
	Location           string `json:"location"`
	Summary            string `json:"summary"`
	PhotoUrl           string `json:"photoUrl"`
	CreationDate       struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	} `json:"creationDate"`
	FollowerCount int `json:"followerCount"`
	Positions     struct {
		PositionsCount  int `json:"positionsCount"`
		PositionHistory []struct {
			Title        string `json:"title"`
			CompanyName  string `json:"companyName"`
			Description  string `json:"description"`
			StartEndDate struct {
				Start *struct {
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"start"`
				End *struct {
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"end"`
			} `json:"startEndDate"`
			CompanyLogo string `json:"companyLogo"`
			LinkedInUrl string `json:"linkedInUrl"`
			LinkedInId  string `json:"linkedInId"`
		} `json:"positionHistory"`
	} `json:"positions"`
	Schools struct {
		EducationsCount  int `json:"educationsCount"`
		EducationHistory []struct {
			DegreeName   string      `json:"degreeName"`
			FieldOfStudy string      `json:"fieldOfStudy"`
			Description  interface{} `json:"description"` // Can be null, so use interface{}
			LinkedInUrl  string      `json:"linkedInUrl"`
			SchoolLogo   string      `json:"schoolLogo"`
			SchoolName   string      `json:"schoolName"`
			StartEndDate struct {
				Start struct {
					Month *int `json:"month"` // Can be null, so use pointer
					Year  *int `json:"year"`  // Can be null, so use pointer
				} `json:"start"`
				End struct {
					Month *int `json:"month"` // Can be null, so use pointer
					Year  *int `json:"year"`  // Can be null, so use pointer
				} `json:"end"`
			} `json:"startEndDate"`
		} `json:"educationHistory"`
	} `json:"schools"`
	Skills    []interface{} `json:"skills"`    // Can be empty, so use interface{}
	Languages []interface{} `json:"languages"` // Can be empty, so use interface{}
}

type ScrapinCompanyDetails struct {
	LinkedInId         string `json:"linkedInId"`
	Name               string `json:"name"`
	UniversalName      string `json:"universalName"`
	LinkedInUrl        string `json:"linkedInUrl"`
	EmployeeCount      int    `json:"employeeCount"`
	EmployeeCountRange struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"employeeCountRange"`
	WebsiteUrl    string      `json:"websiteUrl"`
	Tagline       interface{} `json:"tagline"` // Can be null, so use interface{}
	Description   string      `json:"description"`
	Industry      string      `json:"industry"`
	Phone         interface{} `json:"phone"` // Can be null, so use interface{}
	Specialities  []string    `json:"specialities"`
	FollowerCount int         `json:"followerCount"`
	Headquarter   struct {
		City           string      `json:"city"`
		Country        string      `json:"country"`
		PostalCode     string      `json:"postalCode"`
		GeographicArea string      `json:"geographicArea"`
		Street1        string      `json:"street1"`
		Street2        interface{} `json:"street2"` // Can be null, so use interface{}
	} `json:"headquarter"`
	Logo      string `json:"logo"`
	FoundedOn struct {
		Year int `json:"year"`
	}
}

func (c ScrapinCompanyDetails) HeadquarterIsEmpty() bool {
	return c.Headquarter.City == "" && c.Headquarter.Country == "" && c.Headquarter.PostalCode == "" && c.Headquarter.GeographicArea == "" && c.Headquarter.Street1 == ""
}

func (c ScrapinCompanyDetails) GetEmployeeCount() int64 {
	if c.EmployeeCount > 0 {
		return int64(c.EmployeeCount)
	} else if c.EmployeeCountRange.Start > 0 {
		return int64(c.EmployeeCountRange.Start)
	} else if c.EmployeeCountRange.End > 0 {
		return int64(c.EmployeeCountRange.End)
	}
	return 0
}
