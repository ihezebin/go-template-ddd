package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigurationVariables(t *testing.T) {
	previous := gConfig
	gConfig = &Config{}
	t.Cleanup(func() { gConfig = previous })
	path := filepath.Join(t.TempDir(), "config.json")
	data := `{"service_name": "variable-test", "port": 12345, "mongo_dsn": "mongodb://localhost/${service_name}-server", "logger": {"level": "info", "filename": "log/${service_name}.log"}}`
	for _, name := range []string{"variable-test", "renamed-service"} {
		gConfig = &Config{}
		raw := strings.ReplaceAll(data, "variable-test", name)
		if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
		c, err := Load(path)
		if err != nil {
			t.Fatal(err)
		}
		if c.MongoDsn != "mongodb://localhost/"+name+"-server" || c.Logger == nil || c.Logger.Filename != "log/"+name+".log" {
			t.Fatal("configuration variables did not resolve through project loader")
		}
	}
	if err := os.WriteFile(path, []byte(strings.ReplaceAll(data, "${service_name}", "${missing}")), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("missing variable accepted")
	}
}

func TestExampleConfigurationVariables(t *testing.T) {
	previous := gConfig
	t.Cleanup(func() { gConfig = previous })
	for _, name := range []string{"config.json", "config.toml"} {
		t.Run(name, func(t *testing.T) {
			gConfig = &Config{}
			c, err := Load(name)
			if err != nil {
				t.Fatal(err)
			}
			if c.Logger == nil || c.Logger.Filename != "log/"+c.ServiceName+".log" {
				t.Fatal("example log path did not resolve")
			}
		})
	}
}

func TestConfigurationVariableRules(t *testing.T) {
	previous := gConfig
	t.Cleanup(func() { gConfig = previous })
	for _, tc := range []struct {
		name    string
		data    string
		want    string
		wantErr string
	}{
		{"nested and chained", `service_name = "${logger.level}-${port}"
port = 8080
[logger]
level = "info"
filename = "log/${service_name}.log"`, "log/info-8080.log", ""},
		{"escaped", `[logger]
filename = "log/$${missing}.log"`, "log/${missing}.log", ""},
		{"missing", `service_name = "${missing}"`, "", "unknown config variable"},
		{"cycle", `service_name = "${logger.filename}"
[logger]
filename = "${service_name}"`, "", "cyclic reference"},
		{"object", `service_name = "${logger}"
[logger]
level = "info"`, "", "must be a string, number or boolean"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gConfig = &Config{}
			path := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(path, []byte(tc.data), 0600); err != nil {
				t.Fatal(err)
			}
			c, err := Load(path)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("want error containing %q, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if c.Logger == nil || c.Logger.Filename != tc.want {
				t.Fatalf("want log path %q, got %+v", tc.want, c.Logger)
			}
		})
	}
}
