package database

import (
	"sync"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/util"
	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
	once   sync.Once
)

func CreateConnection() {
	conf := dbConfig{
		User: util.Getenv("DB_USER", "user"),
		Pass: util.Getenv("DB_PASS", "#!UserAcc25"),
		Host: util.Getenv("DB_HOST", "localhost"),
		Port: util.Getenv("DB_PORT", "3306"),
		Name: util.Getenv("DB_NAME", "reservations_db"),
	}

	mysql := mysqlConfig{dbConfig: conf}
	once.Do(func() {
		mysql.Connect()
	})
}

func GetConnection() *gorm.DB {
	if dbConn == nil {
		CreateConnection()
	}
	return dbConn
}
