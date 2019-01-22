/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Credentials represents the google-stackdriver-credentials helper application.
type Credentials struct {
	buildpack buildpack.Buildpack
	layer     layers.Layer
}

// Contributes makes the contribution to launch
func (c Credentials) Contribute() error {
	return c.layer.Contribute(marker{c.buildpack.Info}, func(layer layers.Layer) error {
		if err := os.RemoveAll(layer.Root); err != nil {
			return err
		}

		if err := helper.CopyFile(filepath.Join(c.buildpack.Root, "bin", "google-stackdriver-credentials"), filepath.Join(layer.Root, "bin", "google-stackdriver-credentials")); err != nil {
			return err
		}

		return layer.WriteProfile("google-stackdriver-credentials", `printf "Configuring Google Stackdriver Credentials\n"

google-stackdriver-credentials %[1]s
export GOOGLE_APPLICATION_CREDENTIALS=%[1]s
`, filepath.Join(layer.Root, "google-stackdriver-credentials.json"))
	}, layers.Launch)
}

// String makes Credentials satisfy the Stringer interface.
func (c Credentials) String() string {
	return fmt.Sprintf("Credentials{ buildpack: %s, layer: %s }", c.buildpack, c.layer)
}

// NewCredentials creates a new Credentials instance.
func NewCredentials(build build.Build) Credentials {
	return Credentials{build.Buildpack, build.Layers.Layer("google-stackdriver-credentials")}
}

type marker struct {
	buildpack.Info
}

func (m marker) Identity() (string, string) {
	return "Google Stackdriver Credentials", m.Version
}
