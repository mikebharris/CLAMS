# CLAMS

CLAMS => "BAMS in the Cloud"

A personal learning project using a connection to a legacy event management system as a way of illustrating serverless architectures using Go and Terraform and AWS.

![The architecture of CLAMS](CLAMS-architecture.png)

To use CLAMS, get the API Gateway endpoint via AWS Console; it's also displayed as the output of the deployment script (see below).

# Running Tests

There are service/integration-level tests that use Gherkin syntax to test integration between the Lambda and other dependent AWS servies.  The tests make use of Docker containers to emulate the various services locally, and therefore you need Docker Desktop running.

To run the service tests, change to the service in the _functions_ directory and type:

```shell
AWS_SECRET_ACCESS_KEY=x AWS_ACCESS_KEY_ID=y make int-test
```

# Deploying

There is a Python Fabric 2 script to help you do this.  First authenticate with AWS and then run the following from the command line (changing _mode_ from _plan_, _apply_, or _destroy_ and setting the other variables:

```shell
AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=YYYY fab terraform --account-number=111111111111 --contact=your@email.com --mode=plan

AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=YYYY fab terraform --account-number=111111111111 --contact=your@email.com --mode=apply

AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=YYYY fab terraform --account-number=111111111111 --contact=your@email.com --mode=destroy

```

The command line supports the following:

```shell
Usage: fab [--core-opts] terraform [--options] [other tasks here ...]

Docstring:
  none

Options:
  -a STRING, --account-number=STRING
  -c STRING, --contact=STRING
  -d STRING, --distribution-bucket=STRING
  -e STRING, --environment=STRING
  -i STRING, --input-queue=STRING
  -m STRING, --mode=STRING
  -p STRING, --project-name=STRING
  -r STRING, --region=STRING
  -t STRING, --attendees-table=STRING
```

# TODO list

* Write a decent front-end
* Add authentication to the API
* Put behind CloudFront
* Add provisioning for DNS hostname in Route53
