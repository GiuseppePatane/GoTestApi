package app

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/dstroot/postgres-api/middleware/connlimit"
	env "github.com/joeshaw/envdecode"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

type App struct {
	Router *httprouter.Router
	DB     *sql.DB
	Server *http.Server
	Stats  *stats.Stats
	Cfg    config
}
type config struct {
	HostName string
	Debug    bool   `env:"DEBUG,default=true"`
	Port     string `env:"PORT,default=8000"`

	SQL struct {
		Host     string `env:"SQL_HOST,default=localhost"`
		Port     string `env:"SQL_PORT,default=5432"`
		User     string `env:"SQL_USER,default=postgres"`
		Password string `env:"SQL_PASSWORD,default=root"`
		Database string `env:"SQL_DATABASE,default=Identity"`
	}
}

func Initialize() (app App, err error) {

	err = env.Decode(&app.Cfg)
	if err != nil {
		return app, errors.Wrap(err, "configuration decode failed")
	}
	app.Cfg.HostName, _ = os.Hostname()
	connString := "postgres://" + app.Cfg.SQL.User +
		":" + app.Cfg.SQL.Password +
		"@" + app.Cfg.SQL.Host +
		":" + app.Cfg.SQL.Port +
		"/" + app.Cfg.SQL.Database +
		"?sslmode=disable"

	app.DB, err = sql.Open("postgres", connString)
	if err != nil {
		return app, errors.Wrap(err, "database connection failed")
	}
	err = app.DB.Ping()
	if err != nil {
		return app, errors.Wrap(err, "error pinging database")
	}
	app.Router = httprouter.New()
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	// setup stats https://github.com/thoas/stats
	app.Stats = stats.New()
	n.Use(app.Stats)

	// Connections limiter
	// Manage connections before rate?
	n.Use(connlimit.MaxAllowed(50))

	// Rate limiter
	//limiter := tollbooth.NewLimiter(50, time.Second)
	//n.Use(tollbooth_negroni.LimitHandler(limiter))

	n.UseHandler(app.Router)

	/**
	 * Server
	 */

	app.Server = &http.Server{
		Addr:           ":" + app.Cfg.Port,
		Handler:        n, // pass in router
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return app, nil
}
