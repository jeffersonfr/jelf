package main

import (
  "fmt"
  "os"
  "log"
  "jelf"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println("usage: ", os.Args[0], " <binary>")

    return
  }

  debugger, err := jelf.NewDebugger(os.Args[1])

  if err != nil {
    log.Fatal(err)
  }

  debugger.Process()
}

