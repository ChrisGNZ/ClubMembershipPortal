package appContextConfig

import (
	"ClubMembershipPortal/lib/auth"
	"database/sql"
	syslog "github.com/RackSec/srslog"
)

type Application struct {
	SysLog *syslog.Writer
	DBconn *sql.DB
	Auth   *auth.Authenticator
	Config ApplicationConfiguration
}

// -------------------------------------------------------------------------------------------------------------------
func CreateApplicationContext(cfg ApplicationConfiguration, auth *auth.Authenticator) (*Application, string, error) {

	appContext := Application{}

	//initialise syslog connection to papertrail
	sl, err := syslog.Dial("udp", cfg.PapertrailEndPoint, syslog.LOG_ERR, cfg.ApplicationName)
	if err != nil {
		//panic("failed to dial syslog")
		return nil, "Error opening connection to PaperTrail", err
	}

	db, err := OpenDB(cfg.DBServerName, cfg.DBDatabaseName, cfg.DBUsername, cfg.DBPassword, cfg.ApplicationName)
	if err != nil {
		sl.Err(LogEntry(err.Error(), 1))
		sl.Err("Error opening database connection")
		return nil, "Error opening database connection", err
	}

	appContext.SysLog = sl
	appContext.DBconn = db
	appContext.Config = cfg
	appContext.Auth = auth

	sl.Info(LogEntry("Starting Application: "+cfg.ApplicationName+", with http root: "+cfg.HttpRootPath, 1))
	return &appContext, "OK", nil
}
