package models

import (
	"fmt"
	"time"
)

type User struct {
	Model    model.Model
	Name     string
	Password string
	Encr
	CreatedAt time.Time
}
