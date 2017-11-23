# AWS Assume

Command line application to assume a role using a config.ini and config.creds file, the config.ini utilises the same file format as [aws-extend-switch-roles](https://github.com/tilfin/aws-extend-switch-roles). This then merges a config.creds file to authorise assume roles

## Getting Started

Create a config file that can be used in the format from aws-extend-switch-roles, create a config.creds file with the following values.

    # Master account
    [profile master]
    aws_access_key_id = AWS_ACCESS_KEY_ID
    aws_secret_access_key = AWS_SECRET_ACCESS_KEY
    secret = TOTP_SECRET
    region = ap-southeast-2


[![Codefresh build status]( https://g.codefresh.io/api/badges/build?repoOwner=s3than&repoName=assume&branch=master&pipelineName=assume&accountName=s3than&type=cf-1)]( https://g.codefresh.io/repositories/s3than/assume/builds?filter=trigger:build;branch:master;service:5a14ca51907b050001a46ac0~assume)

<!-- ### Prerequisites

NA

### Installing

A step by step series of examples that tell you have to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo

## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system -->

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/s3than/assume/tags).

## Authors

* **Tim Colbert ** - *Initial work* - [S3than](https://github.com/s3than)

See also the list of [contributors](https://github.com/s3than/assume/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for detail

