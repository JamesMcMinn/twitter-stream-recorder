# Twitter Stream Recorder
Twitter Stream Recorder saves the output from a stream to disk, creating a new file every hour.
By default, Twitter Stream Recorded creates GZIP compressed files, with one Tweet in JSON format
per line.

## Installing
Twitter Stream Recorder can be installed using a number of methods.

The easiest method is to simply run `main.go` using the `go run` command: 

    go run main.go

To build the application binary, `go build` can be used and will generate
the binary in your current directory.

    go build mirgit.dcs.gla.ac.uk/JamesMcMinn/twitter-stream-recorder


## Usage
The stream recorded has a number of paramaters, all of which are optional depending on your setup.

    -compress=true      Perform gzip compression before writing to disk.
    -dir="./"           Path to store downloaded tweets.
    -host="localhost"   The hostname to dial to.
    -port=8053          Port to dial on.

### Example Usage
The following example connects to a machine called juvented on a non-default port, specifies a directory
to store the tweets and enables compression.

    twitter-stream-recorder -host=juventud -port=56874 -dir=/local/jjnas01/Public/Commonwealth -compress=true

To following disbales compression and writes files to the current directory:

    twitter-stream-recorder -host=juventud -port=56874 -dir=./ -compress=false