package database

import (
	"backend/database/entity"
	"backend/loggers"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateFinalTestinationDB(user, password, host, port, dbName string) *FinalTestinationDB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbName, port)
	config := gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:                  logger.Info,
				Colorful:                  true,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      false,
				SlowThreshold:             time.Second,
			},
		),
	}
	db, err := gorm.Open(postgres.Open(dsn), &config)
	if err != nil {
		loggers.Error.Fatalf("Error connecting to the database: %s", err)
	}

	return &FinalTestinationDB{db}
}

type FinalTestinationDB struct {
	Orm *gorm.DB
}

func (db *FinalTestinationDB) CreateSchemas() {
	if err := db.Orm.AutoMigrate(
		&entity.Block{},
		&entity.Game{},
		&entity.Player{},
		&entity.PlayerGame{},
		&entity.Icon{},
	); err != nil {
		loggers.Error.Fatalf("Error creating schemas: %s", err)
	}
}
