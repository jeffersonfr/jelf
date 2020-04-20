package info

type IInformation interface {
  ShowInformation()
  ShowSymbols()
  ShowSections()
  ShowAssemble(addr uint64, lines int) error
}

