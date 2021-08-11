/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, BlueprintVersion 2.0 (the "License");
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

package laraboot

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/cloudfoundry/libcfbuildpack/helper"
)

// ProcessTemplateToFile writes out a specific template to the given file name
func ProcessTemplateToFile(templateBody string, outputPath string, data interface{}) error {
	template, err := template.New(filepath.Base(outputPath)).Parse(templateBody)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = template.Execute(&b, data)
	if err != nil {
		return err
	}

	return helper.WriteFileFromReader(outputPath, 0644, &b)
}
