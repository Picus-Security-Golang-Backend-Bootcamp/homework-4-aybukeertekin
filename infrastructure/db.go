package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewPostgreDB(conString string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "go_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})

	if err != nil {
		panic(fmt.Sprintln("Cannot connect to db!"))
	}

	return db
}
