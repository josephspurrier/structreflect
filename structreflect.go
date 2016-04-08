package structreflect

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrorNotStruct means the value passed is not a struct
	ErrNotStruct = errors.New("Not a struct")
)

// Generate returns byte array of the struct, a string array of the import paths
// and an error if there was a problem
func Generate(s interface{}) ([]byte, []string, error) {
	// Buffer for struct
	var buffer bytes.Buffer

	// Import paths
	var imports []string

	// Value of the struct for reflection
	sValue := reflect.ValueOf(s)

	// Write the beginning of the struct
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", sValue.Type().Name()))

	// Prevent running on types other than struct
	if sValue.Type().Kind() != reflect.Struct {
		return nil, nil, ErrNotStruct
	}

	// Traverse the struct
	traverseStructField(&buffer, &imports, sValue.Type())

	return buffer.Bytes(), imports, nil
}

func traverseStructField(buffer *bytes.Buffer, imports *[]string, sValue reflect.Type) {
	// Look through each field of the struct
	for i := 0; i < sValue.NumField(); i++ {
		// Field
		f := sValue.Field(i)

		// Package/import path
		pkgPath := f.Type.PkgPath()

		// Get real path from pointer
		if f.Type.Kind() == reflect.Ptr {
			pkgPath = f.Type.Elem().PkgPath()
		}

		// Determine if field requires package name
		external := false
		if len(pkgPath) > 0 && pkgPath != "main" {
			external = true

			// Store the import path
			*imports = appendIfMissing(*imports, pkgPath)
		}

		// Only struct can be anonymous
		if f.Anonymous {
			// Pointer
			if f.Type.Kind() == reflect.Ptr {
				// External package
				if external {
					buffer.WriteString(fmt.Sprintf("\t%s", f.Type))
				} else {
					buffer.WriteString(fmt.Sprintf("\t*%s", f.Name))
				}
			} else { // No pointer
				if external {
					buffer.WriteString(fmt.Sprintf("\t%s", f.Type))
				} else {
					buffer.WriteString(fmt.Sprintf("\t%s", f.Name))
				}

			}
		} else { // All other variable types
			// Write the name of the variable
			buffer.WriteString(fmt.Sprintf("\t%s", f.Name))

			// Nested struct
			if f.Type.Kind() == reflect.Struct && len(f.Type.Name()) == 0 {
				buffer.WriteString(" struct {\n")
				traverseStructField(buffer, imports, f.Type)
			} else if f.Type.Kind() == reflect.Slice {
				buffer.WriteString(fmt.Sprintf(" %s", f.Type))
			} else if f.Type.Kind() == reflect.Map {
				buffer.WriteString(fmt.Sprintf(" %s", f.Type))
			} else if f.Type.Kind() == reflect.Interface {
				buffer.WriteString(fmt.Sprintf(" %s", f.Type))
			} else if f.Type.Kind() == reflect.Ptr {
				// Struct
				if strings.Contains(f.Type.String(), ".") {
					if external {
						buffer.WriteString(fmt.Sprintf(" %s", f.Type))
					} else {
						buffer.WriteString(fmt.Sprintf(" *%s", strings.Split(f.Type.String(), ".")[1]))
					}
				} else {
					buffer.WriteString(fmt.Sprintf(" %s", f.Type))
				}
			} else {
				if external {
					buffer.WriteString(fmt.Sprintf(" %s", f.Type))
				} else {
					buffer.WriteString(fmt.Sprintf(" %s", f.Type.Name()))
				}

			}
		}

		tag := string(f.Tag)

		// If a tag exists
		if len(tag) > 0 {
			// If the tag contains both double quotes and backticks
			if strings.Contains(tag, "`") && strings.Contains(tag, `"`) {
				buffer.WriteString(fmt.Sprintf(` "%s"`, strings.Replace(tag, `"`, `\"`, -1)))
			} else if strings.Contains(tag, "`") {
				buffer.WriteString(fmt.Sprintf(` "%s"`, tag))
			} else {
				buffer.WriteString(fmt.Sprintf(" `%s`", tag))
			}

		}

		buffer.WriteString("\n")
	}
	buffer.WriteString("}")
}

// Source: http://stackoverflow.com/a/9561388
func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
