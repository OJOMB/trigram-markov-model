// CLI app for generating pseudo-random novel text learnt from a given input corpus

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/OJOMB/trigram-markov-model/markovmodel"
	"github.com/OJOMB/trigram-markov-model/trigram"
)

var outputNumWords = flag.Int("n", 100, "desired length in words of output")
var output = flag.String("output", "out.txt", "Name of file in which to persist output")
var inputText = flag.String("input", "text/trump.txt", "Filepath of the input text from which the model is built")
var cpuProfile = flag.String("cpuprofile", "", "Profile the CPU usage")
var memoryProfile = flag.String("memprofile", "", "Profile the memory usage")
var tracing = flag.Bool("trace", false, "Run go tool execution tracing")

func main() {
	flag.Parse()

	// CPU profiling code if applicable
	if *cpuProfile != "" {
		cpuProfFile, err := os.Create(*cpuProfile)
		defer cpuProfFile.Close()
		if err != nil {
			log.Fatal("Failed to create cpu profile file: ", err)
		}
		if err = pprof.StartCPUProfile(cpuProfFile); err != nil {
			log.Fatal("Failed to start CPU profile: ", err)
		}
		fmt.Printf("Running CPU profiling, outputting to: %s\n", cpuProfFile.Name())
		defer pprof.StopCPUProfile()
	}

	// start tracing if applicable
	if *tracing {
		traceFile, err := os.Create("trace.out")
		defer traceFile.Close()
		if err != nil {
			log.Fatal("Failed to create trace file: ", err)
		}
		trace.Start(traceFile)
	}

	// set the random seed
	rand.Seed(time.Now().UnixNano())
	trigrams, err := trigram.ParseFileToNormalisedTrigrams(*inputText)
	if err != nil {
		log.Fatal(err)
	}

	// Init the model
	var mm *markovmodel.Model = markovmodel.New(rand.Intn)
	for _, t := range trigrams {
		mm.Add(t)
	}

	// generate output
	result, err := mm.Generate(*outputNumWords)
	if err != nil {
		log.Fatal("Failed to generate novel text from markov model: ", err)
	}

	// Write generated output to file
	outputFile, err := os.Create(*output)
	defer outputFile.Close()
	if err != nil {
		log.Fatal("Failed to create output file: ", err)
	}
	w := bufio.NewWriter(outputFile)
	bytesWritten, err := w.WriteString(result)
	if err != nil {
		log.Fatal("Failed to write output to file: ", err)
	}
	fmt.Printf("%d bytes written to output file: %s\n", bytesWritten, outputFile.Name())

	// flush the buffer to check all content has been written to the output file
	w.Flush()

	if *tracing {
		trace.Stop()
	}

	// memory profiling code if applicable
	if *memoryProfile != "" {
		memProfFile, err := os.Create(*memoryProfile)
		if err != nil {
			log.Fatal("Failed to create memory profile: ", err)
		}
		defer memProfFile.Close()
		runtime.GC() // get up-to-date statistics
		if err = pprof.WriteHeapProfile(memProfFile); err != nil {
			log.Fatal("Failed to write memory profile: ", err)
		}
		fmt.Printf("Running memory profiling, outputting to: %s\n", memProfFile.Name())
	}
}
