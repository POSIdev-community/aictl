#!/usr/bin/env bash
set -euo pipefail

URL="${1:?URL required}"
TOKEN="${2:?TOKEN required}"
PROJECT_NAME="${3:?PROJECT_NAME required}"

AICTL="${AICTL:-aictl}"
FIXTURES_DIR="${FIXTURES_DIR:?FIXTURES_DIR required}"
WORK_DIR="${WORK_DIR:?WORK_DIR required}"

SBOM_FILE="${FIXTURES_DIR}/sbom/sbom.json"
if [[ ! -f "${SBOM_FILE}" ]]; then
	echo "error: SBOM fixture missing: ${SBOM_FILE}" >&2
	echo "Place an SBOM JSON file at tests/e2e/fixtures/sbom/sbom.json before running the SBOM pipeline." >&2
	exit 1
fi

project_id=""
scan_id=""

CONN=( -u "${URL}" -t "${TOKEN}" --tls-skip )

cleanup() {
	local code=$?
	if [[ -n "${project_id}" ]]; then
		"${AICTL}" delete "${CONN[@]}" projects "${project_id}" || true
	fi
	exit "${code}"
}
trap cleanup EXIT

project_id=$("${AICTL}" create "${CONN[@]}" sbom-project "${PROJECT_NAME}" --file "${SBOM_FILE}" --safe -v)
scan_id=$("${AICTL}" scan "${CONN[@]}" sbom -p "${project_id}" -v)
"${AICTL}" scan "${CONN[@]}" await "${scan_id}" -p "${project_id}" -v

aie_version=$("${AICTL}" get "${CONN[@]}" version)
printf '{"project_id":"%s","scan_id":"%s","aie_version":"%s"}\n' \
	"${project_id}" "${scan_id}" "${aie_version}" >"${WORK_DIR}/meta.json"
