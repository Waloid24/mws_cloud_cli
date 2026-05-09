package profile

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCreateCreatesProfileFile(t *testing.T) {
	chdirTemp(t)

	err := Create("test", Profile{User: "example", Project: "new-project"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if _, err := os.Stat("test.yaml"); err != nil {
		t.Fatalf("expected test.yaml to exist: %v", err)
	}
}

func TestCreateWritesUserAndProjectFields(t *testing.T) {
	chdirTemp(t)

	err := Create("test", Profile{User: "example", Project: "new-project"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	data, err := os.ReadFile("test.yaml")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "user: example") {
		t.Fatalf("expected user field in YAML, got:\n%s", content)
	}
	if !strings.Contains(content, "project: new-project") {
		t.Fatalf("expected project field in YAML, got:\n%s", content)
	}
}

func TestGetReadsExistingProfile(t *testing.T) {
	chdirTemp(t)

	want := Profile{User: "example", Project: "new-project"}
	if err := Create("test", want); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := Get("test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != want {
		t.Fatalf("Get() = %+v, want %+v", got, want)
	}
}

func TestListReturnsCreatedProfiles(t *testing.T) {
	chdirTemp(t)

	for name, p := range map[string]Profile{
		"dev":  {User: "dev-user", Project: "dev-project"},
		"prod": {User: "prod-user", Project: "prod-project"},
	} {
		if err := Create(name, p); err != nil {
			t.Fatalf("Create(%q) error = %v", name, err)
		}
	}

	got, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []string{"dev", "prod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %#v, want %#v", got, want)
	}
}

func TestDeleteRemovesProfileFile(t *testing.T) {
	chdirTemp(t)

	if err := Create("test", Profile{User: "example", Project: "new-project"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := Delete("test"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := os.Stat("test.yaml"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected test.yaml to be removed, stat error = %v", err)
	}
}

func TestGetReturnsErrorForMissingProfile(t *testing.T) {
	chdirTemp(t)

	_, err := Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestDeleteReturnsErrorForMissingProfile(t *testing.T) {
	chdirTemp(t)

	err := Delete("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestCreateReturnsErrorIfProfileAlreadyExists(t *testing.T) {
	chdirTemp(t)

	p := Profile{User: "example", Project: "new-project"}
	if err := Create("test", p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := Create("test", p)
	if !errors.Is(err, ErrExists) {
		t.Fatalf("Create() error = %v, want ErrExists", err)
	}
}

func chdirTemp(t *testing.T) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir(%q) error = %v", tempDir, err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory to %q: %v", previousDir, err)
		}
	})

}
