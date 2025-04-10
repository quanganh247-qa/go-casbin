# Project my-casbin

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```


```
my-casbin
├─ .air.toml
├─ .goreleaser.yml
├─ Makefile
├─ README.md
├─ cmd
│  └─ api
│     ├─  rbac_model.conf
│     └─ policy.csv
├─ go.mod
├─ go.sum
└─ internal
   ├─ database
   │  └─ database.go
   └─ server
      ├─ routes.go
      ├─ routes_test.go
      └─ server.go

```