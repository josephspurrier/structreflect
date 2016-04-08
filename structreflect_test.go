package structreflect

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/josephspurrier/structreflect/test/ball"
)

// *****************************************************************************
// Structs
// *****************************************************************************

// Foo test struct
type Foo struct {
	A    string  //string
	B    *string // pointer
	c    string  // not exported
	d    *string
	E    int
	f    *int
	G    uint `structs:"y"`
	H    uint `json:"h"`
	I    bool
	j    bool `xml:"j"` // not exported, with tag
	K    Bar
	L    *Bar
	Buzz `cool`
	*Bar "beans"
	Fizz struct {
		Legs int `json:"legs"`
		*ball.Ball
		bb *ball.Ball
	}
	ball.Ball
	M []string
	n *[]string
	O map[string]interface{}
	p *map[string]interface{}
	q interface{}
	r *interface{}
	s []interface{}
	t *[]interface{}
	u time.Time
	v *time.Time
}

// Buzz test struct
type Buzz struct {
	A *string
	B int
}

// Bar test struct
type Bar struct {
	E string
	F int
	g []string
}

// *****************************************************************************
// Logic Testing
// *****************************************************************************

func TestGenerate(t *testing.T) {
	// Load test file
	filename := "correct.bin"
	path, _ := filepath.Abs("./test/" + filename)
	byteExpected, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename, err)
	}

	// Generate the struct
	byteArray, strTest, err := Generate(Foo{})
	if err != nil {
		log.Println(err)
		return
	}

	// Format the buffer using go fmt
	byteTest, err := format.Source(byteArray)
	if err != nil {
		log.Println(err)
		return
	}

	if !bytes.Equal(byteTest, byteExpected) {
		t.Error("Struct does not match " + filename)
	}

	strExpected := []string{
		"github.com/josephspurrier/structreflect",
		"github.com/josephspurrier/structreflect/test/ball",
		"time",
	}

	if !reflect.DeepEqual(strTest, strExpected) {
		t.Error("Imports paths do not match")
	}
}

func TestBadStruct(t *testing.T) {
	// Attempt to generate the struct
	_, _, err := Generate("string")
	if err != ErrNotStruct {
		t.Error("Generate should fail")
	}
}
