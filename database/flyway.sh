#!/usr/bin/env bash
set -e

PATH_ESCAPE_CHAR=""

function set_path_escape_char() {
    if [[ "$UNAME" == CYGWIN* || "$UNAME" == MINGW* ]] ; then
        PATH_ESCAPE_CHAR="/"
    fi
}

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

CWD=$(wd)

function intro_message() {
    echo -e "
        ${green}${bold}Elsevier Submission Tracker${reset}: ${bold}Flyway EM Replication DB Migrator${reset}
    "
}

function complete_message() {
    echo -e "
        flyway ${FLYWAY_COMMAND}: [${bold}${green}COMPLETE${reset}]
    "
}

function sleeping_for() {
    echo -e "Sleeping for ${bold}$1 seconds${reset}."
    sleep $1
}

# Get configuration for this script to run
function source_script_config() {
    source ${CWD}/utils/.colors
    source ${CWD}/script.conf
}

# Cleanup existing flyway configuration
function clear_flyway_conf() {
    echo -e "Cleaning up configuration: [${bold}${FLYWAY_CONF_LOCATION}${reset}]"
    rm -rf ${FLYWAY_CONF_LOCATION}/*
}

# Create flyway.conf file from variables in script.conf and environment variables
function create_flyway_conf() {
    echo -e "Creating ${bold}${yellow}flyway.conf${reset}: [${bold}${white}${FLYWAY_CONF_LOCATION}${reset}]"
    mkdir -p ${FLYWAY_CONF_LOCATION}/ && clear_flyway_conf
    eval "echo \"$(cat ${FLYWAY_CONF_TEMPLATE_LOCATION})\"" > ${FLYWAY_CONF_LOCATION}/flyway.conf
}

# Create local flyway.conf file from variables in script.conf and environment variables
function create_local_flyway_conf() {
    echo -e "Creating local ${bold}${yellow}flyway.conf${reset}: [${bold}${white}${FLYWAY_CONF_LOCATION}${reset}]"
    mkdir -p ${FLYWAY_CONF_LOCATION}/ && clear_flyway_conf
    eval "echo \"$(cat ${FLYWAY_CONF_LOCAL_TEMPLATE_LOCATION})\"" > ${FLYWAY_CONF_LOCATION}/flyway.conf
}

function run_flyway_command() {
    echo -e "Running [${bold}${yellow}flyway ${FLYWAY_COMMAND}${reset}]"
    docker run --rm --net=host -v ${CWD}/sql/:/flyway/sql -v ${CWD}/conf/:/flyway/conf flyway/flyway $1
}

set_path_escape_char
source_script_config

if [ -z "${FLYWAY_COMMAND}" ] 
then
    echo -e "No Flyway Command Provided. [${bold}${yellow}Skipping...${reset}]"
else
    intro_message
    create_flyway_conf
    run_flyway_command ${FLYWAY_COMMAND}
    clear_flyway_conf
    complete_message
fi
