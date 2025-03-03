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
	"testing"

	core "k8s.io/api/core/v1"
)

// Checks K8 Docker Json Config Secret is properly extracted and
// credentials returned when Provide is called
func TestProvide_K8SecretDockerConfigJson_ReturnsExpected(t *testing.T) {
	var testSecret core.Secret
	js := `{
		"auths": {
			"index.docker.io": {
				"auth": "am9lam9lOmhlbGxv"
			}
		}
	}`
	testSecret.Data = make(map[string][]byte)
	testSecret.Data[core.DockerConfigJsonKey] = []byte(js)
	testSecret.Type = core.SecretTypeDockerConfigJson

	var k8secretprovider k8SecretAuthProvider

	authConfig, err := k8secretprovider.resolveCredentialFromSecret("index.docker.io", &testSecret)
	if err != nil {
		t.Fatalf("resolveCredentialFromSecret failed to get credential with err %v", err)
	}

	if authConfig.Username != "joejoe" || authConfig.Password != "hello" {
		t.Fatalf("resolveCredentialFromSecret returned incorrect credentials (username: %s, password: %s)", authConfig.Username, authConfig.Password)
	}
}

func TestProvide_K8SecretDockerConfigJsonWithIdentityToken_ReturnsExpected(t *testing.T) {
	var testSecret core.Secret
	js := `{
		"auths": {
			"index.docker.io": {
				"auth": "MDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwOg==",
				"identitytoken": "OPAQUE_TOKEN"
			}
		}
	}`
	testSecret.Data = make(map[string][]byte)
	testSecret.Data[core.DockerConfigJsonKey] = []byte(js)
	testSecret.Type = core.SecretTypeDockerConfigJson

	var k8secretprovider k8SecretAuthProvider

	authConfig, err := k8secretprovider.resolveCredentialFromSecret("index.docker.io", &testSecret)
	if err != nil {
		t.Fatalf("resolveCredentialFromSecret failed to get credential with err %v", err)
	}

	if authConfig.Username != "00000000-0000-0000-0000-000000000000" || authConfig.IdentityToken != "OPAQUE_TOKEN" {
		t.Fatalf("resolveCredentialFromSecret returned incorrect credentials (username: %s, identitytoken: %s)", authConfig.Username, authConfig.IdentityToken)
	}
}

// Checks K8 DockerCfg Secret is properly extracted and
// credentials returned when Provide is called
func TestProvide_K8SecretDockerCfg_ReturnsExpected(t *testing.T) {
	var testSecret core.Secret
	js := `{
		"index.docker.io": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}`
	testSecret.Data = make(map[string][]byte)
	testSecret.Data[core.DockerConfigKey] = []byte(js)
	testSecret.Type = core.SecretTypeDockercfg

	var k8secretprovider k8SecretAuthProvider

	authConfig, err := k8secretprovider.resolveCredentialFromSecret("index.docker.io", &testSecret)
	if err != nil {
		t.Fatalf("resolveCredentialFromSecret failed to get credential with err %v", err)
	}

	if authConfig.Username != "joejoe" || authConfig.Password != "hello" {
		t.Fatalf("resolveCredentialFromSecret returned incorrect credentials (username: %s, password: %s)", authConfig.Username, authConfig.Password)
	}
}

func TestProvide_K8SecretDockerCfgWithIdentityToken_ReturnsExpected(t *testing.T) {
	var testSecret core.Secret
	js := `{
		"index.docker.io": {
			"auth": "MDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwOg==",
			"identitytoken": "OPAQUE_TOKEN"
		}
	}`
	testSecret.Data = make(map[string][]byte)
	testSecret.Data[core.DockerConfigKey] = []byte(js)
	testSecret.Type = core.SecretTypeDockercfg

	var k8secretprovider k8SecretAuthProvider

	authConfig, err := k8secretprovider.resolveCredentialFromSecret("index.docker.io", &testSecret)
	if err != nil {
		t.Fatalf("resolveCredentialFromSecret failed to get credential with err %v", err)
	}

	if authConfig.Username != "00000000-0000-0000-0000-000000000000" || authConfig.IdentityToken != "OPAQUE_TOKEN" {
		t.Fatalf("resolveCredentialFromSecret returned incorrect credentials (username: %s, identitytoken: %s)", authConfig.Username, authConfig.IdentityToken)
	}
}

// Checks an error is returned for non-existent registry credential
func TestProvide_K8SecretNonExistentRegistry_ReturnsExpected(t *testing.T) {
	var testSecret core.Secret
	js := `{
		"index.docker.io": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}`
	testSecret.Data = make(map[string][]byte)
	testSecret.Data[core.DockerConfigKey] = []byte(js)
	testSecret.Type = core.SecretTypeDockercfg

	var k8secretprovider k8SecretAuthProvider

	_, err := k8secretprovider.resolveCredentialFromSecret("nonexistent.ghcr.io", &testSecret)
	if err != ErrorNoMatchingCredential {
		t.Fatalf("resolveCredentialFromSecret should have failed to get credential with err %v but returned err %v", ErrorNoMatchingCredential, err)
	}
}
