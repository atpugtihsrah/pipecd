// Copyright 2025 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package toolregistry

var terraformInstallScript = `
cd {{ .TmpDir }}
curl https://releases.hashicorp.com/terraform/{{ .Version }}/terraform_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip -o terraform_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip
unzip terraform_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip
mv terraform {{ .OutPath }}
`
