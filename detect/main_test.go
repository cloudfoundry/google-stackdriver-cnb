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

package main

import (
	"testing"

	"github.com/buildpacks/libbuildpack/v2/buildplan"
	"github.com/cloudfoundry/google-stackdriver-cnb/java"
	"github.com/cloudfoundry/libcfbuildpack/v2/detect"
	"github.com/cloudfoundry/libcfbuildpack/v2/services"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("fails without service", func() {
			g.Expect(d(f.Detect)).To(gomega.Equal(detect.FailStatusCode))
		})

		it("passes with debugger service", func() {
			f.AddService("google-stackdriver-debugger", services.Credentials{"PrivateKeyData": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: java.DebuggerDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: java.DebuggerDependency},
				},
			}))
		})

		it("passes with profiler service", func() {
			f.AddService("google-stackdriver-profiler", services.Credentials{"PrivateKeyData": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: java.ProfilerDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: java.ProfilerDependency},
				},
			}))
		})

		it("passes with debugger and profiler services", func() {
			f.AddService("google-stackdriver-debugger", services.Credentials{"PrivateKeyData": "test-value"})
			f.AddService("google-stackdriver-profiler", services.Credentials{"PrivateKeyData": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: java.DebuggerDependency},
					{Name: java.ProfilerDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: java.DebuggerDependency},
					{Name: java.ProfilerDependency},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
