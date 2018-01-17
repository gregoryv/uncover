package main

import "github.com/gregoryv/stamp"

func init() {
    s := &stamp.Stamp{
	    Package: "main",
	    Revision: "c5b0562",
	    ChangelogVersion: "unknown",
    }
    stamp.Use(s)
}

