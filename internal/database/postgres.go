package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
)

const opPostgresName string = "database.Postgres"

type Postgres struct {
	log      *slog.Logger
	host     string // host of database
	port     int    // port of database
	username string // username for connect ot database
	password string // password for connect to database
	dbname   string // name of database
	db       *sql.DB
}

func NewPostgres(log *slog.Logger, host string, port int, username string, password string, dbname string) *Postgres {
	return &Postgres{
		log:      log,
		host:     host,
		port:     port,
		username: username,
		password: password,
		dbname:   dbname,
	}
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.host,
		p.port,
		p.username,
		p.password,
		p.dbname,
	)
}

func (p *Postgres) Connect() (*sql.DB, error) {
	const _op = opPostgresName + ".Connect"
	var err error
	log := p.log.With("op", _op)

	log.Debug("Connecting to database", "name", p.dbname)
	p.db, err = sql.Open("postgres", p.DSN())

	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return p.db, nil
}

func (p *Postgres) GetDB() *sql.DB {
	return p.db
}

func (p *Postgres) MustConnect() *sql.DB {
	_, err := p.Connect()
	if err != nil {
		panic(err)
	}
	return p.db
}

func (p *Postgres) Ping() error {
	const op = opPostgresName + ".Ping"
	log := p.log.With("op", op)
	log.Debug("Pinging database")
	if err := p.db.Ping(); err != nil {
		log.Error("failed to ping database", "error", err)
		return err
	}
	return nil
}

func (p *Postgres) Close() {
	const op = opPostgresName + ".Close"
	log := p.log.With("op", op)

	log.Debug("Closing database")

	if err := p.db.Close(); err != nil {
		log.Error("failed to close database", "error", err)
	}
	log.Debug("Closed database")
}
