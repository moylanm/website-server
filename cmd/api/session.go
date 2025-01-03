package main

import (
	"database/sql"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

func newSessionManager(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.IdleTimeout = 1 * time.Hour
	sessionManager.Lifetime = 12 * time.Hour

	return sessionManager
}
