package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		flagString string
		flagInt    int
		flagBool   bool
	)

	flag.StringVar(&flagString, "stringvar", "default", "an example string var")
	flag.IntVar(&flagInt, "intval", 42, "an example int var")
	flag.BoolVar(&flagBool, "boolval", false, "an example bool var")

	flag.Parse()

	fmt.Println("String var:", flagString)
	fmt.Println("Int var:", flagInt)
	fmt.Println("Boolean var:", flagBool)

	// os.Args содержит все аргументы командной строки,
	// включая имя исполняемого файла
	args := os.Args
	fmt.Println(args)

	// Определение флагов
	verbose := flag.Bool("v", false, "verbose output")
	port := flag.Int("port", 8080, "server port")

	// Парсинг флагов
	flag.Parse()

	// Аргументы командной строки (не флаги)
	args = flag.Args()

	fmt.Printf("Verbose: %v\n", *verbose)
	fmt.Printf("Port: %d\n", *port)
	fmt.Printf("Arguments: %v\n", args)
	fmt.Printf("All command line args (including flags): %v\n", os.Args)
}
