package term

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

type Term struct {
  active bool
}

func NewTerminal() *Term {
  exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
  exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

  return &Term{
    active: true}
}

func (p *Term) Release() {
  exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func (p *Term) ClearScreen() {
  fmt.Print("\033c")
}

func (p *Term) ClearLine() {
  fmt.Print("\r\033[K")
}

func (p *Term) Read() byte {
  if p.active == false {
    panic("Terminal not activated")
  }

  var b []byte = make([]byte, 1)

  os.Stdin.Read(b)

  return b[0]
}

func (p *Term) SetCursor(col, row int) {
  fmt.Println("\033[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H")
}

