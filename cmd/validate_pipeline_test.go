// Copyright 2022 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

//go:build unit

package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	hd "github.com/MakeNowJust/heredoc"
	"github.com/open-policy-agent/conftest/output"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	output2 "github.com/hacbs-contract/ec-cli/internal/output"
	"github.com/hacbs-contract/ec-cli/internal/policy/source"
)

func TestValidatePipelineCommandOutput(t *testing.T) {
	validate := func(_ context.Context, _ afero.Fs, fpath string, _ []source.PolicySource, namespace string) (*output2.Output, error) {
		return &output2.Output{
			PolicyCheck: []output.CheckResult{
				{
					FileName:  fpath,
					Namespace: namespace,
				},
			},
		}, nil
	}

	cmd := validatePipelineCmd(validate)

	var out bytes.Buffer
	cmd.SetOut(&out)

	cmd.SetArgs([]string{
		"--pipeline-file",
		"/path/file1.yaml",
		"--pipeline-file",
		"/path/file2.yaml",
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	assert.JSONEq(t, `[
		{
		  "filename": "/path/file1.yaml",
		  "namespace": "pipeline.main",
		  "success": true,
		  "violations": [],
		  "warnings": []
		},
		{
		  "filename": "/path/file2.yaml",
		  "namespace": "pipeline.main",
		  "success": true,
		  "violations": [],
		  "warnings": []
		}
	  ]`, out.String())
}

func TestValidatePipelinePolicySources(t *testing.T) {
	expected := []source.PolicySource{
		&source.PolicyUrl{Url: "spam-policy-source", Kind: source.PolicyKind},
		&source.PolicyUrl{Url: "ham-policy-source", Kind: source.PolicyKind},
		&source.PolicyUrl{Url: "bacon-data-source", Kind: source.DataKind},
		&source.PolicyUrl{Url: "eggs-data-source", Kind: source.DataKind},
	}
	validate := func(_ context.Context, _ afero.Fs, fpath string, sources []source.PolicySource, namespace string) (*output2.Output, error) {
		assert.Equal(t, expected, sources)
		return &output2.Output{}, nil
	}

	cmd := validatePipelineCmd(validate)

	var out bytes.Buffer
	cmd.SetOut(&out)

	cmd.SetArgs([]string{
		"--pipeline-file",
		"/path/file1.yaml",
		"--policy",
		"spam-policy-source",
		"--policy",
		"ham-policy-source",
		"--data",
		"bacon-data-source",
		"--data",
		"eggs-data-source",
	})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestOutputFormats(t *testing.T) {
	testJSONText := (`[{"filename":"/path/file1.yaml","namespace":"pipeline.main",` +
		`"violations":[],"warnings":[],"success":true}]`)

	testYAMLTest := hd.Doc(`
	- filename: /path/file1.yaml
	  namespace: pipeline.main
	  success: true
	  violations: []
	  warnings: []
	`)

	cases := []struct {
		name           string
		output         []string
		expectedFiles  map[string]string
		expectedStdout string
	}{
		{
			name:           "default output",
			expectedStdout: testJSONText,
		},
		{
			name:           "json stdout",
			output:         []string{"--output", "json"},
			expectedStdout: testJSONText,
		},
		{
			name:           "yaml stdout",
			output:         []string{"--output", "yaml"},
			expectedStdout: testYAMLTest,
		},
		{
			name:           "json and yaml to file",
			output:         []string{"--output", "json=out.json", "--output", "yaml=out.yaml"},
			expectedStdout: "",
			expectedFiles: map[string]string{
				"out.json": testJSONText,
				"out.yaml": testYAMLTest,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			validate := func(_ context.Context, _ afero.Fs, fpath string, sources []source.PolicySource, namespace string) (*output2.Output, error) {
				return &output2.Output{
					PolicyCheck: []output.CheckResult{
						{
							FileName:  fpath,
							Namespace: namespace,
						},
					},
				}, nil
			}

			cmd := validatePipelineCmd(validate)

			var out bytes.Buffer
			cmd.SetOut(&out)

			cmd.SetArgs(append([]string{
				"--pipeline-file",
				"/path/file1.yaml",
			}, c.output...))

			cmd.SetContext(withFs(context.Background(), fs))

			err := cmd.Execute()
			assert.NoError(t, err)
			assert.Equal(t, c.expectedStdout, out.String())

			for name, expectedText := range c.expectedFiles {
				actualText, err := afero.ReadFile(fs, name)
				assert.NoError(t, err)
				assert.Equal(t, expectedText, string(actualText))
			}
		})
	}
}

func TestValidatePipelineCommandErrors(t *testing.T) {
	validate := func(_ context.Context, _ afero.Fs, fpath string, _ []source.PolicySource, _ string) (*output2.Output, error) {
		return nil, errors.New(fpath)
	}

	cmd := validatePipelineCmd(validate)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SilenceUsage = true

	cmd.SetArgs([]string{
		"--pipeline-file",
		"/path/file1.yaml",
		"--pipeline-file",
		"/path/file2.yaml",
	})

	err := cmd.Execute()
	assert.Error(t, err, "2 errors occurred:\n\t* /path/file1.yaml\n\t* /path/file2.yaml\n")
	assert.Equal(t, "", out.String())
}
