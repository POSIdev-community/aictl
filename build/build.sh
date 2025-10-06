#!/bin/bash

go build -o ./aictl ../cmd/run/main.go && \
sudo cp aictl /usr/bin/aictl && \
aictl completion bash > /etc/bash_completion.d/aictl
