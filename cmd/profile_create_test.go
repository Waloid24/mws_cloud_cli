package cmd

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestProfileCreateCreatesYAMLFile(t *testing.T) {
	chdirTemp(t)

	output, err := executeCommand("profile", "create", "--name=test", "--user=example", "--project=new-project")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}
	if output != "profile \"test\" created\n" {
		t.Fatalf("output = %q, want create success message", output)
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

func TestProfileCreateRequiresName(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "create", "--user=example", "--project=new-project")
	if err == nil || !strings.Contains(err.Error(), `required flag "name" not set`) {
		t.Fatalf("error = %v, want required name error", err)
	}
}

func TestProfileCreateRequiresUser(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "create", "--name=test", "--project=new-project")
	if err == nil || !strings.Contains(err.Error(), `required flag "user" not set`) {
		t.Fatalf("error = %v, want required user error", err)
	}
}

func TestProfileCreateRequiresProject(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "create", "--name=test", "--user=example")
	if err == nil || !strings.Contains(err.Error(), `required flag "project" not set`) {
		t.Fatalf("error = %v, want required project error", err)
	}
}

func TestProfileCreateFailsWhenProfileAlreadyExists(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "create", "--name=test", "--user=example", "--project=new-project")
	if err != nil {
		t.Fatalf("first create error = %v", err)
	}

	_, err = executeCommand("profile", "create", "--name=test", "--user=example", "--project=new-project")
	if err == nil || !strings.Contains(err.Error(), "profile already exists") {
		t.Fatalf("error = %v, want profile already exists error", err)
	}
}

func TestProfileGetPrintsExistingProfile(t *testing.T) {
	chdirTemp(t)
	writeFile(t, "test.yaml", "user: example\nproject: new-project\n")

	output, err := executeCommand("profile", "get", "--name=test")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}

	want := "user: example\nproject: new-project\n"
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestProfileGetRequiresName(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "get")
	if err == nil || !strings.Contains(err.Error(), `required flag "name" not set`) {
		t.Fatalf("error = %v, want required name error", err)
	}
}

func TestProfileGetReturnsErrorForMissingProfile(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "get", "--name=missing")
	if err == nil || !strings.Contains(err.Error(), "profile not found") {
		t.Fatalf("error = %v, want profile not found error", err)
	}
}

func TestProfileListPrintsCreatedProfiles(t *testing.T) {
	chdirTemp(t)
	writeFile(t, "dev.yaml", "user: andrey\nproject: alpha\n")
	writeFile(t, "prod.yaml", "user: andrey\nproject: beta\n")

	output, err := executeCommand("profile", "list")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}

	want := "dev\nprod\n"
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestProfileListReturnsSortedProfiles(t *testing.T) {
	chdirTemp(t)
	writeFile(t, "prod.yaml", "user: andrey\nproject: beta\n")
	writeFile(t, "dev.yaml", "user: andrey\nproject: alpha\n")
	writeFile(t, "test.yaml", "user: andrey\nproject: gamma\n")

	output, err := executeCommand("profile", "list")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}

	got := strings.Fields(output)
	want := []string{"dev", "prod", "test"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("profiles = %#v, want %#v", got, want)
	}
}

func TestProfileListIgnoresNonYAMLFiles(t *testing.T) {
	chdirTemp(t)
	writeFile(t, "dev.yaml", "user: andrey\nproject: alpha\n")
	writeFile(t, "notes.txt", "not a profile\n")
	writeFile(t, "prod.yml", "user: andrey\nproject: beta\n")

	output, err := executeCommand("profile", "list")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}

	want := "dev\n"
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestProfileListWorksWithEmptyDirectory(t *testing.T) {
	chdirTemp(t)

	output, err := executeCommand("profile", "list")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}
	if output != "" {
		t.Fatalf("output = %q, want empty output", output)
	}
}

func TestProfileDeleteRemovesExistingProfile(t *testing.T) {
	chdirTemp(t)
	writeFile(t, "test.yaml", "user: example\nproject: new-project\n")

	output, err := executeCommand("profile", "delete", "--name=test")
	if err != nil {
		t.Fatalf("executeCommand() error = %v", err)
	}
	if output != "profile \"test\" deleted\n" {
		t.Fatalf("output = %q, want delete success message", output)
	}

	if _, err := os.Stat("test.yaml"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected test.yaml to be removed, stat error = %v", err)
	}
}

func TestProfileDeleteRequiresName(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "delete")
	if err == nil || !strings.Contains(err.Error(), `required flag "name" not set`) {
		t.Fatalf("error = %v, want required name error", err)
	}
}

func TestProfileDeleteReturnsErrorForMissingProfile(t *testing.T) {
	chdirTemp(t)

	_, err := executeCommand("profile", "delete", "--name=missing")
	if err == nil || !strings.Contains(err.Error(), "profile not found") {
		t.Fatalf("error = %v, want profile not found error", err)
	}
}

func executeCommand(args ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	rootCmd := NewRootCommand()
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return stdout.String() + stderr.String(), err
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
			t.Errorf("restore working directory to %q: %v", previousDir, err)
		}
	})
}

func writeFile(t *testing.T, name, content string) {
	t.Helper()

	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", name, err)
	}
}
