package database

import (
	"github.com/codemicro/nota/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var (
	Conn *gorm.DB
)

func InitDatabase() {
	var err error // so it writes to Conn on the global level
	Conn, err = gorm.Open("sqlite3", "nota.db")
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connected to database")

	Conn.AutoMigrate(&models.Session{})
	Conn.AutoMigrate(&models.File{})
	log.Println("Migrated database")
}
