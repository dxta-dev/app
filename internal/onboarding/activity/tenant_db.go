package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dxta-dev/app/internal/onboarding"
)

type DatabaseData struct {
	DBID     string `json:"DbId"`
	Hostname string `json:"Hostname"`
	Name     string `json:"Name"`
	DBURL    string
}
type CreateTenantDBRes struct {
	Database DatabaseData `json:"database"`
}

type CreateTenantActivities struct {
	config onboarding.CreateTenantConfig
}

type SeedOpts struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
type TenantDBRequest struct {
	Name  string   `json:"name"`
	Group string   `json:"group"`
	Seed  SeedOpts `json:"seed"`
}

func NewCreateTenantActivities(config onboarding.CreateTenantConfig) *CreateTenantActivities {
	return &CreateTenantActivities{config}
}

func (cta CreateTenantActivities) CreateTenantDB(
	ctx context.Context,
	dbDomainName string,
) (*CreateTenantDBRes, error) {
	reqBody := TenantDBRequest{
		Name:  dbDomainName,
		Group: cta.config.TursoDBGroupName,
		Seed: SeedOpts{
			Type: "database",
			Name: cta.config.TenantSeedDBURL,
		},
	}

	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return nil, errors.New("failed while marshalling request body: " + err.Error())
	}

	apiUrl := fmt.Sprintf("%s/organizations/%s/databases", cta.config.TursoApiURL, cta.config.TursoOrganizationSlug)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, errors.New("failed to create HTTP request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cta.config.TursoAuthToken))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return nil, errors.New("failed to send HTTP request: " + err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to create new tenant db with status code: " + fmt.Sprint(response.StatusCode))
	}

	responseBodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	var body CreateTenantDBRes

	err = json.Unmarshal(responseBodyBytes, &body)

	if err != nil {
		return nil, errors.New("failed to unmarshal response body: " + err.Error())
	}

	fmt.Printf("Success! Data: %v", body)

	body.Database.DBURL = fmt.Sprintf("libsql://%s", body.Database.Hostname)

	return &body, nil
}

func (cta CreateTenantActivities) AddTenantDBToMap(
	ctx context.Context,
	authId string,
	DBName string,
	DBURL string,
	DBDomainName string,
) (bool, error) {
	db, err := onboarding.GetDB(ctx, cta.config.OrganizationsTenantMapDBURL)

	if err != nil {
		return false, errors.New("failed to get organizations-tenant-map db: " + err.Error())
	}

	_, err = db.QueryContext(ctx, `
		INSERT INTO tenants 
			(organization_id, db_url, name, domain) 
		VALUES (?, ?, ?, ?);`, authId, DBURL, DBName, DBDomainName)

	if err != nil {
		return false, errors.New("failed to store tenant db data to organizations-tenant-map db: " + err.Error())
	}

	return true, nil
}
