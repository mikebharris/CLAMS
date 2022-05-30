import contextlib
import os
import shutil
from os import path

from invoke import task, run as local


@task
def terraform(context, account_number="", contact="", cost_code="", distribution_bucket="",
              attendees_table="attendees-datastore", input_queue="attendee-input-queue", project_name="ehams",
              region="us-east-1", environment="nonprod", mode="plan"):
    if mode not in ['plan', 'apply', 'destroy']:
        print("No action to take.  Try passing --mode plan|apply|destroy")
        exit(-1)

    if mode == 'apply':
        build_lambdas()

    bucket = 'elsevier-tio-aws-rap-ssi-{environment}-{account_number}' \
        .format(environment='nonprod', account_number=account_number)

    key = 'tfstate/{environment}-{project_name}.json' \
        .format(environment=environment,
                project_name=project_name)

    print("Remote state is {bucket}/{key}".format(bucket=bucket, key=key))
    terraform_init(bucket, key, 'us-east-1')

    command = 'terraform {mode} -input=false ' \
              '-var "product={project_name}" -var "region={region}" -var "cost_code={cost_code}" ' \
              '-var "contact={contact}" -var "distribution_bucket={distribution_bucket}" ' \
              '-var "attendees_table_name={attendees_table_name}" -var "input_queue_name={input_queue}" ' \
              '-var "account_number={account_number}" -var "environment={environment}" --refresh=true' \
        .format(mode=mode,
                project_name=project_name,
                region=region,
                cost_code=cost_code,
                contact=contact,
                distribution_bucket=distribution_bucket,
                attendees_table_name=attendees_table,
                input_queue=input_queue,
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
