# rt-retention

A JFrog CLI plugin to facilitate artifact retention in Artifactory.

⚠️ **Work in progress** ⚠️
Don't point this at your production instance

## TL;DR

Deletes artifacts matching all [File Specs](https://www.jfrog.com/confluence/display/JFROG/Using+File+Specs) found in a given directory.

Allows for easy templating of retention policies through a JSON configuration file.

## Installation

This plugin isn't currently hosted anywhere yet, so you'll be building it locally.

You can use the [build.sh](scripts/build.sh) and [install.sh](scripts/install.sh) scripts.

## Usage

### Commands

- run
  - Usage: `jf rt-retention run [command options] <filespecs-path>`

  - Arguments:
      - filespecs-path    _(Path to the filespecs file/dir)_

  - Options:
    - --dry-run    _disable communication with Artifactory [Default: **true**]_
    - --verbose    _output verbose logging [Default: false]_
    - --recursive    _recursively find filespecs files in the given dir [Default: false]_

- expand
  - Usage: `jf rt-retention expand [command options] <subscriptions-path> <templates-path> <output-path>`
  
  - Arguments:
    - subscriptions-path    _(Path to the subscriptions JSON file)_
    - templates-path    _(Path to the templates dir)_
    - output-path    _(Path to output the generated filespecs)_

  - Options:
    - --verbose      _output verbose logging [Default: false]_
    - --recursive    _recursively find templates in the given dir [Default: false]_
