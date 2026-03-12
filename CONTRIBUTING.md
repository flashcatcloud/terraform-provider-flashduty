# Contributing

[fork]: https://github.com/flashcatcloud/terraform-provider-flashduty/fork
[pr]: https://github.com/flashcatcloud/terraform-provider-flashduty/compare
[style]: https://github.com/flashcatcloud/terraform-provider-flashduty/blob/main/.golangci.yml

Hi there! We're thrilled that you'd like to contribute to this project. Your help is essential for keeping it great.

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE).

## Prerequisites for running and testing code

These are one time installations required to be able to test your changes locally as part of the pull request (PR) submission process.

1. Install [Go](https://go.dev/doc/install) >= 1.25.5
1. Install [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
1. Install [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation)

## Submitting a pull request

1. [Fork][fork] and clone the repository
1. Make sure the tests pass on your machine: `make test`
1. Make sure linter passes on your machine: `make lint`
1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and make sure the tests and linter still pass
1. If you've added or changed resources/data sources, update documentation: `make generate`
1. Push to your fork and [submit a pull request][pr]
1. Pat yourself on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Follow the [style guide][style].
- Write tests. See [Running Acceptance Tests](#running-acceptance-tests) below.
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests.
- Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Running Acceptance Tests

Acceptance tests create real resources against the Flashduty API. You will need a valid APP key:

```shell
export FLASHDUTY_APP_KEY="your-app-key"
```

Run the full suite:

```shell
make testacc
```

Run a specific test:

```shell
make testacc-run TEST=TestAccTeamResource
```

Run only data source or resource tests:

```shell
make testacc-datasources
make testacc-resources
```

### Test prerequisites

Some tests require existing resources in your Flashduty account:
- `integration_id`: Required for Route resource/data source tests (created via Flashduty UI)
- `template_id`: Required for EscalateRule resource tests
- `member_id` / `person_ids`: Required for Member data source and Schedule resource tests

The following test files contain hardcoded IDs that should be updated before running tests:
- `internal/provider/flashduty_escalate_rule_resource_test.go`: `template_id`
- `internal/provider/flashduty_schedule_resource_test.go`: `person_ids`
- `internal/provider/flashduty_route_resource_test.go`: `integration_id`
- `internal/provider/flashduty_data_sources_test.go`: `member_id`, `integration_id`

## API implementation notes

These notes are for contributors working on the provider internals.

| Resource | Read Implementation |
|----------|---------------------|
| `flashduty_inhibit_rule` | Uses List API + filter by ID |
| `flashduty_silence_rule` | Uses List API + filter by ID |
| `flashduty_escalate_rule` | Uses `/info` API directly |
| `flashduty_route` | No delete API - test skipped, manual cleanup required |
| `flashduty_template` | Uses `/info` API directly |
| `flashduty_alert_pipeline` | Uses `/info` API directly; delete upserts empty rules |

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
