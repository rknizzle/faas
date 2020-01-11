package api

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fnData struct {
	File string `json:"file"`
	Name string `json:"name"`
}

func main() {
	Init()
}

func Init() {
	fmt.Println("Starting...")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// add a new function
	r.POST("/functions", func(c *gin.Context) {

		// get function data from request body
		rawData, err := c.GetRawData()
		if err != nil {
			panic(err)
		}

		var data fnData
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			panic(err)
		}

		// the input "file" is expected to be base64 encoded
		// decode the zip file
		decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data.File))
		if err != nil {
			panic(err)
		}
		// name to be given to the zip file thats written to disk
		zipFile := data.Name + ".zip"

		output, err := os.Create(zipFile)
		if err != nil {
			panic(err)
		}
		// Close output file
		defer output.Close()

		// write decoded data to zip file
		io.Copy(output, decoder)

		// unzip the input zip file to extract the directory containing the function code
		files, err := Unzip(zipFile, data.Name)
		if err != nil {
			panic(err)
		}
		// print files in unzipped directory
		fmt.Println(files)
	})

	// Listen and serve on localhost
	r.Run()
}

// Function found at https://golangcode.com/unzip-files-in-go/ (MIT License)
// Unzip a file
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
