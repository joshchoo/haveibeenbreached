# haveibeenbreached

A simple clone of [haveibeenpwned](https://haveibeenpwned.com/) written in Go, and implemented using the AWS stack and the [Serverless Framework](https://www.serverless.com/):

- Serverless framework
- Cloudformation
- Lambda
- DynamoDB
- Simple Queue Service

## Getting started

You'll need an AWS account before proceeding with the remaining steps.

**Install the Serverless Framework**

```bash
$ npm install -g serverless
```

**Configure your AWS credentials**

https://www.serverless.com/framework/docs/providers/aws/cli-reference/config-credentials/

**Deploy haveibeenbreached**

```bash
$ make deploy
```

## Lambda functions

### createBreach

Creates a new breach.

```bash
$ sls invoke -f createBreach -l --path exampleEvents/createBreach.json
```

### addAccountsToBreach

Adds a list of email accounts to an existing breach.

```bash
$ sls invoke -f addAccountsToBreach -l --path exampleEvents/addAccountsToBreach.json
```

### findAccount

Retrieves an existing email account that has been involved in any breaches.

```bash
$ sls invoke -f findAccount -l --path exampleEvents/findAccount.json
```

### notifyMe

Subscribes an email for breach notifications. It sends a subscription message to SQS.

```bash
$ sls invoke -f notifyMe -l --path exampleEvents/notifyMe.json
```

### sendSubscriptionEmail

Processes subscription messages on SQS, and sends subscription confirmation emails (not implemented) to subscribers.

### notifySubscribersOfBreach

Send emails to subscribers whose emails are involved in a given breach.

```bash
$ sls invoke -f notifySubscribersOfBreach -l --path exampleEvents/notifySubscribersOfBreach.json
```
