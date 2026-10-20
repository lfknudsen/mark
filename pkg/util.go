package pkg

import (
	"encoding/json/v2"
	"os"
	"path/filepath"
)

// Convert marshals in into JSON, then unmarshals into out.
// Both in and out should be pointers.
func Convert(in any, out any) error {
	intermediate, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(intermediate, out)
}

// CapLength limits the length of a string. If above the limit, it
// replaces the last three characters with "...".
func CapLength(str string, limit int) string {
	if len(str) > limit+3 {
		str = str[:limit] + "..."
	}
	return str
}

func ToJSON(structure any, outputFile string) error {
	// Write to a new JSON file.
	err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	return json.MarshalWrite(out, structure)
}
