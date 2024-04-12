// Copyright The Enterprise Contract Contributors
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

package asciidoc

import (
	_ "embed"

	"github.com/enterprise-contract/ec-cli/internal/documentation/asciidoc/cli"
	"github.com/enterprise-contract/ec-cli/internal/documentation/asciidoc/rego"
)

func GenerateAsciidoc(module string) error {
	if err := cli.GenerateCommandLineDocumentation(module); err != nil {
		return err
	}

	if err := rego.GenerateRegoReference(module); err != nil {
		return err
	}

	return nil
}
