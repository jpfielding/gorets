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

  * [gorets/pkg/rets](rets) - provides a Go based client for RETS

  * [gorets/pkg/metadata](metadata) - provides the common structure for reading in properly formed RETS metadata

  * [gorets/pkg/syndication](syndication) - provides the RETS syndication struct for processing syndication feeds 

  * [gorets/pkg/util](util) - helper tools for dealing with RETS

## Web Tools

  * [gorets/pkg/explorer](explorer) - provides a Go backend for a ReactJS UI for browsing RETS servers

  * [gorets/pkg/proxy](proxy) - provides a mechanism for proxying multiple RETS connections through a single endpoint

