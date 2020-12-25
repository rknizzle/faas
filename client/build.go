package client

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// make the current directory into a tar file buffer
	content, err := Tar(path)
	if err != nil {
		return "", err
	}

	// Encode tar as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	// get the name of the current directory
	name := filepath.Base(path)

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

	var buildErr struct {
		Message string `json:"message"`
	}

	var buildRes struct {
		Invoke string `json:"invoke"`
	}

	// get the JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// check status code here and unmarshal into the appropriate struct then return the correct value
	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &buildRes)
		if err != nil {
			return "", err
		}

		// return the invocation name
		return filepath.Base(buildRes.Invoke), nil
	} else {
		err = json.Unmarshal(body, &buildErr)
		if err != nil {
			return "", err
		}
		return "", errors.New(buildErr.Message)
	}
}

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer; the purpose for accepting multiple writers is to allow
// for multiple outputs (for example a file, or md5 hash)
func Tar(src string) ([]byte, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// walk path
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// return on non-regular files
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})

	// TODO: these files required for building should be added on the backend instead. That way the
	// unaltered tar containing the function code can be uploaded to storage before adding the extra
	// files required for building.
	// For now just place the extra required files for building the nodejs image here and pass it
	// right into the Docker ImageBuild call.

	// place the server.js file for handling the HTTP request from the runner into the tar
	err := addInMemoryFileToTar(tw, "server.js", serverFile)
	if err != nil {
		return nil, err
	}

	// place the Dockerfile for building the image into the tar
	err = addInMemoryFileToTar(tw, "Dockerfile", dockerFile)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// addInMemoryFileToTar adds a file to a tar without the file having to exist on disk. Just pass in
// a filenae and the file contents
func addInMemoryFileToTar(tw *tar.Writer, filename string, fileContents string) error {
	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(fileContents)),
	}

	err := tw.WriteHeader(hdr)
	if err != nil {
		return err
	}

	_, err = tw.Write([]byte(fileContents))
	if err != nil {
		return err
	}

	return nil
}
