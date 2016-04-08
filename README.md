structreflect
==========
[![Build Status](https://travis-ci.org/josephspurrier/structreflect.svg)](https://travis-ci.org/josephspurrier/structreflect) [![Coverage Status](https://coveralls.io/repos/josephspurrier/structreflect/badge.svg)](https://coveralls.io/r/josephspurrier/structreflect) [![GoDoc](https://godoc.org/github.com/josephspurrier/structreflect?status.svg)](https://godoc.org/github.com/josephspurrier/structreflect)

This packages allows you to pass a struct to Generate() and it will return a byte array and the import paths for the struct. This is useful for code generation.

## Usage

Below is a simple example of how to generate the User struct and save it to a file.

```go
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"time"
	
	"github.com/josephspurrier/structreflect"
)

// User table contains the user information
type User struct {
	Id         uint32    `db:"id"`
	First_name string    `db:"first_name"`
	Last_name  string    `db:"last_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Status_id  uint8     `db:"status_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

func main() {
	var err error
	var buffer bytes.Buffer

	// Generate the User struct
	userStruct, userImport, err := structreflect.Generate(User{})
	if err != nil {
		log.Println(err)
		return
	}
	
	log.Println(userImport)
	
	// Write to a file
	err = ioutil.WriteFile("model.go", userStruct, 0644)
	if err != nil {
		log.Println("File save error")
		return
	}
}
```

The file, model.go, will contain:

```
type User struct {
	Id         uint32    `db:"id"`
	First_name string    `db:"first_name"`
	Last_name  string    `db:"last_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Status_id  uint8     `db:"status_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}
```

Take a look at example_test.go for ideas on how to create the rest of the file.