package tests

import (
	. "github.com/deevatech/manager/db"
	"time"
)

type Test struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Source      string    `json:"source"`
	Spec        string    `json:"spec"`
}

func FindById(id uint64) *Test {
	var test Test
	DB.First(&test, id)
	return &test
}
