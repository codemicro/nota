package models


type Session struct {
	Time int32 "json:`time`"
	Subject string "json:`subject`"
	Files []File "json:`files`"
	Title string "json:`title`"
}


type File struct {
	Name string "json:`name`"
	Path string "json:`path`"
}