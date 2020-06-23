package models

import "github.com/jinzhu/gorm"

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

// Response models
type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SessionWithFiles struct {
	Session
	Files []File `json:"files"`
}
