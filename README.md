# Introduction

This is the example app that I built as part of a blog tutorial series available here: https://www.kriskelly.me/part1-go-k8s/

This is not meant to be a production-ready application. It is meant to showcase a basic integration of several technologies, namely Go, Kubernetes, Tilt, Dgraph, and S3, with potentially others thrown in there as I go.

# Installation

If you have been following the tutorial and want to install the app, first clone the project:

```bash
$ git clone git@github.com:kriskelly/dating-app-example.git
$ cd dating-app-example
```

Then follow the instructions in the post above for installing Docker, Kubernetes, and Tilt. Then you should be able to run:

```
$ tilt up
```

And then access the GraphQL playground at `localhost:3000`.

# Contributing

If you see anything horribly wrong with my code, please let me know by opening a pull request here, and I'll be sure to update the code and my blog posts to reflect the fix.
