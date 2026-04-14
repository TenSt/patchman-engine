package evaluator

import (
	"app/base/inventory"
	"app/base/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// systemPlatformV2ScanRow is a flat scan target for a joined inventory + patch load.
// Every selected column is aliased with an si_ or sp_ prefix so rh_account_id and
// other overlapping names are never ambiguous.
type systemPlatformV2ScanRow struct {
	SIID                               int64          `gorm:"column:si_id"`
	SIInventoryID                      string         `gorm:"column:si_inventory_id"`
	SIRhAccountID                      int            `gorm:"column:si_rh_account_id"`
	SIVmaasJSON                        *string        `gorm:"column:si_vmaas_json"`
	SIJSONChecksum                     *string        `gorm:"column:si_json_checksum"`
	SILastUpdated                      *time.Time     `gorm:"column:si_last_updated"`
	SIUnchangedSince                   *time.Time     `gorm:"column:si_unchanged_since"`
	SILastUpload                       *time.Time     `gorm:"column:si_last_upload"`
	SIStale                            bool           `gorm:"column:si_stale"`
	SIDisplayName                      string         `gorm:"column:si_display_name"`
	SIReporterID                       *int           `gorm:"column:si_reporter_id"`
	SIYumUpdates                       []byte         `gorm:"column:si_yum_updates"`
	SIYumChecksum                      *string        `gorm:"column:si_yum_checksum"`
	SISatelliteManaged                 bool           `gorm:"column:si_satellite_managed"`
	SIBuiltPkgcache                    bool           `gorm:"column:si_built_pkgcache"`
	SIArch                             *string        `gorm:"column:si_arch"`
	SIBootc                            bool           `gorm:"column:si_bootc"`
	SITags                             []byte         `gorm:"column:si_tags"`
	SICreated                          time.Time      `gorm:"column:si_created"`
	SIStaleTimestamp                   *time.Time     `gorm:"column:si_stale_timestamp"`
	SIStaleWarningTimestamp            *time.Time     `gorm:"column:si_stale_warning_timestamp"`
	SICulledTimestamp                  *time.Time     `gorm:"column:si_culled_timestamp"`
	SIOSName                           *string        `gorm:"column:si_os_name"`
	SIOSMajor                          *int16         `gorm:"column:si_os_major"`
	SIOSMinor                          *int16         `gorm:"column:si_os_minor"`
	SIRhsmVersion                      *string        `gorm:"column:si_rhsm_version"`
	SISubscriptionManagerID            *uuid.UUID     `gorm:"column:si_subscription_manager_id"`
	SISapWorkload                      bool           `gorm:"column:si_sap_workload"`
	SISapWorkloadSIDs                  pq.StringArray `gorm:"column:si_sap_workload_sids"`
	SIAnsibleWorkload                  bool           `gorm:"column:si_ansible_workload"`
	SIAnsibleWorkloadControllerVersion *string        `gorm:"column:si_ansible_workload_controller_version"`
	SIMssqlWorkload                    bool           `gorm:"column:si_mssql_workload"`
	SIMssqlWorkloadVersion             *string        `gorm:"column:si_mssql_workload_version"`
	SIWorkspaces                       []byte         `gorm:"column:si_workspaces"`

	SPSystemID                         int64      `gorm:"column:sp_system_id"`
	SPRhAccountID                      int        `gorm:"column:sp_rh_account_id"`
	SPLastEvaluation                   *time.Time `gorm:"column:sp_last_evaluation"`
	SPInstallableAdvisoryCountCache    int        `gorm:"column:sp_installable_advisory_count_cache"`
	SPInstallableAdvisoryEnhCountCache int        `gorm:"column:sp_installable_advisory_enh_count_cache"`
	SPInstallableAdvisoryBugCountCache int        `gorm:"column:sp_installable_advisory_bug_count_cache"`
	SPInstallableAdvisorySecCountCache int        `gorm:"column:sp_installable_advisory_sec_count_cache"`
	SPPackagesInstalled                int        `gorm:"column:sp_packages_installed"`
	SPPackagesInstallable              int        `gorm:"column:sp_packages_installable"`
	SPPackagesApplicable               int        `gorm:"column:sp_packages_applicable"`
	SPThirdParty                       bool       `gorm:"column:sp_third_party"`
	SPApplicableAdvisoryCountCache     int        `gorm:"column:sp_applicable_advisory_count_cache"`
	SPApplicableAdvisoryEnhCountCache  int        `gorm:"column:sp_applicable_advisory_enh_count_cache"`
	SPApplicableAdvisoryBugCountCache  int        `gorm:"column:sp_applicable_advisory_bug_count_cache"`
	SPApplicableAdvisorySecCountCache  int        `gorm:"column:sp_applicable_advisory_sec_count_cache"`
	SPTemplateID                       *int64     `gorm:"column:sp_template_id"`
}

const loadSystemPlatformV2SQL = `
SELECT
	si.id AS si_id,
	si.inventory_id::text AS si_inventory_id,
	si.rh_account_id AS si_rh_account_id,
	si.vmaas_json AS si_vmaas_json,
	si.json_checksum AS si_json_checksum,
	si.last_updated AS si_last_updated,
	si.unchanged_since AS si_unchanged_since,
	si.last_upload AS si_last_upload,
	si.stale AS si_stale,
	si.display_name AS si_display_name,
	si.reporter_id AS si_reporter_id,
	si.yum_updates AS si_yum_updates,
	si.yum_checksum AS si_yum_checksum,
	si.satellite_managed AS si_satellite_managed,
	si.built_pkgcache AS si_built_pkgcache,
	si.arch AS si_arch,
	si.bootc AS si_bootc,
	si.tags AS si_tags,
	si.created AS si_created,
	si.stale_timestamp AS si_stale_timestamp,
	si.stale_warning_timestamp AS si_stale_warning_timestamp,
	si.culled_timestamp AS si_culled_timestamp,
	si.os_name AS si_os_name,
	si.os_major AS si_os_major,
	si.os_minor AS si_os_minor,
	si.rhsm_version AS si_rhsm_version,
	si.subscription_manager_id AS si_subscription_manager_id,
	si.sap_workload AS si_sap_workload,
	si.sap_workload_sids AS si_sap_workload_sids,
	si.ansible_workload AS si_ansible_workload,
	si.ansible_workload_controller_version AS si_ansible_workload_controller_version,
	si.mssql_workload AS si_mssql_workload,
	si.mssql_workload_version AS si_mssql_workload_version,
	si.workspaces AS si_workspaces,
	sp.system_id AS sp_system_id,
	sp.rh_account_id AS sp_rh_account_id,
	sp.last_evaluation AS sp_last_evaluation,
	sp.installable_advisory_count_cache AS sp_installable_advisory_count_cache,
	sp.installable_advisory_enh_count_cache AS sp_installable_advisory_enh_count_cache,
	sp.installable_advisory_bug_count_cache AS sp_installable_advisory_bug_count_cache,
	sp.installable_advisory_sec_count_cache AS sp_installable_advisory_sec_count_cache,
	sp.packages_installed AS sp_packages_installed,
	sp.packages_installable AS sp_packages_installable,
	sp.packages_applicable AS sp_packages_applicable,
	sp.third_party AS sp_third_party,
	sp.applicable_advisory_count_cache AS sp_applicable_advisory_count_cache,
	sp.applicable_advisory_enh_count_cache AS sp_applicable_advisory_enh_count_cache,
	sp.applicable_advisory_bug_count_cache AS sp_applicable_advisory_bug_count_cache,
	sp.applicable_advisory_sec_count_cache AS sp_applicable_advisory_sec_count_cache,
	sp.template_id AS sp_template_id
FROM system_inventory AS si
JOIN system_patch AS sp
	ON sp.system_id = si.id AND sp.rh_account_id = si.rh_account_id
WHERE si.rh_account_id = ? AND si.inventory_id = ?::uuid
LIMIT 1`

// loadSystemPlatformV2 loads one system row from system_inventory joined to system_patch.
// It returns (nil, nil) when no matching row exists (no gorm.ErrRecordNotFound).
func loadSystemPlatformV2(db *gorm.DB, rhAccountID int, inventoryID string) (*models.SystemPlatformV2, error) {
	var row systemPlatformV2ScanRow
	if err := db.Raw(loadSystemPlatformV2SQL, rhAccountID, inventoryID).Scan(&row).Error; err != nil {
		return nil, err
	}
	// Match prior Model(...).Find behavior: no row yields an empty aggregate (ID 0), not an error.
	if row.SIID == 0 {
		return nil, nil
	}
	return scanRowToSystemPlatformV2(&row), nil
}

func scanRowToSystemPlatformV2(row *systemPlatformV2ScanRow) *models.SystemPlatformV2 {
	inv := models.SystemInventory{
		ID:                               row.SIID,
		InventoryID:                      row.SIInventoryID,
		RhAccountID:                      row.SIRhAccountID,
		VmaasJSON:                        row.SIVmaasJSON,
		JSONChecksum:                     row.SIJSONChecksum,
		LastUpdated:                      row.SILastUpdated,
		UnchangedSince:                   row.SIUnchangedSince,
		LastUpload:                       row.SILastUpload,
		Stale:                            row.SIStale,
		DisplayName:                      row.SIDisplayName,
		ReporterID:                       row.SIReporterID,
		YumUpdates:                       row.SIYumUpdates,
		YumChecksum:                      row.SIYumChecksum,
		SatelliteManaged:                 row.SISatelliteManaged,
		BuiltPkgcache:                    row.SIBuiltPkgcache,
		Arch:                             row.SIArch,
		Bootc:                            row.SIBootc,
		Tags:                             row.SITags,
		Created:                          row.SICreated,
		StaleTimestamp:                   row.SIStaleTimestamp,
		StaleWarningTimestamp:            row.SIStaleWarningTimestamp,
		CulledTimestamp:                  row.SICulledTimestamp,
		OSName:                           row.SIOSName,
		OSMajor:                          row.SIOSMajor,
		OSMinor:                          row.SIOSMinor,
		RhsmVersion:                      row.SIRhsmVersion,
		SubscriptionManagerID:            row.SISubscriptionManagerID,
		SapWorkload:                      row.SISapWorkload,
		SapWorkloadSIDs:                  row.SISapWorkloadSIDs,
		AnsibleWorkload:                  row.SIAnsibleWorkload,
		AnsibleWorkloadControllerVersion: row.SIAnsibleWorkloadControllerVersion,
		MssqlWorkload:                    row.SIMssqlWorkload,
		MssqlWorkloadVersion:             row.SIMssqlWorkloadVersion,
	}
	if len(row.SIWorkspaces) > 0 {
		var ws inventory.Groups
		if err := json.Unmarshal(row.SIWorkspaces, &ws); err == nil {
			inv.Workspaces = &ws
		}
	}
	patch := models.SystemPatch{
		SystemID:                         row.SPSystemID,
		RhAccountID:                      row.SPRhAccountID,
		LastEvaluation:                   row.SPLastEvaluation,
		InstallableAdvisoryCountCache:    row.SPInstallableAdvisoryCountCache,
		InstallableAdvisoryEnhCountCache: row.SPInstallableAdvisoryEnhCountCache,
		InstallableAdvisoryBugCountCache: row.SPInstallableAdvisoryBugCountCache,
		InstallableAdvisorySecCountCache: row.SPInstallableAdvisorySecCountCache,
		PackagesInstalled:                row.SPPackagesInstalled,
		PackagesInstallable:              row.SPPackagesInstallable,
		PackagesApplicable:               row.SPPackagesApplicable,
		ThirdParty:                       row.SPThirdParty,
		ApplicableAdvisoryCountCache:     row.SPApplicableAdvisoryCountCache,
		ApplicableAdvisoryEnhCountCache:  row.SPApplicableAdvisoryEnhCountCache,
		ApplicableAdvisoryBugCountCache:  row.SPApplicableAdvisoryBugCountCache,
		ApplicableAdvisorySecCountCache:  row.SPApplicableAdvisorySecCountCache,
		TemplateID:                       row.SPTemplateID,
	}
	return &models.SystemPlatformV2{Inventory: inv, Patch: patch}
}

func systemPlatformV2ToSystemPlatform(v2 *models.SystemPlatformV2) *models.SystemPlatform {
	if v2 == nil {
		return &models.SystemPlatform{}
	}
	inv := v2.Inventory
	pat := v2.Patch
	return &models.SystemPlatform{
		ID:                               inv.ID,
		InventoryID:                      inv.InventoryID,
		RhAccountID:                      inv.RhAccountID,
		VmaasJSON:                        inv.VmaasJSON,
		JSONChecksum:                     inv.JSONChecksum,
		LastUpdated:                      inv.LastUpdated,
		UnchangedSince:                   inv.UnchangedSince,
		LastEvaluation:                   pat.LastEvaluation,
		InstallableAdvisoryCountCache:    pat.InstallableAdvisoryCountCache,
		InstallableAdvisoryEnhCountCache: pat.InstallableAdvisoryEnhCountCache,
		InstallableAdvisoryBugCountCache: pat.InstallableAdvisoryBugCountCache,
		InstallableAdvisorySecCountCache: pat.InstallableAdvisorySecCountCache,
		LastUpload:                       inv.LastUpload,
		StaleTimestamp:                   inv.StaleTimestamp,
		StaleWarningTimestamp:            inv.StaleWarningTimestamp,
		CulledTimestamp:                  inv.CulledTimestamp,
		Stale:                            inv.Stale,
		DisplayName:                      inv.DisplayName,
		PackagesInstalled:                pat.PackagesInstalled,
		PackagesInstallable:              pat.PackagesInstallable,
		PackagesApplicable:               pat.PackagesApplicable,
		ThirdParty:                       pat.ThirdParty,
		ReporterID:                       inv.ReporterID,
		TemplateID:                       pat.TemplateID,
		YumUpdates:                       inv.YumUpdates,
		YumChecksum:                      inv.YumChecksum,
		SatelliteManaged:                 inv.SatelliteManaged,
		BuiltPkgcache:                    inv.BuiltPkgcache,
		ApplicableAdvisoryCountCache:     pat.ApplicableAdvisoryCountCache,
		ApplicableAdvisoryEnhCountCache:  pat.ApplicableAdvisoryEnhCountCache,
		ApplicableAdvisoryBugCountCache:  pat.ApplicableAdvisoryBugCountCache,
		ApplicableAdvisorySecCountCache:  pat.ApplicableAdvisorySecCountCache,
		Arch:                             inv.Arch,
		Bootc:                            inv.Bootc,
	}
}
