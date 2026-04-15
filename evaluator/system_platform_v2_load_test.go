package evaluator

import (
	"app/base/core"
	"app/base/database"
	"app/base/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fixture rows from dev/test_data.sql: system_inventory id 12 + matching system_patch for account 3.
const (
	testFixtureInvUUID012 = "00000000-0000-0000-0000-000000000012"
	testFixtureRhAccount3 = 3
	testFixtureSystemID12 = int64(12)
)

func setupSystemPlatformV2LoadTests(t *testing.T) {
	t.Helper()
	utils.TestLoadEnv("conf/evaluator_common.env", "conf/evaluator_upload.env")
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
}

func TestLoadSystemPlatformV2_Found(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	v2, err := loadSystemPlatformV2(database.DB, testFixtureRhAccount3, testFixtureInvUUID012)
	require.NoError(t, err)
	require.NotNil(t, v2)

	assert.Equal(t, int64(12), v2.InternalSystemID())
	assert.Equal(t, testFixtureInvUUID012, v2.GetInventoryID())
	assert.Equal(t, v2.Inventory.ID, v2.Patch.SystemID)
	assert.Equal(t, v2.Inventory.RhAccountID, v2.Patch.RhAccountID)
	assert.Equal(t, testFixtureSystemID12, v2.Inventory.ID)
	assert.Equal(t, 2, v2.Patch.PackagesInstalled)
	assert.Equal(t, 2, v2.Patch.PackagesInstallable)
	assert.Equal(t, 2, v2.Patch.PackagesApplicable)
}

func TestLoadSystemPlatformV2_NotFound(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	v2, err := loadSystemPlatformV2(database.DB, testFixtureRhAccount3, "00000000-0000-0000-0000-00000000dead")
	require.NoError(t, err)
	assert.Nil(t, v2)
}

func TestLoadSystemPlatformV2_WrongAccount(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	// UUID ...012 exists only for rh_account_id 3 in test_data.sql.
	v2, err := loadSystemPlatformV2(database.DB, 1, testFixtureInvUUID012)
	require.NoError(t, err)
	assert.Nil(t, v2)
}

func TestLoadSystemData_MatchesLoadV2(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	loaded, err := loadSystemData(testFixtureRhAccount3, testFixtureInvUUID012)
	require.NoError(t, err)
	require.NotNil(t, loaded)

	direct, err := loadSystemPlatformV2(database.DB, testFixtureRhAccount3, testFixtureInvUUID012)
	require.NoError(t, err)
	require.NotNil(t, direct)

	assert.Equal(t, direct.Inventory, loaded.Inventory)
	assert.Equal(t, direct.Patch, loaded.Patch)
}

func TestLoadSystemData_NotFound(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	loaded, err := loadSystemData(testFixtureRhAccount3, "00000000-0000-0000-0000-00000000c0de")
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestLoadSystemData_WrongAccount(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	// Same inventory UUID as fixture system 12, but that row belongs to rh_account_id 3 only.
	loaded, err := loadSystemData(1, testFixtureInvUUID012)
	require.NoError(t, err)
	assert.Nil(t, loaded)
}
