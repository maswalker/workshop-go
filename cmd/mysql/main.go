package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseModel
	Name  string `gorm:"column:name;unique;not null;size:256;"`
	Email string `gorm:"column:email;not null;size:256;"`
}

func (User) TableName() string {
	return "user"
}

func main() {
	dsn := "root:T3st1234!@tcp(127.0.0.1:3306)/workshop?charset=utf8mb4&parseTime=True&loc=Local"
	// 1. open db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("open db failed")
	}

	// 2. set parameters
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&User{})

	// 3. create
	user := User{Name: "alice", Email: "alice@abc.com"}
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	err = transactionSample(db)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func transactionSample(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&User{Name: "Bob"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&User{Name: "Cat"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
