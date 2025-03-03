/*
Copyright The Ratify Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package authprovider

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type TestAuthProvider struct{}

func (ap *TestAuthProvider) Enabled(ctx context.Context) bool {
	return true
}

func (ap *TestAuthProvider) Provide(ctx context.Context, artifact string) (AuthConfig, error) {
	return AuthConfig{
		Username: "test",
		Password: "testpw",
	}, nil
}

// Checks for correct credential resolution when external docker config
// path is provided
func TestProvide_ExternalDockerConfigPath_ExpectedResults(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatalf("unexpected error when creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, "config.json")
	js := `{
		"auths": {
			"index.docker.io": {
				"auth": "am9lam9lOmhlbGxv"
			}
		}
	}`

	err = ioutil.WriteFile(fn, []byte(js), 0600)
	if err != nil {
		t.Fatalf("unexpected error when writing config file: %v", err)
	}

	defaultProvider := defaultAuthProvider{
		configPath: fn,
	}

	authConfig, err := defaultProvider.Provide(context.Background(), "index.docker.io/v1/test:v1")
	if err != nil {
		t.Fatalf("unexpected error in Provide: %v", err)
	}

	if authConfig.Username != "joejoe" || authConfig.Password != "hello" {
		t.Fatalf("incorrect username %v or password %v returned", authConfig.Username, authConfig.Password)
	}
}

func TestProvide_ExternalDockerConfigPathWithIdentityToken_ExpectedResults(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatalf("unexpected error when creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, "config.json")
	js := `{
		"auths": {
			"index.docker.io": {
				"auth": "MDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwOg==",
				"identitytoken": "OPAQUE_TOKEN"
			}
		}
	}`

	err = ioutil.WriteFile(fn, []byte(js), 0600)
	if err != nil {
		t.Fatalf("unexpected error when writing config file: %v", err)
	}

	defaultProvider := defaultAuthProvider{
		configPath: fn,
	}

	authConfig, err := defaultProvider.Provide(context.Background(), "index.docker.io/v1/test:v1")
	if err != nil {
		t.Fatalf("unexpected error in Provide: %v", err)
	}

	if authConfig.Username != "00000000-0000-0000-0000-000000000000" || authConfig.IdentityToken != "OPAQUE_TOKEN" {
		t.Fatalf("incorrect username %v or identitytoken %v returned", authConfig.Username, authConfig.IdentityToken)
	}
}
