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
	err = writeYAMLFile()
	if err != nil {
		return err
	}
	return nil
}

func writeIndexFile() error {
	indexContents := `// Put the function logic below.
// context contains the input data and the callback returns a result to the caller
module.exports = (context, cb) => {
  return cb({hello: "world"})
}
`

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

func writeYAMLFile() error {
	yamlContents := `language: node`

	err := ioutil.WriteFile("fn.yaml", []byte(yamlContents), 0755)
	if err != nil {
		return err
	}
	return nil
}
