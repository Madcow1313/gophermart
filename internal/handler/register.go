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
	errorWrongLoginPair    = "Login/password mismatch"
	parseFormError         = "Can't parse form"
	cookiePrefix           = "user_id="
)

func (hh *Handler) CheckTypeAndBody(w http.ResponseWriter, r *http.Request) (RegisterDataJSON, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, errorWrongType, http.StatusBadRequest)
		return RegisterDataJSON{}, errors.New(strings.ToLower(errorWrongType))
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, parseFormError, http.StatusBadRequest)
		return RegisterDataJSON{}, err
	}

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, errorBody, http.StatusBadRequest)
		return RegisterDataJSON{}, err
	}

	var regCreds RegisterDataJSON
	err = json.Unmarshal(b, &regCreds)
	if err != nil {
		http.Error(w, errorUnmarshal, http.StatusBadRequest)
		return RegisterDataJSON{}, err
	}
	return regCreds, nil
}

func (hh *Handler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regCreds, err := hh.CheckTypeAndBody(w, r)
		if err != nil {
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
			http.Error(w, errorDB, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (hh *Handler) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regCreds, err := hh.CheckTypeAndBody(w, r)
		if err != nil {
			return
		}
		cookie, err := r.Cookie(strings.TrimSuffix(cookiePrefix, "="))
		if errors.Is(err, http.ErrNoCookie) {
			cookie = &http.Cookie{
				Value: "",
			}
		}
		err = hh.DBConnector.ConnectToRegisterDB(func(db *sql.DB, args ...interface{}) error {
			return hh.DBConnector.CheckUserCredentials(db, regCreds.Login, regCreds.Password, cookie.Value)
		})
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, errorWrongLoginPair, http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, errorDB, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
