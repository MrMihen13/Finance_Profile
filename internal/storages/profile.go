package storages

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"profile/internal/models"
	"reflect"
)

const opProfileName string = "storages.Profile"

type Profile struct {
	log *slog.Logger
	db  *sql.DB
}

func NewProfileStorage(log *slog.Logger, db *sql.DB) *Profile {
	return &Profile{
		log: log,
		db:  db,
	}
}

func (p *Profile) GetByID(id uuid.UUID, profile *models.Profile) error {
	const op = opProfileName + ".GetByID"
	log := p.log.With("op", op, "id", id)
	query := "SELECT * FROM profile WHERE id = $1"

	log.Debug("getting user profile", "query", query)
	row := p.db.QueryRow(query, id)
	err := row.Scan(&profile.ID, &profile.Email, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("profile not found") // TODO add static error
		}
		log.Error("error scanning profile", "error", err.Error())
		return err
	}

	log.Debug("found profile", "email", profile.Email)

	return nil
}

func (p *Profile) IsExist(id uuid.UUID) (isExist bool) {
	const op = opProfileName + ".IsExist"
	log := p.log.With("op", op, "id", id)

	query := "SELECT EXISTS(SELECT id FROM profile WHERE id = $1)"
	log.Debug("checking user profile", "query", query)

	row := p.db.QueryRow(query, id)
	err := row.Scan(&isExist)

	if err != nil {
		log.Error("error checking user profile", "error", err)
		return false
	}

	log.Info("user profile exists", "isExist", isExist)

	return isExist
}

func (p *Profile) IsEmailExist(email string) (isExist bool) {
	const op = opProfileName + ".IsEmailExist"
	log := p.log.With("op", op, "email", email)

	query := "SELECT EXISTS(SELECT id FROM profile.public.profile WHERE email = $1)"
	log.Debug("checking user profile", "query", op)

	row := p.db.QueryRow(query, email)
	err := row.Scan(&isExist)

	if err != nil {
		log.Error("error checking user profile", "error", err)
		return false
	}

	log.Info("email exists", "isExist", isExist)

	return isExist
}

func (p *Profile) UpdateEmail(id uuid.UUID, newEmail string) error {
	const op = opProfileName + ".UpdateEmail"
	log := p.log.With("op", op, "id", id)

	query := "UPDATE profile.public.profile SET email = $1 WHERE id = $2"
	log.Debug("updating email in profile", "query", query)

	_, err := p.db.Exec(query, newEmail, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("profile not found")
		}
		log.Error("error updating email in profile", "error", err)
		return err
	}

	log.Info("email updated in profile", "email", newEmail)

	return nil
}

func (p *Profile) Insert(profile *models.Profile) error {
	const op = opProfileName + ".Insert"
	log := p.log.With("op", op, "email", profile.Email, "id", profile.ID)
	query := "INSERT INTO profile (id, email) VALUES ($1, $2)"

	if reflect.ValueOf(profile.ID).IsZero() {
		log.Debug("setting value for ID", "id", profile.ID)
		profile.ID = uuid.New()
	}

	log.Debug("inserting profile", "query", query)
	if _, err := p.db.Exec(query, profile.ID, profile.Email); err != nil {
		log.Error("error inserting profile", "error", err.Error())
		return err
	}

	log.Info("user profile inserted", "email", profile.Email)

	return nil
}

func (p *Profile) Delete(id uuid.UUID) error {
	const op = opProfileName + ".Delete"
	log := p.log.With("op", op, "id", id)
	query := "DELETE FROM profile.public.profile WHERE id = $1"

	log.Debug("deleting profile", "query", query)
	if _, err := p.db.Exec(query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("profile not found") // TODO add static err
		}
		log.Error("error deleting profile", "error", err)
		return err
	}

	log.Info("deleted profile", "id", id)

	return nil
}
