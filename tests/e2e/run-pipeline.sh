#!/usr/bin/env bash
set -euo pipefail

URL="${1:?URL required}"
TOKEN="${2:?TOKEN required}"
PROJECT_NAME="${3:?PROJECT_NAME required}"

AICTL="${AICTL:-aictl}"
FIXTURES_DIR="${FIXTURES_DIR:?FIXTURES_DIR required}"
WORK_DIR="${WORK_DIR:?WORK_DIR required}"
AIPROJ_FIXTURE="${AIPROJ_FIXTURE:?AIPROJ_FIXTURE required}"

project_id=""
branch_id=""
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

branch_name="default"

project_id=$("${AICTL}" create "${CONN[@]}" project "${PROJECT_NAME}" --safe -v)

aiproj_path="${WORK_DIR}/aiproj.json"
python3 -c "
import json
with open('${AIPROJ_FIXTURE}') as f:
    data = json.load(f)
data['ProjectName'] = '${PROJECT_NAME}'
with open('${aiproj_path}', 'w') as f:
    json.dump(data, f)
"

"${AICTL}" set "${CONN[@]}" project settings -p "${project_id}" -f "${aiproj_path}" -v
branch_id=$("${AICTL}" create "${CONN[@]}" branch "${branch_name}" -p "${project_id}" --safe -v)

"${AICTL}" update "${CONN[@]}" sources "${FIXTURES_DIR}/project" -p "${project_id}" -b "${branch_id}" --update-languages -v
scan_id=$("${AICTL}" scan "${CONN[@]}" branch "${branch_id}" -p "${project_id}" -v)

"${AICTL}" scan "${CONN[@]}" await "${scan_id}" -p "${project_id}" -v
"${AICTL}" get "${CONN[@]}" scan report sarif "${scan_id}" -p "${project_id}" \
	-o "${WORK_DIR}/sarif.json" --include-glossary --localization en -v

aie_version=$("${AICTL}" get "${CONN[@]}" version)
printf '{"project_id":"%s","branch_id":"%s","scan_id":"%s","aie_version":"%s"}\n' \
	"${project_id}" "${branch_id}" "${scan_id}" "${aie_version}" >"${WORK_DIR}/meta.json"
