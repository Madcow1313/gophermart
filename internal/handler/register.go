package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

const (
	errorWrongType         = "Wrong content type"
	errorBody              = "Can't ready body"
	errorUnmarshal         = "Can't unmarshal json"
	errorLoginAlreadyInUse = "Login already in use"
	errorDB                = "Internal server error: something went wrong while connecting to database"
	parseFormError         = "Can't parse form"
	cookiePrefix           = "user_id="
)

func (hh *Handler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cType := r.Header.Get("Content-Type")
		if cType != "application/json" {
			http.Error(w, errorWrongType, http.StatusBadRequest)
			return
		}
		err := r.ParseForm()
		if err != nil {
			http.Error(w, parseFormError, http.StatusBadRequest)
			return
		}
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, errorBody, http.StatusBadRequest)
			return
		}

		var regCreds RegisterDataJSON
		err = json.Unmarshal(b, &regCreds)
		if err != nil {
			http.Error(w, errorUnmarshal, http.StatusBadRequest)
			return
		}

		err = hh.DBConnector.ConnectToRegisterDB(func(db *sql.DB, args ...interface{}) error {
			return hh.DBConnector.InsertUserCredentials(db, regCreds.Login, regCreds.Password, strings.TrimPrefix(w.Header().Get("Set-Cookie"), cookiePrefix))
		})
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			http.Error(w, errorLoginAlreadyInUse, http.StatusConflict)
			return
		} else if err != nil {
			fmt.Println(err)
			http.Error(w, errorDB, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
