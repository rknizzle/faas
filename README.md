# FAAS

A simple Functions-as-a-Service(FaaS) platform using Docker containers for the purpose of experimentation/learning


### Install CLI tool
```
go get -u github.com/rknizzle/faas/cmd/faas
```

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
