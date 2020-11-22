package models

type Function struct {
	Name  string
	Image string
}

type FnData struct {
	Name string
	File string // base64 encoded zip file
}
