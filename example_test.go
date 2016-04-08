package structreflect

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"time"
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

// UserStatus table contains different user status (active/inactive)
type UserStatus struct {
	Id         uint8     `db:"id"`
	Status     string    `db:"status"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// Logic will generate the struct
func main() {
	var err error
	var buffer bytes.Buffer

	// Generate the User struct
	userStruct, userImport, err := Generate(User{})
	if err != nil {
		log.Println(err)
		return
	}

	// Generate the UserStatus struct
	userStatus, statusImport, err := Generate(UserStatus{})
	if err != nil {
		log.Println(err)
		return
	}

	// Write the package declaration
	buffer.Write([]byte("package model\n"))

	// Write the imports
	buffer.Write([]byte("import (\n"))
	for _, i := range userImport {
		buffer.Write([]byte(fmt.Sprintf(`"%s"`, i) + "\n"))
	}
	for _, i := range statusImport {
		buffer.Write([]byte(fmt.Sprintf(`"%s"`, i) + "\n"))
	}
	buffer.Write([]byte(")\n"))

	// Wrote the structs
	buffer.Write(userStruct)
	buffer.Write([]byte("\n\n"))
	buffer.Write(userStatus)
	buffer.Write([]byte("\n\n"))

	// Format the buffer using go fmt
	out, err := format.Source(buffer.Bytes())
	if err != nil {
		log.Println(err)
		return
	}

	// Write the file
	err = ioutil.WriteFile("model.go", out, 0644)
	if err != nil {
		log.Println("File save error")
		return
	}
}
