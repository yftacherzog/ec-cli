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

package something

import (
	"errors"

	ece "github.com/hacbs-contract/ec-cli/pkg/error"
)

var (
	TE001 = ece.NewError("TE001", "Test error", ece.ErrorExitStatus)
)

func hasImportRaw() error {
	return errors.New("an error") // want `don't return raw error, return pkg/error.Error`
}

func hasImportECE() error {
	return TE001
}