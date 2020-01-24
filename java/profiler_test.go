/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/google-stackdriver-cnb/java"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestProfiler(t *testing.T) {
	spec.Run(t, "Profiler", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns true if build plan does exist", func() {
			f.AddPlan(buildpackplan.Plan{Name: java.ProfilerDependency})
			f.AddDependency(java.ProfilerDependency, filepath.Join("testdata", "stub-profiler.tar.gz"))

			_, ok, err := java.NewProfiler(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
		})

		it("returns false if build plan does not exist", func() {
			_, ok, err := java.NewProfiler(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("contributes agent", func() {
			f.AddPlan(buildpackplan.Plan{Name: java.ProfilerDependency})
			f.AddDependency(java.ProfilerDependency, filepath.Join("testdata", "stub-profiler.tar.gz"))

			d, ok, err := java.NewProfiler(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())

			g.Expect(d.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("google-stackdriver-profiler-java")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "profiler_java_agent.so")).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HaveProfile("google-stackdriver-profiler", `if [[ -z "${BPL_GOOGLE_STACKDRIVER_MODULE+x}" ]]; then
    MODULE="default-module"
else
	MODULE=${BPL_GOOGLE_STACKDRIVER_MODULE}
fi

if [[ -z "${BPL_GOOGLE_STACKDRIVER_VERSION+x}" ]]; then
	VERSION=""
else
	VERSION=${BPL_GOOGLE_STACKDRIVER_VERSION}
fi

printf "Google Stackdriver Profiler enabled for ${MODULE}"

if [[ "${VERSION}" != "" ]]; then
	printf ":${VERSION}\n"
else
	printf "\n"
fi

AGENT="-agentpath:%s=--logtostderr=1,-cprof_service=${MODULE}"

if [[ "${VERSION}" != "" ]]; then
    AGENT="${AGENT},-cprof_service_version=${VERSION}"
fi

export JAVA_OPTS="${JAVA_OPTS} ${AGENT}"

`, filepath.Join(layer.Root, "profiler_java_agent.so")))
		})
	}, spec.Report(report.Terminal{}))
}
