package config

import (
	"fmt"
	"github.com/savioruz/mikti-task/week-3/internal/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os/exec"
	"time"
)

// NewDatabase creates a new database connection
func NewDatabase(log *logrus.Logger) *gorm.DB {
	dbName := "todo.db"
	dbPath := fmt.Sprintf("./db/%s", dbName)

	exec.Command("rm", dbName)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get connection: %v", err)
	}

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(100)
	connection.SetConnMaxLifetime(time.Second * time.Duration(300))

	if err := db.AutoMigrate(&entities.Todo{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
