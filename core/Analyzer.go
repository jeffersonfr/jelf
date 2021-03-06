package core

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "debug/elf"
  "regexp"
  "encoding/hex"
  "os"
  "os/exec"
  "syscall"
  "bytes"

  "jelf/core/state"
  "jelf/core/info"
  "jelf/core/misc"
  "jelf/core/err"
)

type Analyzer struct {
  state.State
}

func NewAnalyzer(path string) (*Analyzer, error) {
  file, err := elf.Open(path)

  if err != nil {
    return nil, err
  }

  state := state.State {
    Path: path, File: file}

  return &Analyzer{
    state}, nil
}

func (p *Analyzer) Analyze() {
  symbols, err := p.File.Symbols()

  if err == nil {
    p.Symbols = symbols
  }

  dynamicSymbols, err := p.File.DynamicSymbols()

  if err == nil {
    p.DynamicSymbols = dynamicSymbols
  }

  p.Sections = p.File.Sections

  data, err := ioutil.ReadFile(p.Path)

  if err == nil {
    p.Data = data
  }

  r, _ := regexp.Compile(`[\d\w\s,.!?@#$%^&*()-_=+{}\[\];:'"<>~?/\\]+`)
  s := r.FindAllString(string(p.Data), -1)

  for i:=0; i<len(s); i++ {
    if len(s[i]) >= 3 {
      p.Strings = append(p.Strings[:], s[i])
    }
  }

  p.Analyzed = true



  // process dyn func names
  for _, section := range p.Sections {
    if section.Name == ".dynstr" {
      data, _ := section.Data()
      texts := bytes.Split(data, []byte{0x00})

      for n, text := range texts {
        fmt.Println(n, ":", string(text))
      }
    }
  }

  for _, section := range p.Sections {
    if section.Name == ".rela.dyn" {
      data, _ := section.Data()

      fmt.Println(len(data))
      for i:=0; i<len(data); i+=24 {
        fmt.Printf("%d : %s", i/24, hex.Dump(data[i:i+24]))
      }
    }
  }
}

func (p *Analyzer) ShowStrings() {
  if p.Analyzed == false {
    fmt.Println("Call 'analyze' before")

    return
  }

  for i:=0; i<len(p.Strings); i++ {
    fmt.Println(p.Strings[i])
  }
}

func (p *Analyzer) RunProcess(args []string) {
  cmd := exec.Command(p.Path, args...)

  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  cmd.SysProcAttr = &syscall.SysProcAttr{
    Ptrace: true}

  cmd.Start()
  cmd.Wait()

  fmt.Println(cmd.Process.Pid)

  var regs syscall.PtraceRegs

  syscall.PtraceGetRegs(cmd.Process.Pid, &regs)

  fmt.Printf("RBP: %x\n", regs.Rbp);

  // Peek/Poke

  /*
  process, err := os.StartProcess(p.Path, args, nil)

  if err != nil {
		fmt.Println(err)
	}

  fmt.Println("Running: [", p.Path, "], Pid: [", process.Pid, "]")
  */

  p.Running = true
}

func (p *Analyzer) DumpBytes(address, length uint64) {
  if p.Analyzed == false {
    fmt.Println("Call 'analyze' before")

    return
  }

  if address > (uint64)(len(p.Data)) {
    fmt.Printf("Address:[0x%08x] is greater than data code:[0x%08x]\n", address, len(p.Data))

    return
  }

  if (address + length) > (uint64)(len(p.Data)) {
    length = (uint64)(len(p.Data)) - address
  }

  fmt.Printf("%s", hex.Dump(p.Data[address:address + length]))
}

func (p *Analyzer) GetSymbolAddress(name string) (uint64, error) {
  for _, symbol := range p.Symbols {
    if symbol.Name == name {
      return symbol.Value, nil
    }
  }

  return 0, err.SymbolNotFound
}

func (p *Analyzer) GetSectionAddress(name string) (uint64, error) {
  for _, section := range p.Sections {
    if section.Name == name {
      return section.Offset, nil // Addr get address from memory, Offset get "address" from file
    }
  }

  return 0, err.SectionNotFound
}

type debugInfo struct {
  Name string
  Description string
}

func (p *Analyzer) Process() {
  var scanner string
  var words []string
  var addr uint64 = p.File.Entry // INFO:: consider the memory address

  term := misc.NewTerminal()

  defer term.Release()

  info := &info.Information {
    &p.State}

  history := []string{
    "" }

  historyIndex := 0

  cmds := []debugInfo{
    {"analyze", ": process sections, symbols, ..."},
    {"symbols", ": shows symbols of binary"},
    {"sections", ": shows sections of binary"},
    {"seek", "<memory/symbol/section> : seek to the refered pointer address"},
    {"dump", "[number of bytes] : show the number of bytes starting at current address"},
    {"disassemble", "[number of instructions] : disassemble the current address"},
    {"clear", ": clear screen"},
    {"strings", ": show all strings in binary"},
    {"run", ": starts the execution of process"},
    {"quit", ": exit"}}

  expr := []string{
    "=[box] <expression> : evaluate arithmetic expression (b: binary, o: octal, h: hexadecimal)"}

  for true {
    term.ClearLine()

    fmt.Printf("0x%016x >> %s", addr, scanner)

    b := term.Read()

    if b == '\x1b' { // escape
      b = term.Read()

      if b == '\x5b' { // arrow keys
        b = term.Read()

        if b == '\x41' { // up
          if historyIndex > 0 {
            historyIndex = historyIndex - 1
            scanner = history[historyIndex]
          }
        } else if b == '\x42' { // down
          if historyIndex < len(history) - 1 {
            historyIndex = historyIndex + 1
            scanner = history[historyIndex]
          }
        } else if b == '\x43' { // right
        } else if b == '\x44' { // left
        }
      }
    } else {
      if b == '\t' { // tab
        matches := []string{
        }

        r, _ := regexp.Compile("^" + scanner)

        for _, cmd := range cmds {
          if len(r.FindString(cmd.Name)) > 0 {
            matches = append(matches[:], cmd.Name)
          }
        }

        if len(matches) == 1 {
          scanner = matches[0]
          historyIndex = len(history) - 1
          history[historyIndex] = scanner
        } else {
          fmt.Printf("\n%v\n", matches)
        }
      } else if b == '\x08' || b == '\x7f' { // backspace
        if len(scanner) > 0 {
          scanner = scanner[:len(scanner) - 1]
        }
      } else if b == '\n' { // enter
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

        words[0] = strings.Trim(words[0], " \t")

        if words[0] == "help" {
          fmt.Println(":commands:")

          for _, cmd := range cmds {
            fmt.Println("  ", cmd.Name, cmd.Description)
          }

          fmt.Println(":expressions:")

          for _, cmd := range expr {
            fmt.Println("  ", cmd)
          }
        } else if words[0] == "quit" {
          break
        } else if words[0] == "clear" {
          term.ClearScreen()
        } else if words[0] == "analyze" {
          p.Analyze()
        } else if words[0] == "strings" {
          p.ShowStrings()
        } else if words[0] == "run" {
          p.RunProcess([]string{""})
        } else if words[0] == "info" {
          info.ShowInformation()
        } else if words[0] == "dump" {
          var n uint64 = 32

          if len(words) > 1 {
            i, err := strconv.ParseUint(words[1], 10, 64)

            if err == nil {
              n = i
            }
          }

          p.DumpBytes(addr, n)
        } else if words[0] == "symbols" {
          info.ShowSymbols()
        } else if words[0] == "sections" {
          info.ShowSections()
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

            if err == nil {
              n = i
            }

            if n < 0 {
              n = i
            }
          }

          err := info.ShowAssemble(addr, n)

          if err != nil {
            fmt.Println("Invalid operators")
          }
        } else if strings.HasPrefix(words[0], "=") {
          cmd := strings.ToLower(words[0])
          expr := strings.Join(words[1:], "")

          if r, err := misc.ParseAndEval(expr); err == nil {
            if strings.HasPrefix(cmd, "=x") {
              fmt.Printf("0x%s\n", strconv.FormatUint(uint64(r), 16))
            } else if strings.HasPrefix(cmd, "=o") {
              fmt.Printf("0%s\n", strconv.FormatUint(uint64(r), 8))
            } else if strings.HasPrefix(cmd, "=b") {
              fmt.Printf("0b%s\n", strconv.FormatUint(uint64(r), 2))
            } else {
              fmt.Println(r)
            }
          } else {
            fmt.Println("invalid expression: ", err)
          }
        } else {
          fmt.Println("command not found")
        }
      } else {
        if b >= 0x20 && b < 0x7f {
          scanner = scanner + string(b)

          history[len(history) - 1] = scanner
        }
      }
    }
  }
}
