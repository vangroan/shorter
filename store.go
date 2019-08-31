package main

import (
	"github.com/jinzhu/gorm"
)

type Storage interface {
	SaveLocation(location *Location) error
	GetLocation(id uint64) (Location, error)
}

type SQLiteStorage struct {
	db *gorm.DB
}

func NewSQLitestorage(db *gorm.DB) SQLiteStorage {
	return SQLiteStorage{db: db}
}

func (s *SQLiteStorage) SaveLocation(location *Location) error {
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

func (s *SQLiteStorage) GetLocation(id uint64) (Location, error) {
	var location Location

	if err := s.db.First(&location, id).Error; err != nil {
		return Location{}, err
	}

	return location, nil
}
