// Copyright 2024 The PipeCD Authors.
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

package git

import (
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommit(t *testing.T) {
	log, err := os.ReadFile("testdata/log.txt")
	require.NoError(t, err)

	expected := []Commit{
		{
			Author:          "nghialv",
			Committer:       "kapetanios-robot",
			CreatedAt:       1565752022,
			Hash:            "74e20ede0242fdc7fd75b5be56e8d7fa72060707",
			AbbreviatedHash: "74e20ed",
			Message:         "wip",
		},
		{
			Author:          "Le Van Nghia",
			Committer:       "kapetanios-robot",
			CreatedAt:       1565749682,
			Hash:            "c9a7596e7e92ea5e3f03eeb951f632acb02b88a3",
			AbbreviatedHash: "c9a7596",
			Message:         `Add implementation of inplug service (#648)`,
			Body: `**What this PR does / why we need it**:

**Which issue(s) this PR fixes**:

Fixes #

**Does this PR introduce a user-facing change?**:
<!--
If no, just write "NONE" in the release-note block below.
-->
` + "```" + `release-note
NONE
` + "```" + `

This PR was merged by Kapetanios.`,
		},
		{
			Author:          "nghialv",
			Committer:       "kapetanios-robot",
			CreatedAt:       2565752022,
			Hash:            "24e20ede0242fdc7fd75b5be56e8d7fa72060707",
			AbbreviatedHash: "24e20ed",
			Message:         `Added commands to "kapectl" for creating, updating project secret (#475)`,
		},
	}
	commits, err := parseCommits(string(log))
	require.NoError(t, err)
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Hash > expected[j].Hash
	})
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Hash > commits[j].Hash
	})
	assert.Equal(t, expected, commits)
}

func TestGetTrailerValueByKey(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		body string
		key  string
		want string
	}{
		{
			name: "single line trailer",
			body: "Some commit message\nKey: value",
			key:  "Key",
			want: "value",
		},
		{
			name: "multiple line trailers",
			body: "Some commit message\nKey1: value1\nKey2: value2",
			key:  "Key2",
			want: "value2",
		},
		{
			name: "trailer with spaces",
			body: "Some commit message\nKey:    spaced value    ",
			key:  "Key",
			want: "spaced value",
		},
		{
			name: "no matching trailer",
			body: "Some commit message\nOtherKey: value",
			key:  "Key",
			want: "",
		},
		{
			name: "empty body",
			body: "",
			key:  "Key",
			want: "",
		},
		{
			name: "multiline value",
			body: "Some commit message\nKey: first line\n second line",
			key:  "Key",
			want: "first line",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := &Commit{Body: tc.body}
			got := c.GetTrailerValueByKey(tc.key)
			assert.Equal(t, tc.want, got)
		})
	}
}
