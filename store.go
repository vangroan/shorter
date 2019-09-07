package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Storage is a repository that can persistent
// URL mappings, called location records.
type Storage interface {
	SaveLocation(location *Location) error
	GetLocation(id uint64) (Location, error)
	DeleteSince(now time.Time) error
}

// DBStorage persists location records in
// a database.
type DBStorage struct {
	db *gorm.DB
}

// NewDBStorage creates a location storage
// backed by a database.
func NewDBStorage(db *gorm.DB) DBStorage {
	return DBStorage{db: db}
}

// SaveLocation creates a new location record
// in the database.
//
// On successful create, the primary key of
// the record will be written to the referenced
// location argument.
func (s *DBStorage) SaveLocation(location *Location) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&location).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetLocation retrieves an existing
// location record.
func (s *DBStorage) GetLocation(id uint64) (Location, error) {
	var location Location

	if err := s.db.First(&location, id).Error; err != nil {
		return Location{}, err
	}

	return location, nil
}

// DeleteSince hard deletes URLs the given time,
// minus each URLs time-to-live.
func (s *DBStorage) DeleteSince(now time.Time) error {
	err := s.db.
		Where("ttl <> 0").
		Where("created_at < datetime(?, -ttl || ' seconds')", now).
		Delete(Location{}).
		Error

	return err
}
