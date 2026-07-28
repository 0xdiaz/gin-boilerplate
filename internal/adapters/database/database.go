package database

import (
	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

// DbConnection create database connection
func DbConnection(masterDSN, replicaDSN string) error {
	logMode := false
	debug := false
	if c := config.Get(); c != nil {
		logMode = c.Database.LogMode
		debug = c.Server.Debug
	}

	loglevel := gormlogger.Silent
	if logMode {
		loglevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(loglevel),
	})
	if err != nil {
		logger.Errorf("database connection error: %v", err)
		return err
	}

	if !debug {
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				postgres.Open(replicaDSN),
			},
			Policy: dbresolver.RandomPolicy{},
		}))
	}

	DB = db
	return nil
}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}
