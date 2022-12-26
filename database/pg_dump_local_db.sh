#!/usr/bin/env bash

function wd() {
    UNAME=$(uname)
    if [ "$UNAME" == "Linux" ] || [ "$UNAME" == "Darwin" ] ; then
        echo ${PWD}
    elif [[ "$UNAME" == CYGWIN* ]] ; then
        echo "/$(pwd)"
    elif [[ "$UNAME" == MINGW* ]] ; then
        echo "$(pwd)"
    else
        echo -e "[${red}${bold}ERROR${reset}] Failed to resolve the current working directory"
        exit 1
    fi
}

CWD=$(wd)
FILENAME=$1

function source_script_config() {
    source ${CWD}/utils/.colors
}

source_script_config

if [ -z ${DOCKER_CONTAINER_NAME} ] || [ -z ${SCHEMA_NAME} ] || [ -z ${FILENAME} ]; then
    echo -e "[${red}${bold}ERROR${reset}] Some of the required variables have not been provided."
    exit 1
fi

docker exec ${DOCKER_CONTAINER_NAME} pg_dump -U postgres -s --schema=${SCHEMA_NAME} --disable-dollar-quoting postgres > ${FILENAME}