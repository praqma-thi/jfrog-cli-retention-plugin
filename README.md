> **Note**
>
> This repository has moved to [verifa/rt-retention](https://github.com/verifa/rt-retention)

# rt-retention

A JFrog CLI plugin to facilitate artifact retention in Artifactory.

## TL;DR

Deletes artifacts matching all [File Specs](https://www.jfrog.com/confluence/display/JFROG/Using+File+Specs) found in a given directory.

Allows for generation of FileSpecs files through Go templates and a JSON configuration file.

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
  - Usage: `jf rt-retention expand [command options] <config-path> <templates-path> <output-path>`
  
  - Arguments:
    - config-path    _(Path to the JSON config file)_
    - templates-path    _(Path to the templates dir)_
    - output-path    _(Path to output the generated filespecs)_

  - Options:
    - --verbose      _output verbose logging [Default: false]_
    - --recursive    _recursively find templates in the given dir [Default: false]_

## Templating

This plugins allows you to generate retention policies using Go templates and a JSON config file.

### Templates

Templates use values from the JSON config file to generate FileSpec files.

`delete-older-than.json`:
```json
{
    "files": [{
        "aql": {
            "items.find": {
                "repo": "{{.Repo}}",
                "created" : {"$before" : "{{.Time}}"}
            }
        }
    }]
}
```

### JSON config

The JSON config file contains a key for each template, with an array of entries for that template.
Each entry will result in a FileSpecs file being generated.

If the entry has a **Name** property, it's value will be used as the FileSpecs file name.

`config.json`:
```json
{
    "delete-everything": [
        { "Name": "foo-dev", "Repo": "foo-dev-local" },
        { "Name": "bar-dev", "Repo": "bar-dev-local" }
    ],
    "delete-older-than": [
        { "Name": "baz-dev", "Repo": "baz-dev-local", "Time": "30d" }
    ]
}
```
