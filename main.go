package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pboyd/unirecode/codec"
)

func main() {
	var decoderName, encoderName, output string
	flag.StringVar(&decoderName, "d", "", "decoder name")
	flag.StringVar(&encoderName, "e", "", "encoder name")
	flag.StringVar(&output, "o", "", "output file name")
	flag.Parse()

	if decoderName == "" || encoderName == "" {
		if decoderName == "" {
			fmt.Printf("%s: no decoder\n", os.Args[0])
		}

		if encoderName == "" {
			fmt.Printf("%s: no encoder\n", os.Args[0])
		}

		flag.Usage()
		os.Exit(1)
	}

	decoder := codec.GetDecoder(decoderName)
	encoder := codec.GetEncoder(encoderName)
	if decoder == nil || encoder == nil {
		if decoder == nil {
			fmt.Printf("%s: no decoder named %s\n", os.Args[0], decoderName)
		}
		if encoder == nil {
			fmt.Printf("%s: no encoder named %s\n", os.Args[0], encoderName)
		}
		os.Exit(1)
	}

	inFH := os.Stdin
	if flag.NArg() > 0 {
		var err error
		inFH, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Printf("%s: unable to open %s: %v\n", os.Args[0], flag.Arg(0), err)
			os.Exit(1)
		}
		defer inFH.Close()
	}

	outFH := os.Stdout
	if output != "" {
		var err error
		outFH, err = os.Create(output)
		if err != nil {
			fmt.Printf("%s: unable to create %s: %v\n", os.Args[0], output, err)
			os.Exit(1)
		}
		defer outFH.Close()
	}

	br := bufio.NewReader(inFH)
	bw := bufio.NewWriter(outFH)
	defer bw.Flush()

	for {
		char, err := decoder.Decode(br)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "error decoding character: %v", err)
			}
			break
		}

		encoder.Encode(bw, char)
	}
}
