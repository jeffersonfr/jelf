package state

import (
  "debug/elf"
)

type State struct {
  Path string
  File *elf.File
  Symbols []elf.Symbol
  DynamicSymbols []elf.Symbol
  Sections []*elf.Section
  Data []byte
  Strings []string
  Analyzed bool
}

