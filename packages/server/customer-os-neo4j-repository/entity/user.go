package entity

import "time"

type UserEntity struct {
	DataLoaderKey
	Id              string
	FirstName       string     `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName        string     `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Name            string     `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt       time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt       time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	Source          DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth   DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE;supportCaseSensitive:false"`
	AppSource       string     `neo4jDb:"property:appSource;lookupName:APP_SOURCE;supportCaseSensitive:false"`
	Roles           []string   `neo4jDb:"property:roles;lookupName:ROLES;supportCaseSensitive:false"`
	Timezone        string     `neo4jDb:"property:timezone;lookupName:TIMEZONE;supportCaseSensitive:true"`
	ProfilePhotoUrl string     `neo4jDb:"property:profilePhotoUrl;lookupName:PROFILE_PHOTO_URL;supportCaseSensitive:true"`
	Internal        bool
	Bot             bool

	// Indirect properties
	DefaultForPlayer bool
	Tenant           string
}

type UserEntities []UserEntity

func (u UserEntity) GetFullName() string {
	fullName := u.FirstName
	if u.LastName != "" {
		fullName += " " + u.LastName
	}
	if fullName == "" {
		fullName = u.Name
	}
	return fullName
}
