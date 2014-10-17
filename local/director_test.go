package local

import (
	"path/filepath"
	"testing"
)

func TestLoadBoshConfig(t *testing.T) {
	configPath, err := filepath.Abs("../testhelpers/fixtures/bosh_config.yml")
	if err != nil {
		t.Fatalf("error raised: %s", err)
	}

	config, err := LoadBoshConfig(configPath)
	if err != nil {
		t.Fatalf("error raised: %s", err)
	}
	if config == nil {
		t.Fatalf("expect BoshConfig not nil")
	}
	if config.Name != "Bosh Lite Director" {
		t.Fatalf("config.Name not correct: '%s'", config.Name)
	}
	if config.Authentication["https://192.168.50.4:25555"].Username != "admin" {
		t.Fatalf("config.Authentication not correct: '%#v'", config.Authentication)
	}
}
