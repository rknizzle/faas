# FAAS
![Tests](https://github.com/rknizzle/faas/workflows/Go/badge.svg)

A simple Functions-as-a-Service(FaaS) platform using Docker containers for the purpose of experimentation/learning


### Dependencies
Must have Docker daemon installed and running 

### Install CLI tool
Releases for OSX and Linux can be found [here](https://github.com/rknizzle/faas/releases)

### Local demo:

#####  Start the server:
`faas start`

##### Create and invoke a function
In another terminal:
```
mkdir examplefn
cd examplefn
faas init
faas build
faas invoke examplefn
```

Output: `HELLO WORLD!`
