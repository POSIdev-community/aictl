#!/bin/bash

url="https://application-inspector"
token="pPZi9DLtir/LayJFmdk9sa2r/YJwcb+5"

project_name="project"
aiproj_path="./aiproj.json"

branch_name="default"
sources_path="./project"

report_path="./sarif.json"

clearAndExit() {
    aictl ctx clear -y
    exit $1
}

aictl ctx clear -y || clearAndExit $?
aictl ctx set -u $url -t $token --tls-skip || clearAndExit $?

project_id=$(aictl create project $project_name --safe -v) || clearAndExit $?
aictl ctx set -p $project_id || clearAndExit $?

aictl set project settings -f $aiproj_path -v || clearAndExit $?
branch_id=$(aictl create branch $branch_name --safe -v) || clearAndExit $?
aictl ctx set -b $branch_id || clearAndExit $?

aictl update sources $sources_path -v || clearAndExit $?
scan_id=$(aictl scan start branch $branch_id -v) || clearAndExit $?

aictl scan await $scan_id -v || clearAndExit $?
aictl get scan report sarif $scan_id -o $report_path --include-glossary --localization en -v || clearAndExit $?

clearAndExit 0
