#!/bin/bash

url="https://application-inspector"
token="pPZi9DLtir/LayJFmdk9sa2r/YJwcb+5"

project_name="project"
aiproj_path="./aiproj.json"

branch_name="default"
sources_path="./project"

retryTime=${retryTime:-3600}
retryInterval=5

report_path="./sarif.json"

clearAndExit() {
    aictl ctx clear -y
    exit $1
}

aictl ctx clear -y || clearAndExit $?
aictl ctx set -u $url -t $token --tls-skip || clearAndExit $?

project_id=$(aictl create project $project_name --safe -v) || clearAndExit $?
aictl ctx set -p $project_id || clearAndExit $?

branch_id=$(aictl create branch $branch_name --safe -v) || clearAndExit $?
aictl ctx set -b $branch_id || clearAndExit $?

scan_id=""
start_time=$(date +%s)
deadline=$((start_time + retryTime))

while true; do
    echo "⏳ Attempting project setup and scan start..."

    if aictl set project settings -f "$aiproj_path" -v &&
       aictl update sources "$sources_path" -v &&
       scan_id=$(aictl scan start branch "$branch_id" -v); then

        echo "✅ Project settings, sources, and scan start succeeded."
        break
    fi

    now=$(date +%s)
    if (( now >= deadline )); then
        echo "❌ Failed to complete setup + scan start within $retryTime seconds."
        clearAndExit 1
    fi

    echo "⏸️ Retrying in $retryInterval seconds..."
    sleep "$retryInterval"
done

aictl scan await $scan_id -v || clearAndExit $?
aictl get scan report sarif $scan_id -o $report_path --include-glossary --localization en -v || clearAndExit $?

clearAndExit 0
