package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidName = errors.New("invalid profile name")
	ErrExists      = errors.New("profile already exists")
	ErrNotFound    = errors.New("profile not found")
)

func Create(name string, p Profile) error {
	path, err := profilePath(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%w: %q", ErrExists, name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func Get(name string) (Profile, error) {
	path, err := profilePath(name)
	if err != nil {
		return Profile{}, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Profile{}, fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	if err != nil {
		return Profile{}, err
	}

	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Profile{}, err
	}

	return p, nil
}

func List() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) == ".yaml" {
			names = append(names, strings.TrimSuffix(name, ".yaml"))
		}
	}

	sort.Strings(names)
	return names, nil
}

func Delete(name string) error {
	path, err := profilePath(name)
	if err != nil {
		return err
	}

	if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %q", ErrNotFound, name)
	} else if err != nil {
		return err
	}

	return nil
}

func profilePath(name string) (string, error) {
	if strings.TrimSpace(name) == "" ||
		strings.Contains(name, "/") ||
		strings.Contains(name, "\\") ||
		strings.Contains(name, "..") {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, name)
	}

	return name + ".yaml", nil
}
