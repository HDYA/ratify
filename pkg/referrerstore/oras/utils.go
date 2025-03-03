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

package oras

import (
	"regexp"
	"strings"

	"github.com/deislabs/ratify/pkg/ocispecs"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	artifactspec "github.com/oras-project/artifacts-spec/specs-go/v1"
)

// Detect the loopback IP (127.0.0.1)
var reLoopback = regexp.MustCompile(regexp.QuoteMeta("127.0.0.1"))

// Detect the loopback IPV6 (::1)
var reipv6Loopback = regexp.MustCompile(regexp.QuoteMeta("::1"))

func isInsecureRegistry(registry string, config *OrasStoreConf) bool {
	if config.UseHttp {
		return true
	}
	if strings.HasPrefix(registry, "localhost:") {
		return true
	}

	if reLoopback.MatchString(registry) {
		return true
	}
	if reipv6Loopback.MatchString(registry) {
		return true
	}

	return false
}

func OciDescriptorToReferenceDescriptor(ociDescriptor oci.Descriptor) ocispecs.ReferenceDescriptor {
	return ocispecs.ReferenceDescriptor{
		Descriptor:   ociDescriptor,
		ArtifactType: ociDescriptor.ArtifactType,
	}
}

func ArtifactManifestToReferenceManifest(artifactManifest artifactspec.Manifest) ocispecs.ReferenceManifest {
	var blobs []oci.Descriptor
	for _, blob := range artifactManifest.Blobs {
		ociBlob := oci.Descriptor{
			MediaType:   blob.MediaType,
			Digest:      blob.Digest,
			Size:        blob.Size,
			URLs:        blob.URLs,
			Annotations: blob.Annotations,
		}
		blobs = append(blobs, ociBlob)
	}

	subjects := []oci.Descriptor{{
		MediaType:   artifactManifest.Subject.MediaType,
		Digest:      artifactManifest.Subject.Digest,
		Size:        artifactManifest.Subject.Size,
		URLs:        artifactManifest.Subject.URLs,
		Annotations: artifactManifest.Subject.Annotations,
	}}

	return ocispecs.ReferenceManifest{
		MediaType:    ocispecs.MediaTypeArtifactManifest,
		ArtifactType: artifactManifest.ArtifactType,
		Blobs:        blobs,
		Subjects:     subjects,
	}
}
