# Contributing to vault-aws-credentials-helper

Thank you for investing your time in contributing to the `vault-aws-credentials-helper`! The following is a set of guidelines for contributing to the tool.

## Reporting Bugs

Bugs are tracked as GitHub issues. When you create an issue, please provide the following information by filling in the template.

- **Be clear in the title of your issue** to help maintainers to identify the problem
- **Give a clear and concise description of your problem** to give more context to maintainers or contributors.
- **Describe the exact steps required to reproduce the problem** in as many details as possible.
- **Add your versions of `vault`, `AWS CLI` and `vault-aws-credentials-helper`**

## Suggesting enhancements

In case you want to suggest features or improvements for the `vault-aws-credentials-helper`, please follow this guideline to help the community and the maintainers understand your suggestion. Please, be sure that your suggestion isn't already requested by checking the [issue list](https://github.com/thomasgouveia/vault-aws-credentials-helper/issues).

- **Be clear in the title of your issue** to help maintainers to identify the problem
- **Give a clear and concise description of your problem** to give more context to maintainers or contributors.
- **Answer by `yes/no` if you are able to make a PR for resolving this issue and if you want to work on it**

## Contribution workflow

Please follow the below guideline in order to make contributions to `vault-aws-credentials-helper`.

### Fork the repository

First off, you need to fork the [thomasgouveia/vault-aws-credentials-helper](https://github.com/thomasgouveia/vault-aws-credentials-helper) repository in order to do your changes freely without affecting the original project.

> See [GitHub - Fork a repo](https://docs.github.com/en/issues/tracking-your-work-with-issues/creating-an-issue).

### Set up your laptop

Clone your fork on your laptop, and set up the tools needed to develop the project: 

- [Go](https://go.dev/doc/install) (1.20+)

For some fixes or features, you will need to configure a [Vault](https://www.vaultproject.io/) in order to test your changes before submitting them. The easiest way to do so is to use one of the [examples](./examples/) provided in the repository. For instance, generally using the [examples/vault-with-userpass](./examples/vault-with-userpass/) will be sufficient. Please refer to the README present in the example folder to learn how to run it.

### Make a branch for your changes

Once everything is set up on your laptop, you can create a new branch in your repository and make your changes. Once your work is done, push it to the branch and create a pull request on the [original repository](https://github.com/thomasgouveia/vault-aws-credentials-helper).

Please review the [convention for commits and pull requests](#conventions).

**If you're creating a PR that is not ready to review, please do not forget to mark it as a draft.**

### Review process

Once the PR will be ready to be reviewed, ensure it is marked as ready. One of the project maintainer will review it when possible. If everything is ok, your PR will be merged into the original repository. 

Congratulations! 🎉 You just achieved your first contribution. How about a do-over? 

## Conventions

### Commits

Please respect the following convention when comitting changes:

#### Message format

```bash
$type: Short description (#1234)

Longer description here if necessary

BREAKING CHANGE: only contain breaking change
Signed-off-by: $username <$email>
```

Where `$type` is one of the following :

* **feat**: A new feature
* **fix**: A bug fix
* **docs**: Documentation changes
* **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **perf**: A code change that improves performance
* **test**: Adding missing or correcting existing tests
* **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation
* **ci**: Changes to the CI process of the project

Below some examples of accepted commit messages:

```
feat: add support for github authentication method (#13)
```

```
ci: add format job (#101)
```

All commits required to be signed (`Signed-off-by:`) to be [DCO compliant](https://developercertificate.org/). To automatically add this line to your commits, append the option `-s` or `--signoff` to your `git commit` command. 

### Pull requests

Please use the following format for your PR title:

```bash
Short Description (#300)
# or if your PR solves multiple issues
Short Description (#301, #300, #302)
```

## Code of conduct

This project and everyone participating in it is governed by the [Code of Conduct](./CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
