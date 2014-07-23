package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
)

var (
	port        *int    = flag.Int("port", 8053, "Port to dial on. Default: 8053")
	host        *string = flag.String("host", "localhost", "The hostname of the server reading from the Twitter gradenhose.")
	outputDir   *string = flag.String("dir", "./", "Path to store downloaded tweets.")
	compression *bool   = flag.Bool("compress", true, "Perform gzip compression before writing to disk.")

	fileExtension string
	outFile       *os.File // The file we're currently writing to
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	flag.Parse()
	log.Println("Using output directory:", *outputDir)

	// Catch Ctrl + C, make sure we clean up.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			outFile.Close()
			os.Exit(0)
		}
	}()

	if *compression == true {
		fileExtension = ".json.gz"
		log.Println("Using GZIP compression.")
	} else {
		fileExtension = ".json"
	}

	// Can be toggled by the program to stop repeated connection failed messages.
	// This should make debugging easier in the case that a connection drops.
	printConnectionError := true

	for {
		// Atempt to connect to the server
		conn, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
		if err != nil {
			if printConnectionError {
				log.Println("Could not dial remote server on " + *host + ":" + strconv.Itoa(*port))
				printConnectionError = false
			}
			time.Sleep(1 * time.Second)
			continue
		}
		printConnectionError = true
		log.Println("Connected to server " + *host + ":" + strconv.Itoa(*port))

		reader := bufio.NewReader(conn)

		// Start wrting the recorded data to a file
		currentHour := -1
		for {
			now := time.Now()
			year, month, day := now.Date()
			hour, min, sec := now.Clock()

			// New hour, new file.
			if hour != currentHour || outFile == nil {
				if outFile != nil {
					outFile.Close()
				}

				filename := fmt.Sprintf("%4d-%02d-%02d_%02d-%02d-%02d"+fileExtension, year, month, day, hour, min, sec)
				outFile, err = os.Create(*outputDir + "/" + filename)
				if err != nil {
					log.(err)
				}
				log.Println("Creating new output file..." + filename)
				currentHour = hour
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			data := []byte(line)
			if *compression == true {
				var b bytes.Buffer
				w, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
				w.Write(data)
				w.Close()
				data = b.Bytes()
			}

			outFile.Write(data)
		}
	}

}
