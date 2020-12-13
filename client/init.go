package client

import (
	"io/ioutil"
)

func Init() error {
	err := writeServerFile()
	if err != nil {
		return err
	}
	err = writeIndexFile()
	if err != nil {
		return err
	}
	err = writePackageFile()
	if err != nil {
		return err
	}
	err = writeDockerfile()
	if err != nil {
		return err
	}
	err = writeYAMLFile()
	if err != nil {
		return err
	}
	return nil
}

func writeServerFile() error {
	serverContents := `const express = require('express')
const app = express()
const bodyParser = require('body-parser')
app.use(bodyParser.json())

app.post('/invoke', (req, res) => {
	console.log('INVOKE TESTING 123')
	res.json({my: 'response', hello: 'world!'})
	app.close()
})

const port = 8080
app.listen(port, () => {})`

	err := ioutil.WriteFile("server.js", []byte(serverContents), 0755)
	if err != nil {
		return err
	}
	return nil

}

func writeIndexFile() error {
	indexContents := `if (require.main === module) {
  console.log('HELLO WORLD!')
}`

	err := ioutil.WriteFile("index.js", []byte(indexContents), 0755)
	if err != nil {
		return err
	}
	return nil
}

func writePackageFile() error {
	packageContents := `
{
  "name": "example",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "dependencies": {
    "body-parser": "^1.19.0",
    "express": "^4.17.1"
  }
}`

	err := ioutil.WriteFile("package.json", []byte(packageContents), 0755)
	if err != nil {
		return err
	}
	return nil
}

func writeDockerfile() error {
	dockerfileContents := `FROM node:12

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
COPY package*.json ./

RUN npm install

# Bundle app source
COPY . .

EXPOSE 8080
CMD [ "node", "server.js" ]`

	err := ioutil.WriteFile("Dockerfile", []byte(dockerfileContents), 0755)
	if err != nil {
		return err
	}
	return nil
}

func writeYAMLFile() error {
	yamlContents := `language: node`

	err := ioutil.WriteFile("fn.yaml", []byte(yamlContents), 0755)
	if err != nil {
		return err
	}
	return nil
}
