# snap collector plugin - procstat

[![Build Status](https://travis-ci.org/Staples-Inc/snap-plugin-collector-procstat.svg?branch=master)](https://travis-ci.org/Staples-Inc/snap-plugin-collector-procstat)

## Getting Started
The procstat collector takes a specified pid file and collects metrics from the process specified by that pid file.
### Build

Build the plugin by running make within the repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Run
#### Configuration
```
{
  "control": {
    "plugins": {
      "collector": {
        "procstat": {
            "all": {
                "files": "/tmp/example1.pid:example1,/tmp/example2.pid:example2"
            }
        }
      }
    }
  }
}
```
* Specify the pid file with the desired name identifier using the format: `"files": "<filepath>:<name>,<filepath>:<name>,..."`

## Roadmap
* Add pgrep functionality for tracking a process.
* Allow for pid files with mutiple pids separated by a new line

If you have suggestions please open up an issue or provide a pull request.
