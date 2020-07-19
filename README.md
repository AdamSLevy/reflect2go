# package reflect2go

[![GoDoc](https://godoc.org/github.com/AdamSLevy/reflect2go?status.svg)](https://godoc.org/github.com/AdamSLevy/reflect2go)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdamSLevy/reflect2go)](https://goreportcard.com/report/github.com/AdamSLevy/reflect2go)
[![Build Status](https://travis-ci.org/AdamSLevy/reflect2go.svg?branch=master)](https://travis-ci.org/AdamSLevy/reflect2go)
[![Coverage Status](https://coveralls.io/repos/github/AdamSLevy/reflect2go/badge.svg?branch=master)](https://coveralls.io/github/AdamSLevy/reflect2go?branch=master)

Package reflect2go generates the Go code for a given reflect.Type.

This is useful when you want to generate the Go code for a dynamically
defined type, like with reflect.StructOf.

An example use case is a program that updates or manipulates the tags on the
fields of an existing struct.
