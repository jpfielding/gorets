gorets
======

RETS client in Go

[![Build Status](https://travis-ci.org/jpfielding/gorets.svg?branch=master)](https://travis-ci.org/jpfielding/gorets)

The attempt is to meet [RETS 1.8.0](https://www.reso.org/specifications/) compliance.

Find me at gophers.slack.com#gorets


There are **multiple projects** in this repository:

[gorets/rets](rets) - provides a Go based client for RETS

[gorets/metadata](metadata) - provides the common structure for reading in properly formed RETS metadata

[gorets/explorer](explorer) - provides a Go backend for a ReactJS UI for browsing RETS servers

[gorets/proxy](proxy) - provides a mechanism for proxying multiple RETS connections through a single endpoint

[gorets/syndication](syndication) - provides the RETS syndication struct for processing syndication feeds 

[gorets/util](util) - helper tools for dealing with RETS