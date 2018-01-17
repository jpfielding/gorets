gorets
======

RETS client in Go

[![Build Status](https://travis-ci.org/jpfielding/gorets.svg?branch=master)](https://travis-ci.org/jpfielding/gorets)

The attempt is to meet 1.8.0 compliance.

http://www.reso.org/assets/RETS/Specifications/rets_1_8.pdf.

Find me at gophers.slack.com#gorets


There are multiple projects in the repository:

gorets/rets - provides a go based client for RETS

gorets/metadata - provides the common structure for reading in properly formed RETS metadata

gorets/explorer - provides a go backend for a reactjs ui for browsing RETS servers

gorets/proxy - provides a mechanism for proxying multiple RETS connections through a single endpoint

gorets/syndication - provides the RETS syndication struct for processing syndication feeds 

gorets/utils - helper tools for dealing with RETS