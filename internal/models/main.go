package models

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

// Database models
type Session struct {
	gorm.Model
	Time    int64  `json:"time" form:"time"`
	Subject string `json:"subject" form:"subject"`
	Title   string `json:"title" form:"title"`
}

type File struct {
	gorm.Model
	Name          string `json:"name" gorm:"unique;"`
	Path          string `json:"path"`
	ParentSession int64  `json:"parentSession"`
}

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique;"`
	PasswordHash []byte `json:"-"`
	Salt         string `json:"-"`
}

// Response models
type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SessionWithFiles struct {
	Session
	Files []File `json:"files"`
}

// Misc. models
type Settings struct {
	AllowRegistration bool `json:"allowRegistration"`
}

func (b Settings) GetFields() []string {
	var resp []string
	val := reflect.ValueOf(b)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)

		switch jsonTag := t.Tag.Get("json"); jsonTag {
		case "":
			resp = append(resp, t.Name)
		case "-":
		default:
			var fieldName string
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}
			resp = append(resp, fieldName)
		}
	}

	return resp
}
