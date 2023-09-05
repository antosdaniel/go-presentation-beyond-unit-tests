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

## Requirements

- Go 1.21
- Docker Desktop 24+

## Start application

```shell
docker compose up
```
