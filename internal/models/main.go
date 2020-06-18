package models

import "github.com/jinzhu/gorm"

type Session struct {
	gorm.Model
	Time int32 "json:`time`"
	Subject string "json:`subject`"
	Files []File "json:`files`"
	Title string "json:`title`"
}


type File struct {
	gorm.Model
	Name string "json:`name`"
	Path string "json:`path`"
}


type GenericResponse struct {
	Status string "json:`status`"
	Message string "json:`message`"
}