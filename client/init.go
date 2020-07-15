package client

import (
	"io/ioutil"
)

func Init() error {
	err := writeIndexFile()
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
  "license": "ISC"
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

CMD [ "node", "index.js" ]`

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
