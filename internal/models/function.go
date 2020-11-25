package models

// Function represents a single function that a user deploys and can invoke
type Function struct {
	Name  string
	Image string
}

// FnData contains the name and the data of a function that a user is deploying. The data is stored
// as a base64 encoded zip file containing the function code
type FnData struct {
	Name string
	// base64 encoded zip file
	File string
}
