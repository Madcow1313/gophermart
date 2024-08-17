package dbconnector

import "database/sql"

const (
	checkUserQuery = "SELECT login FROM gophermart_users WHERE (login = $1 AND password = $2) OR cookie = $3"
)

func (c *Connector) CheckUserCredentials(db *sql.DB, login string, password string, cookie string) error {
	row := db.QueryRow(checkUserQuery, login, password, cookie)
	var temp string
	err := row.Scan(&temp)
	return err
}
