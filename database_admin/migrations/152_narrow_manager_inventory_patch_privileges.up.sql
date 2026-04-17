-- Restore least-privilege grants for manager after migration 145 temporarily
-- broadened them for system_platform view/trigger updates (removed in 151).
REVOKE UPDATE ON system_inventory FROM manager;
GRANT UPDATE (stale) ON system_inventory TO manager;

REVOKE UPDATE ON system_patch FROM manager;
GRANT UPDATE (
    installable_advisory_count_cache,
    installable_advisory_enh_count_cache,
    installable_advisory_bug_count_cache,
    installable_advisory_sec_count_cache,
    applicable_advisory_count_cache,
    applicable_advisory_enh_count_cache,
    applicable_advisory_bug_count_cache,
    applicable_advisory_sec_count_cache,
    template_id) ON system_patch TO manager;
