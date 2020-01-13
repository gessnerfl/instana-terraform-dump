# Instana Terraform Dump

[![Build Status](https://travis-ci.com/gessnerfl/instana-terraform-dump.svg?branch=master)](https://travis-ci.com/gessnerfl/instana-terraform-dump)

This project is a simple command line tool which read resources from Instana (<https://instana.io)> REST API
<https://instana.github.io/openapi/> and creates the corresponding Terraform resources for them. The resources
are created for the Terraform Instana Provider from <https://github.com/gessnerfl/terraform-provider-instana>.

## Supported Resources

* Custom Event Specifications

## Usage

```(bash)
instana-terraform-dump -key=<api-key> -host=<instana-host-name> -out=<output-file-path>
```

### Parameters

* key = the Instana API key to access Instana REST API
* host = the Instana Hostname (customer specific endpoint)
* out = the path to the output file

## How to make

* Install make
  * add make to path
* Install go
  * add go to path "gopath"
* run make from ./instana-terraform-dump
  * `make`

### Potential issues during make:

**Error Message:** `golint: command not found`

**Solution:** If not installed automatically, install `golint` command line tool: `go get -u golang.org/x/lint/golint`