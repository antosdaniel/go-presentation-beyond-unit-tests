# Go presentation: Beyond unit tests

A dive into integration and end-to-end testing in Go.

## Walkthrough

[app_to_test](./app_to_test) contains very simple web application. It's expense
tracker, with small HTTP API. Don't take too many of these ideas to production 🙃
You can play with it through [expenses.http](app_to_test%2Fexpenses.http).

First suite of test to check out is [test_repos](./test_repos). It uses test 
containers to spin up Postgres database and checks if repository works as promised.
It's a good tool for more complex SQL queries, or testing against concurrency issues.

Next one is [test_http_api](test_http_api). It includes database trick from
previous suite, but also spins up API. Explored idea here is how to keep tests
like these easy to read, and not a chore to write (at least, after the first one).

Final one is [test_e2e](test_e2e). We are spinning up whole application, and
running tests against that. To make things even more interesting, we are stubing
external API for better reliability. This checks that our application starts correctly, 
and that all the components work together. It's certainly more complex, not as
quick to write, nor run. Use them sparingly, to make sure that crucial parts
of your application work as expected. These could also be named smoke tests.

## Requirements

- Go 1.21
- Docker Desktop 24+

## Start application

```shell
docker compose up
```

## Run tests

Run unit tests:

```shell
go test ./...
```

Run e2e tests:

```shell
PKGS=$(git grep --files-with-matches "//go:build e2e_tests" -- "*.go" | xargs dirname | sed 's,^,./,g' | sort -u)
go test -tags=e2e_tests $PKGS
```
