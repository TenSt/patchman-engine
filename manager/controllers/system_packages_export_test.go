package controllers

import (
	"app/base/core"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemPackagesExportHandlerJSON(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithParams("GET", "/:inventory_id/packages", "00000000-0000-0000-0000-000000000013", "",
		nil, "application/json", SystemPackagesExportHandler, 3)

	var output []SystemPackageInline
	CheckResponse(t, w, http.StatusOK, &output)
	assert.Equal(t, 4, len(output))
	byName := make(map[string]SystemPackageInline, len(output))
	for _, p := range output {
		byName[p.Name] = p
	}
	kernel := byName["kernel"]
	assert.Equal(t, "5.6.13-200.fc31.x86_64", kernel.EVRA)
	assert.Equal(t, "5.6.13-200.fc31.x86_64", kernel.LatestInstallable)
	assert.Equal(t, "5.6.13-200.fc31.x86_64", kernel.LatestApplicable)
	assert.Equal(t, "The Linux kernel", kernel.Summary)
}

func TestSystemPackagesExportHandlerCSV(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithParams("GET", "/:inventory_id/packages", "00000000-0000-0000-0000-000000000013", "",
		nil, "text/csv", SystemPackagesExportHandler, 3)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	lines := strings.Split(body, "\r\n")

	assert.Equal(t, 6, len(lines))
	assert.Equal(t, "name,evra,summary,description,updatable,update_status,latest_installable,latest_applicable", lines[0])

	// Export applies the same default sort as the list API (package name ascending).
	// nolint:lll
	assert.Equal(t, "bash,4.4.19-8.el8_0.x86_64,The GNU Bourne Again shell,The GNU Bourne Again shell (Bash) is a shell...,false,"+
		"None,4.4.19-8.el8_0.x86_64,4.4.19-8.el8_0.x86_64", lines[1])
	// nolint:lll
	assert.Equal(t, "curl,7.61.1-8.el8.x86_64,A utility for getting files from remote servers...,curl is a command line tool for transferring data...,false,"+
		"None,7.61.1-8.el8.x86_64,7.61.1-8.el8.x86_64", lines[2])
	assert.Equal(t, "firefox,76.0.1-1.fc31.x86_64,Mozilla Firefox Web browser,Mozilla Firefox is an "+
		"open-source web browser...,true,Installable,76.0.1-2.fc31.x86_64,77.0.1-1.fc31.x86_64", lines[3])
	assert.Equal(t, "kernel,5.6.13-200.fc31.x86_64,The Linux kernel,The kernel meta package,false,"+
		"None,5.6.13-200.fc31.x86_64,5.6.13-200.fc31.x86_64", lines[4])
}

func TestSystemPackagesExportUnknown(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithParams("GET", "/:inventory_id/packages", "unknownsystem", "", nil, "text/csv",
		SystemPackagesExportHandler, 3)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
