package err

import (
  "errors"
)

var (
  SymbolNotFound = errors.New("Symbol not found")
  SectionNotFound = errors.New("Section not found")
  NoStringFound = errors.New("No string found")
  NoSymbolFound = errors.New("No symbol found")
  AddressNotFound = errors.New("Address not found")
)
