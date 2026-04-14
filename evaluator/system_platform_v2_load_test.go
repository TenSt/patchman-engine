package evaluator

import (
	"app/base/core"
	"app/base/database"
	"app/base/models"
	"app/base/utils"
	"bytes"
	"testing"
	"time"

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

// loadSystemData tests: regression against system_platform view until the view is removed.
func TestLoadSystemData_MatchesSystemPlatformView(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	var fromView models.SystemPlatform
	err := database.DB.Model(&models.SystemPlatform{}).
		Where("rh_account_id = ?", testFixtureRhAccount3).
		Where("inventory_id = ?::uuid", testFixtureInvUUID012).
		First(&fromView).Error
	require.NoError(t, err)

	loaded, err := loadSystemData(testFixtureRhAccount3, testFixtureInvUUID012)
	require.NoError(t, err)
	require.NotNil(t, loaded)

	assertSystemPlatformViewEquivalent(t, &fromView, loaded)

	// Mapper path is what loadSystemData uses; same fields as view for this fixture.
	mapped := systemPlatformV2ToSystemPlatform(mustLoadV2(t, testFixtureRhAccount3, testFixtureInvUUID012))
	assertSystemPlatformViewEquivalent(t, &fromView, mapped)
}

func TestLoadSystemData_NotFound(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	loaded, err := loadSystemData(testFixtureRhAccount3, "00000000-0000-0000-0000-00000000c0de")
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, int64(0), loaded.ID)
	assert.Empty(t, loaded.InventoryID)
}

func TestLoadSystemData_WrongAccount(t *testing.T) {
	setupSystemPlatformV2LoadTests(t)

	// Same inventory UUID as fixture system 12, but that row belongs to rh_account_id 3 only.
	loaded, err := loadSystemData(1, testFixtureInvUUID012)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, int64(0), loaded.ID)
	assert.Empty(t, loaded.InventoryID)
}

func mustLoadV2(t *testing.T, accountID int, invUUID string) *models.SystemPlatformV2 {
	t.Helper()
	v2, err := loadSystemPlatformV2(database.DB, accountID, invUUID)
	require.NoError(t, err)
	require.NotNil(t, v2)
	return v2
}

func assertSystemPlatformViewEquivalent(t *testing.T, fromView, loaded *models.SystemPlatform) {
	t.Helper()
	assert.Equal(t, fromView.ID, loaded.ID, "ID")
	assert.Equal(t, fromView.InventoryID, loaded.InventoryID, "InventoryID")
	assert.Equal(t, fromView.RhAccountID, loaded.RhAccountID, "RhAccountID")
	assertPtrStringEqual(t, fromView.VmaasJSON, loaded.VmaasJSON, "VmaasJSON")
	assertPtrStringEqual(t, fromView.JSONChecksum, loaded.JSONChecksum, "JSONChecksum")
	assertPtrTimeEqual(t, fromView.LastUpdated, loaded.LastUpdated, "LastUpdated")
	assertPtrTimeEqual(t, fromView.UnchangedSince, loaded.UnchangedSince, "UnchangedSince")
	assertPtrTimeEqual(t, fromView.LastEvaluation, loaded.LastEvaluation, "LastEvaluation")
	assert.Equal(t, fromView.InstallableAdvisoryCountCache, loaded.InstallableAdvisoryCountCache)
	assert.Equal(t, fromView.InstallableAdvisoryEnhCountCache, loaded.InstallableAdvisoryEnhCountCache)
	assert.Equal(t, fromView.InstallableAdvisoryBugCountCache, loaded.InstallableAdvisoryBugCountCache)
	assert.Equal(t, fromView.InstallableAdvisorySecCountCache, loaded.InstallableAdvisorySecCountCache)
	assert.Equal(t, fromView.ApplicableAdvisoryCountCache, loaded.ApplicableAdvisoryCountCache)
	assert.Equal(t, fromView.ApplicableAdvisoryEnhCountCache, loaded.ApplicableAdvisoryEnhCountCache)
	assert.Equal(t, fromView.ApplicableAdvisoryBugCountCache, loaded.ApplicableAdvisoryBugCountCache)
	assert.Equal(t, fromView.ApplicableAdvisorySecCountCache, loaded.ApplicableAdvisorySecCountCache)
	assertPtrTimeEqual(t, fromView.LastUpload, loaded.LastUpload, "LastUpload")
	assertPtrTimeEqual(t, fromView.StaleTimestamp, loaded.StaleTimestamp, "StaleTimestamp")
	assertPtrTimeEqual(t, fromView.StaleWarningTimestamp, loaded.StaleWarningTimestamp, "StaleWarningTimestamp")
	assertPtrTimeEqual(t, fromView.CulledTimestamp, loaded.CulledTimestamp, "CulledTimestamp")
	assert.Equal(t, fromView.Stale, loaded.Stale)
	assert.Equal(t, fromView.DisplayName, loaded.DisplayName)
	assert.Equal(t, fromView.PackagesInstalled, loaded.PackagesInstalled)
	assert.Equal(t, fromView.PackagesInstallable, loaded.PackagesInstallable)
	assert.Equal(t, fromView.PackagesApplicable, loaded.PackagesApplicable)
	assert.Equal(t, fromView.ThirdParty, loaded.ThirdParty)
	assert.Equal(t, fromView.ReporterID, loaded.ReporterID)
	assertPtrInt64Equal(t, fromView.TemplateID, loaded.TemplateID, "TemplateID")
	assert.True(t, bytes.Equal(fromView.YumUpdates, loaded.YumUpdates), "YumUpdates")
	assertPtrStringEqual(t, fromView.YumChecksum, loaded.YumChecksum, "YumChecksum")
	assert.Equal(t, fromView.SatelliteManaged, loaded.SatelliteManaged)
	assert.Equal(t, fromView.BuiltPkgcache, loaded.BuiltPkgcache)
	assertPtrStringEqual(t, fromView.Arch, loaded.Arch, "Arch")
	assert.Equal(t, fromView.Bootc, loaded.Bootc)
}

func assertPtrStringEqual(t *testing.T, want, got *string, field string) {
	t.Helper()
	if want == nil {
		assert.Nil(t, got, field)
		return
	}
	require.NotNil(t, got, field)
	assert.Equal(t, *want, *got, field)
}

func assertPtrInt64Equal(t *testing.T, want, got *int64, field string) {
	t.Helper()
	if want == nil {
		assert.Nil(t, got, field)
		return
	}
	require.NotNil(t, got, field)
	assert.Equal(t, *want, *got, field)
}

func assertPtrTimeEqual(t *testing.T, want, got *time.Time, field string) {
	t.Helper()
	if want == nil {
		assert.Nil(t, got, field)
		return
	}
	require.NotNil(t, got, field)
	assert.True(t, want.UTC().Equal(got.UTC()), "%s: want %v got %v", field, want.UTC(), got.UTC())
}
