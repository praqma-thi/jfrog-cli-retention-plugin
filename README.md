# rt-retention

A JFrog CLI plugin to facilitate retention in Artifactory.

⚠️ **Work in progress** ⚠️
Don't point this at your production instance

## TL;DR

Deletes artifacts matching all [File Specs](https://www.jfrog.com/confluence/display/JFROG/Using+File+Specs) found in a given directory.

Currently just prevents you from having to run [`jfrog rt delete`](https://www.jfrog.com/confluence/display/CLI/CLI+for+JFrog+Artifactory#CLIforJFrogArtifactory-DeletingFiles) multiple times, but a policy templating system is in the works.

## Installation

This plugin isn't currently hosted anywhere, so you'll be building it locally.

You can use the [build.sh](scripts/build.sh) and [install.sh](scripts/install.sh) scripts.

## Usage

### Commands

- run
  - Usage: `jfrog rt-retention run [command options] <filespecs-path>`

  - Arguments:
      - filespecs-path    _(Path to the filespecs file/dir)_

  - Options:
    - --dry-run    _disable communication with Artifactory [Default: **true**]_
    - --verbose    _output verbose logging [Default: false]_
    - --recursive    _recursively find filespecs files in the given dir [Default: false]_

