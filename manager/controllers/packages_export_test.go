package controllers

import (
	"app/base/core"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageExportJSON(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithParams("GET", "/", "", "", nil, "application/json", PackagesExportHandler, 3)

	var output []PackageItem
	CheckResponse(t, w, http.StatusOK, &output)
	assert.Equal(t, 4, len(output))
	byName := make(map[string]PackageItem, len(output))
	for _, p := range output {
		byName[p.Name] = p
	}
	k := byName["kernel"]
	assert.Equal(t, "The Linux kernel", k.Summary)
	assert.Equal(t, 3, k.SystemsInstalled)
	assert.Equal(t, 2, k.SystemsInstallable)
	assert.Equal(t, 2, k.SystemsApplicable)
}

func TestPackageExportCSV(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithParams("GET", "/", "", "", nil, "text/csv", PackagesExportHandler, 3)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	lines := strings.Split(body, "\r\n")

	assert.Equal(t, 6, len(lines))
	assert.Equal(t, "name,summary,systems_installed,systems_installable,systems_applicable", lines[0])

	// Export applies default sort by package name (same as list API).
	data := []string{lines[1], lines[2], lines[3], lines[4]}
	sort.Strings(data)
	assert.Equal(t, []string{
		"bash,The GNU Bourne Again shell,1,0,0",
		"curl,A utility for getting files from remote servers...,1,0,0",
		"firefox,Mozilla Firefox Web browser,2,2,2",
		"kernel,The Linux kernel,3,2,2",
	}, data)
}
