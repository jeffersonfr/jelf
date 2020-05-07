package main

import (
  "fmt"
  "os"
  "log"

  jelf "jelf/core"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println("usage: ", os.Args[0], " <binary>")

    return
  }

  analyzer, err := jelf.NewAnalyzer(os.Args[1])

  if err != nil {
    log.Fatal(err)
  }

  analyzer.Process()
}

