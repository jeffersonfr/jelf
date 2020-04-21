package jelf

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "errors"
  "debug/elf"
  "regexp"
  "encoding/hex"

  "jelf/state"
  "jelf/info"

	"github.com/nsf/termbox-go"
)

type Debugger struct {
  state.State
}

func NewDebugger(path string) (*Debugger, error) {
  file, err := elf.Open(path)

  if err != nil {
    return nil, err
  }

  state := state.State {
    Path: path, File: file}

  return &Debugger{
    state}, nil
}

func (p *Debugger) Analyze() {
  symbols, err := p.File.Symbols()

  if err == nil {
    p.Symbols = symbols
  }

  p.Sections = p.File.Sections

  data, err := ioutil.ReadFile(p.Path)

  if err == nil {
    p.Data = data
  }

  r, _ := regexp.Compile(`[\d\w\s,.!?@#$%^&*()-_=+{}\[\];:'"<>~?/\\]+`)
  s := r.FindAllString(string(p.Data), -1)

  for i:=0; i<len(s); i++ {
    if len(s[i]) > 3 {
      p.Strings = append(p.Strings[:], s[i])
    }
  }

  p.Analyzed = true
}

func (p *Debugger) ShowStrings() {
  for i:=0; i<len(p.Strings); i++ {
    fmt.Println(p.Strings[i])
  }
}

func (p *Debugger) DumpBytes(address, length uint64) {
  if p.Analyzed == false {
    fmt.Println("Call 'analyze' before")

    return
  }

  if address > (uint64)(len(p.Data)) {
    fmt.Printf("Address:[0x%08x] is greater than data code:[0x%08x]\n", address, len(p.Data))

    return
  }

  fmt.Printf("%s", hex.Dump(p.Data[address:address + length]))
}

func (p *Debugger) GetSymbolAddress(name string) (uint64, error) {
  for _, symbol := range p.Symbols {
    if symbol.Name == name {
      return symbol.Value, nil
    }
  }

  return 0, errors.New("Symbol not found")
}

func (p *Debugger) GetSectionAddress(name string) (uint64, error) {
  for _, section := range p.Sections {
    if section.Name == name {
      return section.Offset, nil // Addr get address from memory, Offset get "address" from file
    }
  }

  return 0, errors.New("Section not found")
}

func (p *Debugger) Process() {
  var scanner string
  var words []string
  var addr uint64 = p.File.Entry

  err := termbox.Init()

  if err != nil {
    panic(err)
  }

  defer termbox.Close()

  var infoPtr info.IInformation = &info.Information {
    &p.State}

  history := []string{
    "" }

  historyIndex := 0

  cmds := []string{
    "analyze",
    "symbols",
    "sections",
    "seek",
    "dump",
    "disassemble",
    "clear",
    "strings",
    "quit"}

  for true {
    width, _ := termbox.Size()

    fmt.Printf(fmt.Sprintf("\r0x%%016x >> %%-%ds", width - 22 - len(scanner) - 1), addr, scanner + "_")

    ev := termbox.PollEvent()

    if ev.Type != termbox.EventKey {
      continue
    }

    if ev.Key == termbox.KeyArrowUp {
      if historyIndex > 0 {
        historyIndex = historyIndex - 1
        scanner = history[historyIndex]
      }
    } else if ev.Key == termbox.KeyArrowDown {
      if historyIndex < len(history) - 1 {
        historyIndex = historyIndex + 1
        scanner = history[historyIndex]
      }
    } else if ev.Key == termbox.KeyTab {
      matches := []string{
      }

      r, _ := regexp.Compile("^" + scanner)

      for _, cmd := range cmds {
        if len(r.FindString(cmd)) > 0 {
          matches = append(matches[:], cmd)
        }
      }

      if len(matches) == 1 {
        scanner = matches[0]
        historyIndex = len(history) - 1
        history[historyIndex] = scanner
      } else {
        fmt.Printf("\n%v\n", matches)
      }
    } else if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
      if len(scanner) > 0 {
        scanner = scanner[:len(scanner) - 1]
      }
    } else if ev.Key == termbox.KeyEnter {
      words = strings.Fields(scanner)

      if len(scanner) == 0 {
        continue
      }

      history = append(history[:], "")
      historyIndex = len(history) - 1
      scanner = ""

      fmt.Printf("\n")

      if len(words) == 0 {
        continue
      }

      if words[0] == "help" {
        fmt.Println(`
.: jELF :: Binary Parser :.

  analyze: analyze sections, symbols, ...

  symbols: shows symbols of binary

  sections: shows sections of binary

  seek <memory/symbol/section>: seek to the refered pointer address

  dump [number of bytes]: show the number of bytes starting at current address

  disassemble [number of instructions]: disassemble the current address

  clear: clear screen

  quit: exit
        `)
      } else if words[0] == "quit" {
        break
      } else if words[0] == "clear" {
        termbox.Clear(0x09, 0x00) // fg: white, bg: black
        termbox.Sync()
      } else if words[0] == "analyze" {
        p.Analyze()
      } else if words[0] == "strings" {
        p.ShowStrings()
      } else if words[0] == "info" {
        infoPtr.ShowInformation()
      } else if words[0] == "dump" {
        var n uint64 = 32

        if len(words) > 1 {
          i, err := strconv.ParseUint(words[1], 10, 64)

          if err == nil {
            n = i
          }
        }

        p.DumpBytes(addr, n)
      } else if words[0] == "level" {
        i, err := strconv.ParseUint(words[1], 10, 64)

        if err != nil {
          fmt.Println("Invalid command")

          continue
        }

        if i == 1 {
          infoPtr = &info.Information {
            &p.State }
        } else if i == 2 {
          infoPtr = &info.AdvancedInformation {
            &info.Information {
              &p.State } }
        } else {
          fmt.Println("Invalid level")

          continue
        }
      } else if words[0] == "symbols" {
        infoPtr.ShowSymbols()
      } else if words[0] == "sections" {
        infoPtr.ShowSections()
      } else if words[0] == "seek" {
        if len(words) == 2 {
          i, err := p.GetSymbolAddress(words[1])

          if err != nil {
            i, err := p.GetSectionAddress(words[1])

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
        }

        // fmt.Printf("seek: 0x%08x\n", addr)
      } else if words[0] == "disassemble" {
        n := 32

        if len(words) > 1 {
          i, err := strconv.Atoi(words[1])

          if err != nil {
            n = i
          }

          if n < 0 {
            n = i
          }
        }

        err := infoPtr.ShowAssemble(addr, n)

        if err != nil {
          fmt.Println("Invalid operators")
        }
      } else {
        fmt.Println("command not found")
      }
    } else {
      var ch rune = 0

      if ev.Ch != 0 {
        ch = ev.Ch
      } else {
        ch = (rune)(ev.Key)
      }

      if ch >= 0x20 && ch < 0x7f {
        scanner = scanner + string(ch)

        history[len(history) - 1] = scanner
      }
    }
  }
}
