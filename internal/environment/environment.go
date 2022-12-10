package environment

import "database/sql"

type DeploymentType string

const (
	Development DeploymentType = "development"
	Production  DeploymentType = "production"
)

type Env struct {
	Deployment  DeploymentType
	Db          *sql.DB
	Language    string
	LocalDomain string
	WebDomain   string
}
