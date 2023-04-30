package appContextConfig

import (
	"ClubMembershipPortal/lib/auth"
	"github.com/joho/godotenv"
	"os"
)

type ApplicationConfiguration struct {
	ApplicationName    string
	HttpRootPath       string
	HttpServerPort     string
	DBServerName       string
	DBDatabaseName     string
	DBUsername         string
	DBPassword         string
	SMTPServer         string
	PapertrailEndPoint string
	AuthConfig         auth.AuthenticatorConfig
}

// ----------------------------------------------------------------------------------------------
func LoadEnvFile() (ApplicationConfiguration, error) {
	err := godotenv.Load()
	if err != nil {
		return ApplicationConfiguration{}, err
	}

	ac := ApplicationConfiguration{}
	ac.ApplicationName = os.Getenv("APPLICATIONNAME")
	ac.HttpRootPath = os.Getenv("HTTPROOTPATH")
	ac.HttpServerPort = os.Getenv("HTTPSERVERPORT")
	ac.DBServerName = os.Getenv("DBSERVERNAME")
	ac.DBDatabaseName = os.Getenv("DBDATABASENAME")
	ac.DBUsername = os.Getenv("DBUSERNAME")
	ac.DBPassword = os.Getenv("DBPASSWORD")
	ac.SMTPServer = os.Getenv("SMTPSERVER")
	ac.PapertrailEndPoint = os.Getenv("PAPERTRAILENDPOINT")
	ac.AuthConfig = auth.LoadAuthConfigFromEnv()

	return ac, nil
}
