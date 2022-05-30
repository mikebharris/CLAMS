import contextlib
import os
import shutil
from os import path

from invoke import task, run as local
import vars


@task
def terraform(context, project_name=vars.project_name, application_name=vars.application_name, environment='nonprod', mode='plan'):
    if mode not in ['plan', 'apply', 'destroy']:
        print("No action to take.  Try passing --mode plan|apply|destroy")
        exit(-1)

    if project_name == '' or application_name == '':
        print("I cannot work with empty values for --project_name and --application_name")
        exit(-1)

    if mode == 'apply':
        build_lambdas()

    account_number = '215048116110'

    bucket = 'elsevier-tio-aws-rap-ssi-{environment}-{account_number}' \
        .format(environment='nonprod', account_number=account_number)

    key = 'tfstate/{environment}-{project_name}-{application_name}.json' \
        .format(environment=environment,
                project_name=project_name,
                application_name=application_name)

    print("Remote state is {bucket}/{key}".format(bucket=bucket, key=key))
    terraform_init(bucket, key, 'us-east-1')

    command = 'terraform {mode} -input=false ' \
              '-var "account_number={account_number}" -var "environment={environment}" --refresh=true' \
        .format(mode=mode,
                account_number=account_number,
                environment=environment)

    with do_in_directory('terraform'):
        local(command)


def build_lambdas():
    for f in ['attendees-api', 'registrar']:
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


@contextlib.contextmanager
def do_in_directory(path):
    CWD = os.getcwd()
    os.chdir(path)
    try:
        yield
    finally:
        os.chdir(CWD)
