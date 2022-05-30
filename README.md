# eHAMS
A personal learning project using a connection to a legacy event management system as a way of illustrating serverless architectures using Go

# Deploying

There is a Python Fabric 2 script to help you do this.  First authenticate with AWS and then run the following from the command line (changing _mode_ from _plan_ to _apply_ and setting the other variables:

```shell
fab terraform --environment=nonprod --mode=plan --account-number=12345678901234 --contact=you@yourdomain.com --cost-code=12345 --distribution-bucket=lambda-distributions
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
  -o STRING, --cost-code=STRING
  -p STRING, --project-name=STRING
  -r STRING, --region=STRING
  -t STRING, --attendees-table=STRING
```