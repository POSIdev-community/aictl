package project

import "errors"

const ErrSbomUnsupported = "SBOM projects are supported starting from AIE 6.3"

var (
	ErrNameUsedBySource           = errors.New("name already used by a source project")
	ErrNameUsedBySbom             = errors.New("name already used by an SBOM project")
	ErrProjectNameUsedBySource    = errors.New("project name already used by a source project")
	ErrProjectNameUsedBySbom      = errors.New("project name already used by an SBOM project")
	ErrCannotUpdateSourcesOnSbom  = errors.New("cannot update sources on an SBOM project")
	ErrCannotUpdateSbomOnSource   = errors.New("cannot update SBOM on a source project")
	ErrCannotRunSbomScanOnSource  = errors.New("cannot run an SBOM scan on a source project")
	ErrCannotRunProjectScanOnSbom = errors.New("cannot run a project scan on an SBOM project; use 'aictl scan sbom'")
	ErrCannotCreateBranchOnSbom   = errors.New("cannot create a branch on an SBOM project")
	ErrCannotGetBranchesOnSbom    = errors.New("cannot get branches on an SBOM project")
	ErrCannotGetBranchOnSbom      = errors.New("cannot get a branch on an SBOM project")
	ErrCannotGetSetSettingsOnSbom = errors.New("cannot get/set project settings on an SBOM project")
	ErrCannotStartSbomScan        = errors.New("cannot start SBOM scan for this project")
)
