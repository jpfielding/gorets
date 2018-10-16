gorets
======

RETS client in Go

[![Build Status](https://travis-ci.org/jpfielding/gorets.svg?branch=master)](https://travis-ci.org/jpfielding/gorets)
[![GoDoc](https://godoc.org/github.com/jpfielding/gorets?status.svg)](https://godoc.org/github.com/jpfielding/gorets)
[![Go Report Card](https://goreportcard.com/badge/github.com/jpfielding/gorets)](https://goreportcard.com/report/github.com/jpfielding/gorets)


The attempt is to meet [RETS 1.8.0](https://www.reso.org/specifications/) compliance.

Find me at gophers.slack.com#gorets


There are **multiple projects** in this repository:

## Client Tools

  * [Client](pkg/rets) - provides a Go based client for RETS

  * [Metadata](pkg/metadata) - provides the common structure for reading in properly formed RETS metadata

  * [Syndication](pkg/syndication) - provides the RETS syndication struct for processing syndication feeds 

  * [Util](pkg/util) - helper tools for dealing with RETS

## Web Tools

  * [Explorer Client](web/explorer) - provides a ReactJS UI for browsing RETS servers

  * [Explorer Service](pkg/explorer) - provides a Go backend for browsing RETS servers

  * [Proxy](pkg/proxy) - provides a mechanism for proxying multiple RETS connections through a single endpoint

