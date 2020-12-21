package client

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rknizzle/faas/internal/models"
)

const (
	serverFile = `const express = require('express')
const app = express()
const bodyParser = require('body-parser')
app.use(bodyParser.json())

// load in the function code
const fn = require('./index.js')

app.post('/invoke', (req, res) => {
  function cb(output) {
    res.json(output)
    // exit the container after finishing running the function
    server.close()
  }

  fn(req.body, cb)
})

const port = 8080
const server = app.listen(port, () => {})`

	dockerFile = `FROM node:12

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
COPY package*.json ./

RUN npm install

# Bundle app source
COPY . .

EXPOSE 8080
CMD [ "node", "server.js" ]`
)

func Build() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Get all files in directory
	var fileList []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fileList = append(fileList, f.Name())
	}

	// remove node_modules if it exists from list of files to send to server
	fileList = remove(fileList, "node_modules")

	// get the name of the current directory
	name := filepath.Base(path)
	output := name + ".zip"

	// Combine all files into a zip
	err = ZipFiles(output, fileList)
	if err != nil {
		return "", err
	}

	// Open new zip file
	f, err := os.Open(output)
	if err != nil {
		return "", err
	}

	// Read entire zip file into byte slice
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Encode zip as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	// format the function data to send to the server
	fd := &models.FnData{File: encoded, Name: name}

	// convert function data to JSON body to send in HTTP request to server
	funcByte, _ := json.Marshal(fd)
	funcReader := bytes.NewReader(funcByte)

	client := &http.Client{}
	r, err := http.NewRequest("POST", "http://localhost:5555/functions", funcReader)
	if err != nil {
		return "", err
	}

	// send request to server to submit new function
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// get the JSON response
	var result map[string]string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	f.Close()

	// remove the zip file
	err = os.Remove(output)
	if err != nil {
		return "", err
	}

	// return the invocation name
	return filepath.Base(result["invoke"]), nil
}

func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// add a file to the zip without the file having to exist on disk. Pass in the contents of the file
// as a []byte and the name to use as the filename
func addInMemoryFileToZip(zipWriter *zip.Writer, filename string, content []byte) error {
	zipFile, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}
	_, err = zipFile.Write(content)
	if err != nil {
		return err
	}
	return nil
}

// remove a specific file from a list of files
func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}
