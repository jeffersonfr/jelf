package state

import (
  "debug/elf"
)

type State struct {
  Path string
  File *elf.File
  Symbols []elf.Symbol
  Sections []*elf.Section
  Data []byte
  Analyzed bool
}

