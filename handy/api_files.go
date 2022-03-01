package handy

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type APIFile struct {
	Name    string
	Version string
	Path    string
	Format  string
	// Consumer takes the responsibility to close the file
	File io.ReadCloser
}

var fileRegexp *regexp.Regexp

func init() {
	fileRegexp = regexp.MustCompile(`.*\.(yaml|json)$`)
}

// APIFiles read files from paths and extract titles, filenames and formats from the content.
// Paths could be either files or directories.
// To make it simple, only .yaml and .json are supported.
// It relies on file extension for properly parsing.
func APIFiles(paths []string) ([]APIFile, error) {
	var apis []APIFile

	for _, path := range paths {
		apis = append(apis, walk(path)...)
	}

	var g errgroup.Group

	for i := range apis {
		i := i

		g.Go(func() error {
			err := info(&apis[i])
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return apis, nil
}

// walk walks through the path and returns all .yaml/.json files.
func walk(root string) []APIFile {
	var apis []APIFile

	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		match := fileRegexp.FindStringSubmatch(path)
		if len(match) < 2 {
			return nil
		}

		apis = append(apis, APIFile{
			Path:   path,
			Format: match[1],
		})

		return nil
	})

	return apis
}

func info(api *APIFile) error {
	file, err := os.Open(api.Path)
	if err != nil {
		return err
	}

	var v doc

	switch api.Format {
	case "json":
		if err := json.NewDecoder(file).Decode(&v); err != nil {
			return err
		}
	case "yaml":
		if err := yaml.NewDecoder(file).Decode(&v); err != nil {
			return err
		}
	}

	file.Seek(0, io.SeekStart)
	api.File = file
	ss := strings.Split(v.Info.Title, "/")
	api.Name = ss[len(ss)-1]
	api.Version = v.Info.Version

	return nil
}

type doc struct {
	Info struct {
		Title   string
		Version string
	}
}
