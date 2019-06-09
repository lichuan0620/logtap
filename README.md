# LogTap

Open the tap and get a flow of log! Logtap is a benchmark tool that generates log messages in a controlled way.

## Quick Look

Logtap provides many ways to define how the log messages should be generated. The simplest way to get started without worrying about all the settings is via a template. For example, run:

```
logtap -t Roast
```

LogTap would start with the `Roast` template which produces 40 log messages per second, each contains 1 MiB of random characters.

## Build and Run

#### Prerequsites

[Git](https://git-scm.com/)  
[Go (at least Go 1.11, earlier version not tested)](https://golang.org/dl/)  
[Docker (required only if you want to run via docker)](https://docs.docker.com/install/) 

#### Using Docker

The simplest way.

```
# pull the image
docker pull lichuan0620/logtap:latest

# get help
docker run lichuan0620/logtap:latest -h

# run
docker run lichuan0620/logtap:latest -t Roast
``` 

#### Build from source

If you don't have or don't want to use Docker.

```
# clone the source code from GitHub
git clone https://github.com/lichuan0620/logtap.git

# build the bin
make build-local

# get help
bin/logtap -h

# run
bin/logtap -t Roast
```

## Configuration

If the templates don't satisfy you, you can configure LogTap to generate more specific workloads. 

To get started, first run `logtap -h` to check the help messages. It'll show you all the configurable flags and the special constants such as the names of the templates. 

You can override the default value of almost all command line flags using environment variables. The environment variables are all in the format of `LOGTAP_NAME_OF_THE_FLAG`. For example, to set the default value for `--output.filePath`, which is the path of the log file to which the log messages should be appended, you can set the `LOGTAP_OUTPUT_FILE_PATH` environment variable.
