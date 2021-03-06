// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/goharbor/harbor/src/pkg/retention/res"

	"github.com/stretchr/testify/suite"
)

// IndexTestSuite tests the rule index
type IndexTestSuite struct {
	suite.Suite
}

// TestIndexEntry is entry of IndexTestSuite
func TestIndexEntry(t *testing.T) {
	suite.Run(t, new(IndexTestSuite))
}

// SetupSuite ...
func (suite *IndexTestSuite) SetupSuite() {
	Register(&IndexMeta{
		TemplateID: "fakeEvaluator",
		Action:     "retain",
		Parameters: []*IndexedParam{
			{
				Name:     "fakeParam",
				Type:     "int",
				Unit:     "count",
				Required: true,
			},
		},
	}, newFakeEvaluator)
}

// TestRegister tests register
func (suite *IndexTestSuite) TestGet() {

	params := make(Parameters)
	params["fakeParam"] = 99
	evaluator, err := Get("fakeEvaluator", params)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), evaluator)

	candidates := []*res.Candidate{{
		Namespace:  "library",
		Repository: "harbor",
		Kind:       "image",
		Tag:        "latest",
		PushedTime: time.Now().Unix(),
		Labels:     []string{"L1", "L2"},
	}}

	res, err := evaluator.Process(candidates)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(res))
	assert.Condition(suite.T(), func() bool {
		c := res[0]
		return c.Repository == "harbor" && c.Tag == "latest"
	})
}

// TestIndex tests Index
func (suite *IndexTestSuite) TestIndex() {
	metas := Index()
	require.Equal(suite.T(), 1, len(metas))
	assert.Condition(suite.T(), func() bool {
		m := metas[0]
		return m.TemplateID == "fakeEvaluator" &&
			m.Action == "retain" &&
			len(m.Parameters) > 0
	})
}

type fakeEvaluator struct {
	i int
}

// Process rule
func (e *fakeEvaluator) Process(artifacts []*res.Candidate) ([]*res.Candidate, error) {
	return artifacts, nil
}

// Action of the rule
func (e *fakeEvaluator) Action() string {
	return "retain"
}

// newFakeEvaluator is the factory of fakeEvaluator
func newFakeEvaluator(parameters Parameters) Evaluator {
	i := 10
	if v, ok := parameters["fakeParam"]; ok {
		i = v.(int)
	}

	return &fakeEvaluator{i}
}
