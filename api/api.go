package api

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rknizzle/faas/manager"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FnData struct {
	File string `json:"file"`
	Name string `json:"name"`
}

func Start() {
	r := gin.Default()
	m := manager.New()

	r.GET("/ping", ping)

	r.POST("/functions", addFunctionHandler(m))
	r.POST("/functions/:fn", invokeHandler(m))

	// Listen and serve on localhost
	r.Run()
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func invokeHandler(m *manager.Manager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		fn := c.Param("fn")
		fmt.Println("Executing function: " + fn)
		m.RunContainer(fn)

		c.JSON(200, gin.H{
			"success": "true",
		})
	})
}

func addFunctionHandler(m *manager.Manager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		data, err := fnDataFromReq(c)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		}

		fmt.Println("Creating function...")

		dir, err := dirFromBase64Data(data.Name, data.File)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		}

		tag := data.Name
		m.BuildImage(dir, tag)

		// remove temporary directory used to build the image
		os.RemoveAll(dir)

		c.JSON(200, gin.H{
			"invoke": c.Request.Host + "/functions/" + tag,
		})
	})
}

func fnDataFromReq(c *gin.Context) (FnData, error) {
	rawData, err := c.GetRawData()
	if err != nil {
		return FnData{}, err
	}

	var fnData FnData
	err = json.Unmarshal(rawData, &fnData)
	if err != nil {
		return FnData{}, err
	}
	return fnData, nil
}

func dirFromBase64Data(fnName string, base64Data string) (string, error) {
	filename, err := writeDataToZip(fnName, base64Data)
	if err != nil {
		return "", err
	}

	// unzip the file into a directory called the function name
	dir := fnName
	_, err = Unzip(filename, dir)
	if err != nil {
		return "", err
	}

	// remove the zip file now that the directory has been extracted
	err = os.Remove(filename)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func writeDataToZip(fnName string, fnData string) (string, error) {
	filename := fnName + ".zip"

	zipFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fnData))

	// write decoded data to zip file
	io.Copy(zipFile, decoder)
	return filename, nil
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
