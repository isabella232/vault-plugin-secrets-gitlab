// Copyright  2021 Masahiro Yoshida
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

package gitlabtoken

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccRoleToken(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test (short)")
	}
	req, backend := newGitlabAccEnv(t)

	ID := envAsInt("GITLAB_PROJECT_ID", 1)

	t.Run("successfully create", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"id":     ID,
			"name":   "vault-role-test",
			"scopes": []string{"read_api"},
		}
		roleName := "successful"
		mustRoleCreate(t, backend, req.Storage, roleName, data)
		resp, err := testIssueRoleToken(t, backend, req, roleName, nil)
		require.NoError(t, err)
		require.False(t, resp.IsError())

		assert.NotEmpty(t, resp.Data["token"], "no token returned")
		assert.NotEmpty(t, resp.Data["id"], "no id returned")
		assert.NotEmpty(t, resp.Data["expires_at"], "default is 1d for expires_at")
	})

	t.Run("non-existing role", func(t *testing.T) {
		resp, err := testIssueRoleToken(t, backend, req, "non-existing", nil)
		require.NoError(t, err)
		require.True(t, resp.IsError())
	})

}

// create the token given role name
func testIssueRoleToken(t *testing.T, b logical.Backend, req *logical.Request, roleName string, data map[string]interface{}) (*logical.Response, error) {
	req.Operation = logical.CreateOperation
	req.Path = fmt.Sprintf("%s/%s", pathPatternToken, roleName)
	req.Data = data

	resp, err := b.HandleRequest(context.Background(), req)

	return resp, err
}