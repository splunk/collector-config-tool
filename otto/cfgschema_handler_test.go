// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertSubMapsToStringMaps(t *testing.T) {
	m := map[string]any{
		"a": map[any]any{
			"b": 1,
			"c": map[any]any{
				"d": 1,
			},
		},
		"e": []any{
			map[any]any{
				"f": 1,
			},
		},
	}
	actual, err := convertAnyMapsToStringMaps(m)
	require.NoError(t, err)
	assert.EqualValues(t, m, actual)
	a, ok := m["a"].(map[string]any)
	assert.True(t, ok)
	_, ok = a["c"].(map[string]any)
	assert.True(t, ok)
	e, ok := m["e"].([]any)
	assert.True(t, ok)
	assert.NotNil(t, e)
	v := e[0]
	_, ok = v.(map[string]any)
	assert.True(t, ok)
}
