package main

import (
	"fmt"
	"log"
  "io/ioutil"
  "errors"
  "debug/elf"
  "os"
  "strings"
  "strconv"

	"golang.org/x/arch/x86/x86asm"
	"github.com/nsf/termbox-go"
)

type Decompiler struct {
  path string

  file *elf.File

  Symbols []elf.Symbol
  Sections []*elf.Section

  Data []byte
}

func NewDecompiler(path string) (*Decompiler, error) {
  file, err := elf.Open(path)

  if err != nil {
    return nil, err
  }

  return &Decompiler{
    path: path, file: file}, nil
}

func (p *Decompiler) Analyze() {
  symbols, err := p.file.Symbols()

  if err == nil {
    p.Symbols = symbols
  }

  p.Sections = p.file.Sections

  data, err := ioutil.ReadFile(p.path)

  if err == nil {
    p.Data = data
  }
}

func (p *Decompiler) ShowSymbols() {
  for _, symbol := range p.Symbols {
    if len(symbol.Name) > 0 {
      fmt.Printf(
        "%24s [0x%08x, %d]\n",
        symbol.Name, symbol.Value, symbol.Size)
    }
  }

  if len(p.Symbols) == 0 {
    fmt.Println("no symbols found")
  }
}

func (p *Decompiler) ShowSections() {
  fmt.Printf("Entry Addr: 0x%08x\n", p.file.Entry)

  for _, section := range p.Sections {
    if len(section.Name) > 0 {
      fmt.Printf(
        "%24s [addr:0x%08x size:0x%08x]\n",
        section.Name, section.Addr, section.Size)
    }
  }

  if len(p.Sections) == 0 {
    fmt.Println("no sections found")
  }
}

func (p *Decompiler) GetSymbolAddress(name string) (uint64, error) {
  for _, symbol := range p.Symbols {
    if symbol.Name == name {
      return symbol.Value, nil
    }
  }

  return 0, errors.New("Symbol not found")
}

func (p *Decompiler) GetSectionAddress(name string) (uint64, error) {
  for _, section := range p.Sections {
    if section.Name == name {
      return section.Addr, nil
    }
  }

  return 0, errors.New("Section not found")
}

func (p *Decompiler) ShowAssemble(addr uint64, lines int) {
  if addr > (uint64)(len(p.Data)) {
    fmt.Printf("%24s", "Address is greater than data code\n")

    return
  }

  data := p.Data[addr:]

  for i:=0; i<lines; i++ {
    ins, err := x86asm.Decode(data, 64)

    if err != nil {
      log.Fatalln(err)
    }

    fmt.Printf("0x%08x:  %v\n", addr, ins)

    data = data[ins.Len:]

    addr = addr + (uint64)(ins.Len)
  }
}

func cli() {
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("usage: ", os.Args[0], " <binary>")

    return
  }

  dec, err := NewDecompiler(os.Args[1])

  if err != nil {
    log.Fatal(err)
  }

  var scanner string
  var words []string
  var addr uint64 = 0

  err = termbox.Init()

  if err != nil {
    panic(err)
  }

  defer termbox.Close()

  history := []string{ "" }
  historyIndex := -1

  for true {
    fmt.Print("\r>> ", scanner, "                                                                              ")

    ev := termbox.PollEvent()

    if ev.Type != termbox.EventKey {
      continue
    }

    if ev.Key == termbox.KeyArrowUp {
      if historyIndex > 0 {
        historyIndex = historyIndex - 1
      }

      scanner = history[historyIndex]
    } else if ev.Key == termbox.KeyArrowDown {
      if historyIndex < len(history) - 1 {
        historyIndex = historyIndex + 1
      }

      scanner = history[historyIndex]
    } else if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
      if len(scanner) > 0 {
        scanner = scanner[:len(scanner) - 1]
      }
    } else if ev.Key == termbox.KeyEnter {
      words = strings.Fields(scanner)
      history = append(history[:], scanner)
      history = append(history[:], "")
      historyIndex = len(history) - 1
      scanner = ""

      fmt.Printf("\n%v\n", words)

      if len(words) == 0 {
        continue
      }

      if words[0] == "help" {
        fmt.Println("quit: exit");
        fmt.Println("analyze: analyze sections, symbols and others attributes");
        fmt.Println("symbols: shows the symbols of binary");
        fmt.Println("sections: shows the sections of binary");
        fmt.Println("seek <memory/symbol/section>: seek the pointer to the refered address");
        fmt.Println("disassemble [number of instructions]: disassemble the current address");
      } else if words[0] == "quit" {
        break
      } else if words[0] == "analyze" {
        dec.Analyze()
      } else if words[0] == "symbols" {
        dec.ShowSymbols()
      } else if words[0] == "sections" {
        dec.ShowSections()
      } else if words[0] == "seek" {
        if len(words) == 2 {
          i, err := dec.GetSymbolAddress(words[1])

          if err != nil {
            i, err := dec.GetSectionAddress(words[1])

            if err != nil {
              i, err := strconv.ParseUint(words[1], 8, 64)

              if err != nil {
                i, err := strconv.ParseUint(words[1], 10, 64)

                if err != nil {
                  if strings.HasPrefix(words[1], "0x") {
                    i, err := strconv.ParseUint(words[1][2:], 16, 64)

                    if err != nil {
                      fmt.Println("address not found")
                    } else {
                      addr = i
                    }
                  } else {
                    fmt.Println("address not found")
                  }
                } else {
                  addr = i
                }
              } else {
                addr = i
              }
            } else {
              addr = i
            }
          } else {
            addr = i
          }
        }

        fmt.Printf("seek: 0x%08x\n", addr)
      } else if words[0] == "disassemble" {
        instructions := 32

        if len(words) > 1 {
          i, err := strconv.Atoi(words[1])

          if err != nil {
            instructions = i
          }

          if instructions < 0 {
            instructions = i
          }
        }

        dec.ShowAssemble(addr, instructions)
      } else {
        fmt.Println("command not found")
      }

      fmt.Print(">> ")
    } else {
      if ev.Ch != 0 {
        scanner = scanner + string(ev.Ch)
      } else {
        scanner = scanner + string(ev.Key)
      }

      history[len(history) - 1] = scanner
    }
  }
}

