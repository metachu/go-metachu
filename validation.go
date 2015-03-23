package main

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"os"
	"path/filepath"
)

func (fa FileActionJson) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	// if strings.Contains(cf.Message, "Go needs generics") {
	//     errors = append(errors, binding.Error{
	//         FieldNames:     []string{"message"},
	//         Classification: "ComplaintError",
	//         Message:        "Go has generics. They're called interfaces.",
	//     })
	// }

	return errors
}

func ValidatePath(input_path string) (string, string, error) {
	raw_path, _ := filepath.Abs(filepath.Join(Config.RootPath, input_path))
	raw_path = filepath.Clean(raw_path)
	rel_path := filepath.Clean(raw_path[len(Config.RootPath):])
	if rel_path == "." {
		rel_path = "/"
	}
	_, err := os.Stat(raw_path)
	if err != nil {
		return "", "", err
	}
	return raw_path, rel_path, nil
}
