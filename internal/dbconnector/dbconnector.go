package dbconnector

import "database/sql"

const (
	CreateUsersQuery = "CREATE TABLE IF NOT EXISTS gophermart_users (login varchar unique not null, password varchar not null, cookie varchar, balance float default 0)"
	CreateGophermart = "CREATE TABLE IF NOT EXISTS gophermart (login varchar)"
	sqlDriver        = "postgres"
)

type Connector struct {
	DatabaseDSN  string
	LastResult   string
	IsDeleted    bool
	URLmap       map[string]string
	UserURLS     map[string][]string
	RegisterDB   *sql.DB
	GophermartDB *sql.DB
}

func NewConnector(databaseDSN string) *Connector {
	return &Connector{DatabaseDSN: databaseDSN}
}

func (c *Connector) ConnectToRegisterDB(connectFunc func(db *sql.DB, args ...interface{}) error) error {
	if c.RegisterDB == nil {
		db, err := sql.Open(sqlDriver, c.DatabaseDSN)
		if err != nil {
			return err
		}
		c.RegisterDB = db
	}
	if connectFunc != nil {
		err := connectFunc(c.RegisterDB)
		return err
	}
	return nil
}

func (c *Connector) ConnectToGophermartDB(connectFunc func(db *sql.DB, args ...interface{}) error) error {
	if c.GophermartDB == nil {
		db, err := sql.Open(sqlDriver, c.DatabaseDSN)
		if err != nil {
			return err
		}
		c.GophermartDB = db
	}
	if connectFunc != nil {
		err := connectFunc(c.GophermartDB)
		return err
	}
	return nil
}

func (c *Connector) CreateTable(db *sql.DB, creationQuery string) error {
	_, err := db.Exec(creationQuery)
	if err != nil {
		return err
	}
	return nil
}
