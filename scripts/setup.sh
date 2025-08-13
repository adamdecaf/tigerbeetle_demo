#!/bin/bash

if [ ! -f data/1_0.tigerbeetle ]; then
    docker run --security-opt seccomp=unconfined -v $(pwd)/data:/data ghcr.io/tigerbeetle/tigerbeetle format --cluster=1 --replica=0 --replica-count=3 /data/1_0.tigerbeetle;
fi

if [ ! -f data/1_1.tigerbeetle ]; then
    docker run --security-opt seccomp=unconfined -v $(pwd)/data:/data ghcr.io/tigerbeetle/tigerbeetle format --cluster=1 --replica=1 --replica-count=3 /data/1_1.tigerbeetle;
fi

if [ ! -f data/1_2.tigerbeetle ]; then
    docker run --security-opt seccomp=unconfined -v $(pwd)/data:/data ghcr.io/tigerbeetle/tigerbeetle format --cluster=1 --replica=2 --replica-count=3 /data/1_2.tigerbeetle;
fi
