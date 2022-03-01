package handy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SaveFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	content = append(content, []byte("\n")...)

	switch {
	case strings.HasSuffix(path, ".yaml"):
		if err := os.WriteFile(path, content, 0644); err != nil {
			return fmt.Errorf("error saving file(%s): %v", path, err)
		}
	case strings.HasSuffix(path, ".json"):
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, content, "", "  "); err != nil {
			return fmt.Errorf("error json indenting: %v", err)
		}
		if err := os.WriteFile(path, pretty.Bytes(), 0644); err != nil {
			return fmt.Errorf("error saving file(%s): %v", path, err)
		}
	}

	return nil
}
