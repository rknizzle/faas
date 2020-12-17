# FAAS 
![example workflow name](https://github.com/rknizzle/faas/workflows/Test/badge.svg)

### A simple Functions-as-a-Service(FaaS) platform using Docker containers for the purpose of experimentation/learning

- [Dependencies](#dependencies)
- [Install](#install-cli-tool)
- [Demo](#quick-local-demo)
- [Example](#example)
- [API](#rest-api)
  - [Add Function](#add-function)
  - [Invoke Function](#invoke-function)
- [Roadmap](#roadmap)


# Dependencies
- Only works in Linux and tested on Ubuntu
- Must have Docker daemon installed and running 

# Install CLI tool
NOTE: Currently only Linux is supported due to extra effort required to interact with Docker
containers on Mac OS and Windows because of the Docker virtual machine

The latest Linux binary can be found [here](https://github.com/rknizzle/faas/releases)

# Quick Local Demo:

###  Start the server:
`faas start`

### Create and invoke a function
In another terminal:
```sh
mkdir examplefn
cd examplefn
faas init
faas build
faas invoke examplefn
```

Output: `Response: {"hello":"world"}`

# Example
If you are running locally start the faas server first: `faas start`  
Next, create a new working directory for your function. cd into it and then run faas init.  
After running faas init, open index.js in your code editor to modify your function.  

Here is an example function that takes in an array of numbers are returns the sum:  
```js
// takes an array of numbers and returns the sum
module.exports = (context, cb) => {
  let sum
  for (let num of context.numbers) {
    sum += num
  }

  // respond to the caller with the sum of numbers
  return cb({sum: sum})
}
```
Run `faas build` in your functions working directory to deploy the function to the server.  

Now lets put some input data into a JSON file.  
### input.json:  
```json
{
  "numbers": [6,4,13]
}
```

Now invoke the function and pass it your array of numbers:  
`faas invoke fn -d input.json`  

Output: `Response: {"sum":23}`

# REST API
# Add Function

Adds a new function and returns the url used to invoke it

**URL** : `/functions/`

**Method** : `POST`

**Data constraints**

Provide name of the function and a base64 encoded zip file containing the directory with the
function code inside

```json
{
    "name": "<string>",
    "file": "<base64 encoded zip file containing function code>"
}
```

**Data example** All fields must be sent.

```json
{
    "name": "add",
    "file": "UEsDBBQACAAIAIgEkVEAAAAAAAAAAAAAAAAKAAkARG9ja2VyZmlsZVVUBQABAajaXySOwUrGMBAG7/sUS7yJJNVT8WhboUhNiYiKeAjJIq01Cbut4Nv/lNxnvm8enZ0w5Uj3t3cAV9gx+Z3Ql4JxYQp75n94s+6pHx2aQ9gIB+NLOeExye63rdJUKEVKYSGBzs4fWHz48d90rVfJCbUBcK/PmMovLtU7Jx6OFLf6J/ngQNXVqAGG99m+DNg2bQPd1OMnqrNU3aAS4j9ivYrCr0sAAAD//1BLBwgkCHCtpQAAAMMAAABQSwMEFAAIAAgAiASRUQAAAAAAAAAAAAAAAAcACQBmbi55YW1sVVQFAAEBqNpfyknMSy9NTE+1UsjLT0kFBAAA//9QSwcI+HW5hBQAAAAOAAAAUEsDBBQACAAIAIgEkVEAAAAAAAAAAAAAAAAIAAkAaW5kZXguanNVVAUAAQGo2l80jU3KgzAURedZxcWRwoeZf2DX0C3k57WGvuZJ8oKCuPdSqaMLh8M91uLeFDoTHi0HTZLB8kwBnljW0ViLIFlp03NdyvW0U16aIjp1cDmiVTpxcMzehRdUUEhbyXAoVBvrF10KFfOW2JhG2hYpWjGh/3X+EPyA6Ybd4PoIvt9nYpZ/dKsUjt0xmMN8AgAA//9QSwcI5xQrk5AAAAC+AAAAUEsDBBQACAAIAIgEkVEAAAAAAAAAAAAAAAAMAAkAcGFja2FnZS5qc29uVVQFAAEBqNpfTJAxTwQhEIX7/RUTiquUHImJcVtjYW3paYLwzI3eAplB3Yu5/26AM1ryPZj5HtP3RGSSX2BmMlj9Ug4wFw1+QpRzatzZrd0OGqFBuNRzMuDiuZ84Raz2TQcdF9XM1JYQmQqtfU3YZ9qZO5EsM6VMLSAtCPzKiDtDmw1h5UrOTESnPu0dx68ssY17fOrEf9R9lj+LAwck7UXuH25/dQtSRAqMfyIvOR4vixdFf/7srLs5F6T2C0Wg2pMr667tsJhOPwEAAP//UEsHCGpcPl/GAAAALQEAAFBLAwQUAAgACACIBJFRAAAAAAAAAAAAAAAACQAJAHNlcnZlci5qc1VUBQABAajaX0yQQW7rMAxE9zrF7CwB/nKW2fifoVdwbLpRGlAKKQcpCt+9kOygXgl4GI4eOUbWDHolIVX0EHosQcg2O2qc2SJDSujfQfumlzh9fwyiJMfZQv+lihtnhpT8omT/sv6mka1zxnQd7nGYEBj5SpgXHnOIjDFOtH8x87Had4EnevlbMavVKWq2TRf4Gb+oaWGFHi2E1KH/jx+DQ+vFxiWnJbvKUVKby44r7DrQK+QqNEbOQ2ASDHMmwRw46DXwJ2RhLu9Ru44ryZPEj/eoZEvjaooDFy9fbtBivDizOrMvmKJk9DifzqedbBXoy9H9PWgmtiXVwm5Lre43AAD//1BLBwgrY3WL/AAAALgBAABQSwECFAMUAAgACACIBJFRJAhwraUAAADDAAAACgAJAAAAAAAAAAAA7YEAAAAARG9ja2VyZmlsZVVUBQABAajaX1BLAQIUAxQACAAIAIgEkVH4dbmEFAAAAA4AAAAHAAkAAAAAAAAAAADtgeYAAABmbi55YW1sVVQFAAEBqNpfUEsBAhQDFAAIAAgAiASRUecUK5OQAAAAvgAAAAgACQAAAAAAAAAAAO2BOAEAAGluZGV4LmpzVVQFAAEBqNpfUEsBAhQDFAAIAAgAiASRUWpcPl/GAAAALQEAAAwACQAAAAAAAAAAAO2BBwIAAHBhY2thZ2UuanNvblVUBQABAajaX1BLAQIUAxQACAAIAIgEkVErY3WL/AAAALgBAAAJAAkAAAAAAAAAAADtgRADAABzZXJ2ZXIuanNVVAUAAQGo2l9QSwUGAAAAAAUABQBBAQAATAQAAAAA"
}
```

## Success Response

**Condition** : If the function data is valid.

**Code** : `200 OK`

**Content example**

```json
{
    "invoke": "http://localhost:5555/functions/:name",
}
```

## Error Responses

**Code** : `400 BAD REQUEST`

**Content example**

```json
{
    "message": "Failed to deploy function"
}
```

# Invoke Function

Invokes a function

**URL** : `/functions/:name`

**Method** : `POST`

**Data constraints**

Provide the input data to the function in the request body

```json
{
  "anyInputFields": "anyValues"
}
```

**Data example**

```json
{
    "numbers": [7,33,12,18],
}
```

## Success Response

**Condition** : If the function exists.

**Code** : `200 OK`

**Content example**

```json
{
    "response": "[JSON function result]",
}
```

## Error Responses

**Code** : `400 BAD REQUEST`

**Content example**

```json
{
    "message": "[Error message]"
}
```

# Roadmap:
- [x] Pass input data to the function container and recieve a response to return to the caller
- [x] Run unit tests automatically on push/PR via Github Actions
- [x] Automatically deploy a new binary release on PR merge into master via GitHub Actions
- [ ] Store function name/images and users/apps in a datastore
- [ ] Push deployed images to a docker registry and pull images when invoked
- [ ] Decouple the gateway and runners into separate binaries so faas can be scaled across multiple
  machines
- [ ] Support more programming languages
