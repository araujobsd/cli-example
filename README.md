[![GoDoc](https://godoc.org/github.com/araujobsd/cli-example/plugins?status.svg)](https://godoc.org/github.com/araujobsd/cli-example/)
[![GitHub issues](https://img.shields.io/github/issues/araujobsd/cli-example.svg)](https://github.com/araujobsd/cli-example/issues)
[![GitHub forks](https://img.shields.io/github/forks/araujobsd/cli-example.svg)](https://github.com/araujobsd/cli-example/network)
[![Go Report Card](https://goreportcard.com/badge/github.com/araujobsd/cli-example)](https://goreportcard.com/report/github.com/araujobsd/cli-example)
[![CircleCI](https://circleci.com/gh/araujobsd/cli-example.svg?style=svg)](https://circleci.com/gh/araujobsd/cli-example)

cli-example
================
It is a CRUD CLI implementation using Go1.12 that simulates an user listing products. It uses CSV to store the data, and all CRUD commands and others can be implemented dinamically without the need to restart the main program to load it.

How to build?
================
You can simple:
```sh build.sh```
OR
```make```

How to run?
================
Proceed using ```run.sh```
OR
```./thecarousell```

List of commands implemented
================
- `register`: Register an user. Only registered users can use the additional commands.
```
Usage:
   REGISTER user1
```

- `create_listing`: Add an item for sale
```
Usage:
   CREATE_LISTING user1 'Phone model 8' 'Black color, brand new' 1000 'Electronics'
```

- `delete_listing`: Remove an item for sale
```
Usage:
   DELETE_LISTING user1 itemID
```

- `get_listing`: Find and print to stdout the item
```
Usage:
   GET_LISTING user1 itemID
```

- `get_category`: Find and print to stdout all the items from a follow category
```
Usage:
   GET_CATEGORY user1 category {sort_price|sort_time} {asc|dsc}
```

- `get_top_category`: Find an user category with the most items
```
Usage:
   GET_TOP_CATEGORY user1
```

#### NOTE
We do a normalization on the command line, if you type REGISTER or ReGiStEr, we will find the right command for you. Also you can run the commands as a standalone command, just get into ```commands``` and run it using the same parameters above.

## Copyright and licensing
Distributed under [2-Clause BSD License](https://github.com/araujobsd/cli-example/blob/master/LICENSE).
