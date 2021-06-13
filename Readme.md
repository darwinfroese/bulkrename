# bulkrename

A very simple tool hacked together to solve a very specific problem that I've encountered
too many times: renaming a large number of files and their directories.

## Building the tool
This tool only relies on the standard library so just run `go build -o bulkrename main.go` to build it.

## Usage

The tool is very straight forward to use you just navigate to the root directory that
you want to walk through and rename the files of and execute it. The following flags
can be set to modify behaviour:

* `-r` will recurse through all subdirectories and update them accordingly
* `-w` will remove whitespace from names
* `-p "regex" ` (can be passed multiple times) a regex expression that will be removed
from the file name
* `-n` will do everything except the actual rename so you can see what is going to be
performed.

## Improvements

So many improvements can and should be made. Right now it just does the job and that's
where I'll leave it, for now.
