package appContextConfig

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func openMSSQLDBwithPing(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func OpenDB(DBServerName string, DBDatabaseName string, DBUsername string, DBPassword string, ApplicationName string) (*sql.DB, error) {

	connectionString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;connection timeout=300;app name=%s",
		DBServerName, DBDatabaseName, DBUsername, DBPassword, ApplicationName)

	db, err := openMSSQLDBwithPing(connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
