#!/usr/bin/env bash
PG_SANITY_BASE_DOCKER="python:3.5-alpine"
PG_SANITY_CONTAINER_NAME="pgsanity"
PG_SANITY_PATH="//pgsanity"
PG_SANITY_SQL_PATH="//sql"
IS_ERROR=false

# Get current working directory for relevant OS
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

function clear_existing_docker_containers_if_any() {
    echo "Removing docker containers"
    docker rm -f ${PG_SANITY_CONTAINER_NAME} || true
}

function create_pg_sanity_base_container() {
    echo "Creating pg_sanity base container"
    docker run --name ${PG_SANITY_CONTAINER_NAME} -t -d -v ${CWD}/pgsanity/:${PG_SANITY_PATH} -v ${CWD}/sql/:${PG_SANITY_SQL_PATH} ${PG_SANITY_BASE_DOCKER};
}

function run_shell_script_in_container() {
    echo "Running script [${bold}$1${reset}]"
    docker exec ${PG_SANITY_CONTAINER_NAME} //bin/sh $1
}

function run_pgsanity_in_container() {
    echo "Running script [${bold}$1${reset}]"
    docker exec ${PG_SANITY_CONTAINER_NAME} //usr/local/bin/pgsanity $1
}

function pgsanity_sql_dir() {
    for f in sql/*.sql
    do
        run_pgsanity_in_container "${PG_SANITY_SQL_PATH}/$(basename "$f")"
        EXIT_CODE=$?
        if [ $EXIT_CODE != 0 ]; then
            IS_ERROR=true
        fi
    done
}

CWD=$(wd)

source ${CWD}/utils/.colors

clear_existing_docker_containers_if_any

create_pg_sanity_base_container

run_shell_script_in_container "${PG_SANITY_PATH}/install-pgsanity.sh"

pgsanity_sql_dir

clear_existing_docker_containers_if_any

if [ $IS_ERROR == true ]; then
    echo -e "[${red}${bold}ERROR${reset}] Some of the scripts failed the sanity check."
    exit 1
fi