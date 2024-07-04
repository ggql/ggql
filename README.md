# ggql

[![Build Status](https://github.com/ggql/ggql/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/ggql/ggql/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/ggql/ggql/branch/main/graph/badge.svg?token=El8oiyaIsD)](https://codecov.io/gh/ggql/ggql)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggql/ggql)](https://goreportcard.com/report/github.com/ggql/ggql)
[![License](https://img.shields.io/github/license/ggql/ggql.svg)](https://github.com/ggql/ggql/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/ggql/ggql.svg)](https://github.com/ggql/ggql/tags)



## Introduction

*ggql* is forked from [GQL](https://github.com/AmrDeveloper/GQL) and written in Go.

ggql is a query and mutate language with a syntax very similar to SQL with a tiny engine to
perform queries or mutation on .git files instance of database files, note that all keywords
in ggql are case-insensitive similar to SQL.



## Prerequisites

- Go >= 1.20.0



## Run

```bash
make build

# Query
./bin/ggql -q "select * from branches" -r /path/to/git/repo
./bin/ggql -q "select * from commits where author_name=joe" -r /path/to/git/repo
./bin/ggql -q "select * from projects where name=test" -r /path/to/git/repo

# Mutate
./bin/ggql -m "update project set name=test" -r /path/to/git/repo
./bin/ggql -m "insert into project (change_id, name) values (1, project)" -r /path/to/git/repo
./bin/ggql -m "delete from change where change_id=1" -r /path/to/git/repo
```



## Usage

```
ggql is a SQL like query and mutate language to run on local repositories

Usage: ggql [OPTIONS]

Options:
-r,  --repos <REPOS>        Path for local repositories to run query on
-q,  --query <GQL Query>    GitQL query to run on selected repositories
-m,  --mutate <GQL Mutate>  GitQL mutate to run on selected repositories
-p,  --pagination           Enable print result with pagination
-ps, --pagesize             Set pagination page size [default: 10]
-o,  --output               Set output format [render, json, csv]
-a,  --analysis             Print Query analysis
-h,  --help                 Print GitQL help
-v,  --version              Print GitQL Current Version
```



## License

Project License can be found [here](LICENSE).



## Reference

- [GQL](https://github.com/AmrDeveloper/GQL)
