#!/bin/bash

url="https://application-inspector"
token="pPZi9DLtir/LayJFmdk9sa2r/YJwcb+5"

project_name="sbom-project"
sbom_path="./sbom.json"

report_path="./sarif.json"

clearAndExit() {
    aictl ctx clear -y
    exit $1
}

aictl ctx clear -y || clearAndExit $?
aictl ctx set -u $url -t $token --tls-skip || clearAndExit $?

project_id=$(aictl create sbom-project $project_name --file $sbom_path --safe -v) || clearAndExit $?
aictl ctx set -p $project_id || clearAndExit $?

scan_id=$(aictl scan sbom -v) || clearAndExit $?

aictl scan await $scan_id -v || clearAndExit $?
aictl get scan report sarif $scan_id -o $report_path --include-glossary --localization en -v || clearAndExit $?

clearAndExit 0
