#!/bin/bash

set -eo pipefail
set -u
set +x

SHELL_FOLDER=$(
    cd $(dirname ${BASH_SOURCE[0]:-$0})
    pwd
)

ROOT_DIR=${SHELL_FOLDER}/../

genAwsRegcred() {
    kubectl delete secret --ignore-not-found -n default aws-regcred
    kubectl create secret docker-registry -n default aws-regcred \
        --docker-server=111111111111.dkr.ecr.ap-east-1.amazonaws.com \
        --docker-username=AWS \
        --docker-password=$(aws ecr get-login-password --region ap-east-1) \
        --namespace=default
}

loginDocker() {
    aws ecr get-login-password --region ap-east-1 | docker login --username AWS --password-stdin 111111111111.dkr.ecr.ap-east-1.amazonaws.com
}

recordTimestampFile() {
    local timestampFile=$1

    if [ ! -f "$timestampFile" ]; then
        touch "$timestampFile"
    fi

    local previousTimestamp=$(cat "${timestampFile}")
    local currentTimestamp=$(date +%s%N | cut -b1-13)
    local compareTimestamp=$(($(date +%s%N | cut -b1-13) - 7000000))

    if [ "${previousTimestamp}" = "" ] || [ "${compareTimestamp}" -gt "${previousTimestamp}" ]; then
        echo "${currentTimestamp}" >${timestampFile}
        echo 0
    else
        echo 1
    fi
}

checkGenAwsRegcred() {
    local result=$(recordTimestampFile ".gen-aws-secret-timestamp")
    if [ ${result} -eq 0 ]; then
        genAwsRegcred
    fi
}

checkLoginDocker() {
    local result=$(recordTimestampFile ".login-docker-timestamp")
    if [ ${result} -eq 0 ]; then
        loginDocker
    fi
}

build() {
    git fetch
    git checkout -f main
    git pull origin main

    go mod tidy
    go build -o app-go

    rm -rf ./dev_deploy/packages
    cp -rf /var/package/wallet-front/packages ./dev_deploy/packages

    docker build -f ./dev_deploy/Dockerfile . -t 111111111111.dkr.ecr.ap-east-1.amazonaws.com/wallet:latest
    docker push 111111111111.dkr.ecr.ap-east-1.amazonaws.com/wallet:latest
    echo "build successful ............ "
}

deploy() {
    kubectl rollout restart deployment wallet-deployment -n default

    ATTEMPTS=0
    ROLLOUT_STATUS_CMD="kubectl rollout status deployment wallet-deployment -n default"
    until $ROLLOUT_STATUS_CMD || [ $ATTEMPTS -eq 60 ]; do
        $ROLLOUT_STATUS_CMD
        ATTEMPTS=$((ATTEMPTS + 1))
        sleep 10
    done
}

bootstrap() {
    cd ${ROOT_DIR}

    checkGenAwsRegcred
    checkLoginDocker
    build
    deploy
}

# Start up
bootstrap
