package config

import(
"gorm.io/gorm"
"gorm.io/driver/mysql"
"log"
"MyApp/global"
"time"
)


func InitDB(){

	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to initialize database , got error %v", err)
	}


sqlDB, err := db.DB()

sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
sqlDB.SetConnMaxLifetime(time.Hour)	

if err != nil {
	log.Fatalf("failed to configure database , got error %v", err)
}

global.DB = db


}