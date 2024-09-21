package services

import (
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"profile/internal/models"
)

const opProfileName string = "services.Profile"

type storage interface {
	Insert(profile *models.Profile) error
	GetByID(id uuid.UUID, profile *models.Profile) error
	IsEmailExist(email string) (isExist bool)
	IsExist(id uuid.UUID) (isExist bool)
	UpdateEmail(id uuid.UUID, newEmail string) error
	Delete(id uuid.UUID) error
}

type Profile struct {
	log     *slog.Logger
	storage storage
}

func NewProfileService(log *slog.Logger, storage storage) *Profile {
	return &Profile{
		log:     log,
		storage: storage,
	}
}

func (p *Profile) Create(email string) (*models.Profile, error) {
	const op = opProfileName + ".Create"
	log := p.log.With("op", op, "email", email)

	log.Info("creating new profile")

	log.Debug("check email is exist")
	if p.storage.IsEmailExist(email) {
		log.Info("email already exists")
		return nil, errors.New("email already exists")
	}

	log.Debug("create profile in db")
	profile := &models.Profile{
		ID:    uuid.New(),
		Email: email,
	}
	if err := p.storage.Insert(profile); err != nil {
		log.Error("failed to create profile")
		return nil, err
	}

	log.Debug("filling profile data")
	if err := p.storage.GetByID(profile.ID, profile); err != nil {
		log.Error("failed to get profile")
		return nil, err
	}

	log.Info("successfully created profile")
	return profile, nil
}

func (p *Profile) GetByID(id uuid.UUID) (*models.Profile, error) {
	const op = opProfileName + ".GetByID"
	log := p.log.With("op", op, "id", id)

	log.Info("getting profile by id")

	log.Debug("getting profile")
	profile := &models.Profile{}
	if err := p.storage.GetByID(id, profile); err != nil {
		log.Error("failed to get profile")
		return nil, err
	}

	log.Info("successfully fetched profile")
	return profile, nil
}

func (p *Profile) GetByEmail(email string) (*models.Profile, error) {
	const op = opProfileName + ".GetByEmail"
	return nil, nil
}

func (p *Profile) UpdateEmail(id uuid.UUID, newEmail string) (*models.Profile, error) {
	const op = opProfileName + ".UpdateEmail"
	log := p.log.With("op", op, "id", id, "newEmail", newEmail)

	log.Info("updating profile email")

	log.Debug("checking is new email exist")
	if p.storage.IsEmailExist(newEmail) {
		log.Warn("email already exists")
		return nil, errors.New("email already exists")
	}

	log.Debug("checking is profile exist")
	if !p.storage.IsExist(id) {
		log.Warn("profile not found")
		return nil, errors.New("profile not found")
	}

	log.Debug("updating email in profile")
	if err := p.storage.UpdateEmail(id, newEmail); err != nil {
		log.Error("failed to update email")
		return nil, err
	}

	log.Debug("filling new profile data")
	profile := &models.Profile{}
	if err := p.storage.GetByID(id, profile); err != nil {
		log.Error("failed to get profile")
		return nil, err
	}

	log.Info("successfully updated profile")
	return profile, nil
}

func (p *Profile) DeleteByID(id uuid.UUID) error {
	const op = opProfileName + ".DeleteByID"
	log := p.log.With("op", op, "id", id)

	log.Debug("checking is profile exist")
	if !p.storage.IsExist(id) {
		log.Info("profile not found")
		return errors.New("profile not found")
	}

	log.Debug("deleting profile")
	if err := p.storage.Delete(id); err != nil {
		log.Error("failed to delete profile")
		return err
	}

	log.Info("successfully deleted profile")
	return nil
}
