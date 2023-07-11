package hubspot

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"strconv"
	"time"
)

func MapOrganization(jsonStr string) (entity.OrganizationData, error) {
	var data entity.OrganizationData
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return entity.OrganizationData{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Extract field values from JSON and assign them to OrganizationData struct fields
	if id, ok := jsonData["id"].(string); ok {
		data.ExternalId = id
	}
	if createdAt, ok := jsonData["createdAt"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return entity.OrganizationData{}, fmt.Errorf("failed to parse createdAt field: %v", err)
		}
		data.CreatedAt = parsedTime
	}
	if updatedAt, ok := jsonData["updatedAt"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339Nano, updatedAt)
		if err != nil {
			return entity.OrganizationData{}, fmt.Errorf("failed to parse updatedAt field: %v", err)
		}
		data.UpdatedAt = parsedTime
	}

	// Extract the "properties" field from JSON
	jsonDataProperties, ok := jsonData["properties"].(map[string]interface{})
	if !ok {
		return entity.OrganizationData{}, fmt.Errorf("failed to parse 'properties' field from JSON")
	}
	if name, ok := jsonDataProperties["name"].(string); ok {
		data.Name = name
	}
	if description, ok := jsonDataProperties["description"].(string); ok {
		data.Description = description
	}
	if website, ok := jsonDataProperties["website"].(string); ok {
		data.Website = website
	}
	if industry, ok := jsonDataProperties["industry"].(string); ok {
		data.Industry = industry
	}
	if isPublic, ok := jsonDataProperties["is_public"].(bool); ok {
		data.IsPublic = isPublic
	}
	if employees, ok := jsonDataProperties["numberofemployees"].(float64); ok {
		data.Employees = int64(employees)
	}
	if phoneNumber, ok := jsonDataProperties["phone"].(string); ok {
		data.PhoneNumber = phoneNumber
	}
	if country, ok := jsonDataProperties["country"].(string); ok {
		data.Country = country
	}
	if state, ok := jsonDataProperties["state"].(string); ok {
		data.Region = state
	}
	if city, ok := jsonDataProperties["city"].(string); ok {
		data.Locality = city
	}
	if address, ok := jsonDataProperties["address"].(string); ok {
		data.Address = address
	}
	if address2, ok := jsonDataProperties["address2"].(string); ok {
		data.Address2 = address2
	}
	if zip, ok := jsonDataProperties["zip"].(string); ok {
		data.Zip = zip
	}
	if ownerId, ok := jsonDataProperties["hubspot_owner_id"].(string); ok {
		data.UserExternalOwnerId = ownerId
	}
	if domain, ok := jsonDataProperties["domain"].(string); ok {
		data.Domains = []string{domain}
	}
	if companyType, ok := jsonDataProperties["type"].(string); ok {
		switch companyType {
		case "PROSPECT":
			data.RelationshipName = entity.Customer
			data.RelationshipStage = entity.Prospect
		case "PARTNER":
			data.RelationshipName = entity.Partner
			data.RelationshipStage = entity.Live
		case "RESELLER":
			data.RelationshipName = entity.Reseller
			data.RelationshipStage = entity.Live
		case "VENDOR":
			data.RelationshipName = entity.Vendor
			data.RelationshipStage = entity.Live
		}
	}

	return data, nil
}

func MapUser(jsonStr string) (entity.UserData, error) {
	var data entity.UserData
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return entity.UserData{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Extract field values from JSON and assign them to OrganizationData struct fields
	if userId, ok := jsonData["userId"].(float64); ok {
		data.ExternalId = strconv.FormatInt(int64(userId), 10)
	}
	if createdAt, ok := jsonData["createdAt"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return entity.UserData{}, fmt.Errorf("failed to parse createdAt field: %v", err)
		}
		data.CreatedAt = parsedTime
	}
	if updatedAt, ok := jsonData["updatedAt"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339Nano, updatedAt)
		if err != nil {
			return entity.UserData{}, fmt.Errorf("failed to parse updatedAt field: %v", err)
		}
		data.UpdatedAt = parsedTime
	}
	if id, ok := jsonData["id"].(string); ok {
		data.ExternalOwnerId = id
	}
	if lastName, ok := jsonData["lastName"].(string); ok {
		data.LastName = lastName
	}
	if firstName, ok := jsonData["firstName"].(string); ok {
		data.FirstName = firstName
	}
	if email, ok := jsonData["email"].(string); ok {
		data.Email = email
	}

	return data, nil
}
