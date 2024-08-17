package dbconnector

import "database/sql"

const InsertUserCredentialsQuery = "INSERT INTO gophermart_users (login, password, cookie) VALUES ($1, $2, $3)"

func (c *Connector) InsertUserCredentials(db *sql.DB, login string, password string, cookie string) error {
	stmt, err := db.Prepare(InsertUserCredentialsQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(login, password, cookie)
	if err != nil {
		return err
	}
	err = stmt.Close()
	return err
}
