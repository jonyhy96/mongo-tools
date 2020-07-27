// Copyright (C) MongoDB, Inc. 2014-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Main package for the mongorestore tool.
package main

import (
	"runtime/trace"

	"github.com/jonyhy96/mongo-tools/mongorestore"
	"github.com/mongodb/mongo-tools-common/log"
	"github.com/mongodb/mongo-tools-common/signals"
	"github.com/mongodb/mongo-tools-common/util"

	"os"
)

var (
	VersionStr = "built-without-version-string"
	GitCommit  = "build-without-git-commit"
)

func main() {
	trace.Start(os.Stderr)
	defer trace.Stop()

	opts, err := mongorestore.ParseOptions(os.Args[1:], VersionStr, GitCommit)

	if err != nil {
		log.Logvf(log.Always, "error parsing command line options: %s", err.Error())
		log.Logvf(log.Always, util.ShortUsage("mongorestore"))
		os.Exit(util.ExitFailure)
	}

	// print help or version info, if specified
	if opts.PrintHelp(false) {
		return
	}

	if opts.PrintVersion() {
		return
	}

	restore, err := mongorestore.New(opts)
	if err != nil {
		log.Logvf(log.Always, err.Error())
		os.Exit(util.ExitFailure)
	}
	defer restore.Close()

	finishedChan := signals.HandleWithInterrupt(restore.HandleInterrupt)
	defer close(finishedChan)

	result := restore.Restore()
	if result.Err != nil {
		log.Logvf(log.Always, "Failed: %v", result.Err)
	}

	if restore.ToolOptions.WriteConcern.Acknowledged() {
		log.Logvf(log.Always, "%v document(s) restored successfully. %v document(s) failed to restore.", result.Successes, result.Failures)
	} else {
		log.Logvf(log.Always, "done")
	}

	if result.Err != nil {
		os.Exit(util.ExitFailure)
	}
	os.Exit(util.ExitSuccess)
}
