import contextlib
import os
import shutil
import re
from os import path

from invoke import task, run as local


@task
def terraform(context, account_number="", contact="", distribution_bucket="terraform-deployments",
              project_name="clams", region="us-east-1", environment="nonprod", lambdas="yes",
              frontend="yes", flywayonly="no", database="yes", mode="plan"):
    if mode not in ['init', 'plan', 'apply', 'destroy']:
        print("No action to take.  Try passing --mode init|plan|apply|destroy")
        exit(-1)

    if flywayonly == 'yes':
        apply_database_schema()
        exit(0)

    if mode == 'apply' and lambdas == 'yes':
        build_lambdas()

    bucket = '{account_number}-{distribution_bucket}' \
        .format(account_number=account_number, distribution_bucket=distribution_bucket)

    key = 'tfstate/{environment}-{project_name}.json' \
        .format(environment=environment,
                project_name=project_name)

    print("Remote state is {bucket}/{key}".format(bucket=bucket, key=key))

    if mode == 'init':
        terraform_init(bucket, key, 'us-east-1')
        exit(0)

    command = 'terraform {mode} -input=false ' \
              '-var "product={project_name}" -var "region={region}" ' \
              '-var "contact={contact}" -var "distribution_bucket={distribution_bucket}" ' \
              '-var "account_number={account_number}" -var "environment={environment}" --refresh=true' \
        .format(mode=mode,
                project_name=project_name,
                region=region,
                contact=contact,
                distribution_bucket=distribution_bucket,
                account_number=account_number,
                environment=environment)

    print("executing: {c}".format(c=command))

    with do_in_directory('terraform'):
        local(command)

    if mode == 'apply' and database == 'yes':
        apply_database_schema()

    if mode == 'apply' and frontend == 'yes':
        build_and_deploy_frontend()


@task
def frontend(context):
    build_and_deploy_frontend()


def build_and_deploy_frontend():
    api_url = get_api_url()
    with do_in_directory('frontend'):
        ensure_frontend_dependencies_are_installed()
        build_frontend(api_url)
    deploy_frontend()


def get_api_url() -> str:
    with do_in_directory('terraform'):
        result = local('terraform output')
        search = re.search("(https.*)", result.stdout)
        api_url = '{host}/clams'.format(host=search.group(0)[:len(search.group(0)) - 1])
    return api_url


def build_lambdas():
    for f in ['attendees-api', 'attendee-writer', 'processor', 'db-trigger']:
        lambda_location = 'functions/{function}'.format(function=f)
        print("Building lambda in {l}....".format(l=lambda_location))
        with do_in_directory(lambda_location):
            local('make target')


def terraform_init(bucket, key, region):
    remove_local_terraform_state_files_to_prevent_deploying_in_wrong_environment()
    with do_in_directory('terraform'):
        local('terraform init -backend-config="bucket={bucket}" -backend-config="key={key}" ' \
              '-backend-config="region={region}"'
              .format(bucket=bucket, key=key, region=region))


def remove_local_terraform_state_files_to_prevent_deploying_in_wrong_environment():
    with do_in_directory('terraform'):
        if path.exists('.terraform'):
            shutil.rmtree('.terraform')


def ensure_frontend_dependencies_are_installed():
    print("Installing frontend dependencies...")
    local('npm install')


def build_frontend(api_url: str):
    print("Building frontend using API url {api_url}...".format(api_url=api_url))
    local('echo API_GATEWAY_URL={api_url}>.env'.format(api_url=api_url))
    local('npm run build')


def deploy_frontend():
    local('aws s3 cp frontend/public s3://clams.events.hacktionlab.org --recursive')


def apply_database_schema():
    with do_in_directory('terraform'):
        db_host = local('terraform output -raw rds_database_host').stdout
        db_name = local('terraform output -raw rds_database_name').stdout
        db_username = local('terraform output -raw rds_database_username').stdout
        db_password = local('terraform output -raw rds_database_password').stdout

    with do_in_directory('flyway'):
        flyway_command = 'flyway -user={db_username} -password={db_password} ' \
                         '-url=jdbc:postgresql://{db_host}:5432/{db_name} migrate' \
            .format(db_username=db_username,
                    db_password=db_password,
                    db_host=db_host,
                    db_name=db_name)
        print("Executing command: ", flyway_command, "\n")
        local(flyway_command)


@contextlib.contextmanager
def do_in_directory(directory: str):
    CWD = os.getcwd()
    os.chdir(directory)
    try:
        yield
    finally:
        os.chdir(CWD)
