package info

import (
	"fmt"
  "debug/elf"
  "encoding/hex"
  "strings"
  "regexp"
  "strconv"

  "jelf/state"
  "jelf/err"

	"golang.org/x/arch/x86/x86asm"
)

type Information struct {
  *state.State
}

func (p *Information) ShowInformation() {
  fmt.Print("Class:[", p.File.Class.String(), "]: ")

  if (p.File.Class == elf.ELFCLASS32) {
    fmt.Println("32 bits")
  } else if (p.File.Class == elf.ELFCLASS64) {
    fmt.Println("64 bits")
  } else {
    fmt.Println("Unknown")
  }

  fmt.Print("Endianess:[", p.File.Data.String(), "]: ")

  if (p.File.Data == elf.ELFDATA2LSB) {
    fmt.Println("Little Endian")
  } else if (p.File.Data == elf.ELFDATA2MSB) {
    fmt.Println("Big Endian")
  } else {
    fmt.Println("Unknown")
  }

  fmt.Print("System ABI:[", p.File.OSABI.String(), "]: ")

  if (p.File.OSABI == elf.ELFOSABI_HPUX) {
    fmt.Println("HP-UX operating system")
  } else if (p.File.OSABI == elf.ELFOSABI_NETBSD) {
    fmt.Println("NetBSD")
  } else if (p.File.OSABI == elf.ELFOSABI_LINUX) {
    fmt.Println("GNU/Linux")
  } else if (p.File.OSABI == elf.ELFOSABI_HURD) {
    fmt.Println("GNU/Hurd")
  } else if (p.File.OSABI == elf.ELFOSABI_86OPEN) {
    fmt.Println("86Open common IA32 ABI")
  } else if (p.File.OSABI == elf.ELFOSABI_SOLARIS) {
    fmt.Println("Solaris")
  } else if (p.File.OSABI == elf.ELFOSABI_AIX) {
    fmt.Println("AIX")
  } else if (p.File.OSABI == elf.ELFOSABI_IRIX) {
    fmt.Println("IRIX")
  } else if (p.File.OSABI == elf.ELFOSABI_FREEBSD) {
    fmt.Println("FreeBSD")
  } else if (p.File.OSABI == elf.ELFOSABI_TRU64) {
    fmt.Println("TRU64 UNIX")
  } else if (p.File.OSABI == elf.ELFOSABI_MODESTO) {
    fmt.Println("Novell Modesto")
  } else if (p.File.OSABI == elf.ELFOSABI_OPENBSD) {
    fmt.Println("OpenBSD")
  } else if (p.File.OSABI == elf.ELFOSABI_OPENVMS) {
    fmt.Println("Open VMS")
  } else if (p.File.OSABI == elf.ELFOSABI_NSK) {
    fmt.Println("HP Non-Stop Kernel")
  } else if (p.File.OSABI == elf.ELFOSABI_AROS) {
    fmt.Println("Amiga Research OS")
  } else if (p.File.OSABI == elf.ELFOSABI_FENIXOS) {
    fmt.Println("The FenixOS highly scalable multi-core OS")
  } else if (p.File.OSABI == elf.ELFOSABI_CLOUDABI) {
    fmt.Println("Nuxi CloudABI")
  } else if (p.File.OSABI == elf.ELFOSABI_ARM) {
    fmt.Println("ARM")
  } else if (p.File.OSABI == elf.ELFOSABI_STANDALONE) {
    fmt.Println("Standalone (embedded) application")
  } else {
    fmt.Println("Unknown")
  }

  fmt.Print("Type:[", p.File.Type.String(), "]: ")

  if (p.File.Type == elf.ET_REL) {
    fmt.Println("Relocatable")
  } else if (p.File.Type == elf.ET_EXEC) {
    fmt.Println("Executable")
  } else if (p.File.Type == elf.ET_DYN) {
    fmt.Println("Shared object")
  } else if (p.File.Type == elf.ET_CORE) {
    fmt.Println("Core file")
  } else if (p.File.Type == elf.ET_LOOS) {
    fmt.Println("First operating system specific")
  } else if (p.File.Type == elf.ET_HIOS) {
    fmt.Println("Last operating system-specific")
  } else if (p.File.Type == elf.ET_LOPROC) {
    fmt.Println("First processor-specific")
  } else if (p.File.Type == elf.ET_HIPROC) {
    fmt.Println("Last processor-specific")
  } else {
    fmt.Println("Unknown")
  }

  fmt.Print("Machine:[", p.File.Machine.String(), "]: ")

  if (p.File.Machine == elf.EM_M32) {
    fmt.Println("AT&T WE32100")
  } else if (p.File.Machine == elf.EM_SPARC) {
    fmt.Println("Sun SPARC")
  } else if (p.File.Machine == elf.EM_386) {
    fmt.Println("Intel i386")
  } else if (p.File.Machine == elf.EM_68K) {
    fmt.Println("Motorola 68000")
  } else if (p.File.Machine == elf.EM_88K) {
    fmt.Println("Motorola 88000")
  } else if (p.File.Machine == elf.EM_860) {
    fmt.Println("Intel i860")
  } else if (p.File.Machine == elf.EM_MIPS) {
    fmt.Println("MIPS R3000 Big-Endian only")
  } else if (p.File.Machine == elf.EM_S370) {
    fmt.Println("IBM System/370")
  } else if (p.File.Machine == elf.EM_MIPS_RS3_LE) {
    fmt.Println("MIPS R3000 Little-Endian")
  } else if (p.File.Machine == elf.EM_PARISC) {
    fmt.Println("HP PA-RISC")
  } else if (p.File.Machine == elf.EM_VPP500) {
    fmt.Println("Fujitsu VPP500")
  } else if (p.File.Machine == elf.EM_SPARC32PLUS) {
    fmt.Println("SPARC v8plus")
  } else if (p.File.Machine == elf.EM_960) {
    fmt.Println("Intel 80960")
  } else if (p.File.Machine == elf.EM_PPC) {
    fmt.Println("PowerPC 32-bit")
  } else if (p.File.Machine == elf.EM_PPC64) {
    fmt.Println("PowerPC 64-bit")
  } else if (p.File.Machine == elf.EM_S390) {
    fmt.Println("IBM System/390")
  } else if (p.File.Machine == elf.EM_V800) {
    fmt.Println("NEC V800")
  } else if (p.File.Machine == elf.EM_FR20) {
    fmt.Println("Fujitsu FR20")
  } else if (p.File.Machine == elf.EM_RH32) {
    fmt.Println("TRW RH-32")
  } else if (p.File.Machine == elf.EM_RCE) {
    fmt.Println("Motorola RCE")
  } else if (p.File.Machine == elf.EM_ARM) {
    fmt.Println("ARM")
  } else if (p.File.Machine == elf.EM_SH) {
    fmt.Println("Hitachi SH")
  } else if (p.File.Machine == elf.EM_SPARCV9) {
    fmt.Println("SPARC v9 64-bit")
  } else if (p.File.Machine == elf.EM_TRICORE) {
    fmt.Println("Siemens TriCore embedded processor")
  } else if (p.File.Machine == elf.EM_ARC) {
    fmt.Println("Argonaut RISC Core")
  } else if (p.File.Machine == elf.EM_H8_300) {
    fmt.Println("Hitachi H8/300")
  } else if (p.File.Machine == elf.EM_H8_300H) {
    fmt.Println("Hitachi H8/300H")
  } else if (p.File.Machine == elf.EM_H8S) {
    fmt.Println("Hitachi H8S")
  } else if (p.File.Machine == elf.EM_H8_500) {
    fmt.Println("Hitachi H8/500")
  } else if (p.File.Machine == elf.EM_IA_64) {
    fmt.Println("Intel IA-64 Processor")
  } else if (p.File.Machine == elf.EM_MIPS_X) {
    fmt.Println("Stanford MIPS-X")
  } else if (p.File.Machine == elf.EM_COLDFIRE) {
    fmt.Println("Motorola ColdFire")
  } else if (p.File.Machine == elf.EM_68HC12) {
    fmt.Println("Motorola M68HC12")
  } else if (p.File.Machine == elf.EM_MMA) {
    fmt.Println("Fujitsu MMA")
  } else if (p.File.Machine == elf.EM_PCP) {
    fmt.Println("Siemens PCP")
  } else if (p.File.Machine == elf.EM_NCPU) {
    fmt.Println("Sony nCPU")
  } else if (p.File.Machine == elf.EM_NDR1) {
    fmt.Println("Denso NDR1 microprocessor")
  } else if (p.File.Machine == elf.EM_STARCORE) {
    fmt.Println("Motorola Star*Core processor")
  } else if (p.File.Machine == elf.EM_ME16) {
    fmt.Println("Toyota ME16 processor")
  } else if (p.File.Machine == elf.EM_ST100) {
    fmt.Println("STMicroelectronics ST100 processor")
  } else if (p.File.Machine == elf.EM_TINYJ) {
    fmt.Println("Advanced Logic Corp. TinyJ processor")
  } else if (p.File.Machine == elf.EM_X86_64) {
    fmt.Println("Advanced Micro Devices x86-64")
  } else if (p.File.Machine == elf.EM_PDSP) {
    fmt.Println("Sony DSP Processor")
  } else if (p.File.Machine == elf.EM_PDP10) {
    fmt.Println("Digital Equipment Corp. PDP-10")
  } else if (p.File.Machine == elf.EM_PDP11) {
    fmt.Println("Digital Equipment Corp. PDP-11")
  } else if (p.File.Machine == elf.EM_FX66) {
    fmt.Println("Siemens FX66 microcontroller")
  } else if (p.File.Machine == elf.EM_ST9PLUS) {
    fmt.Println("STMicroelectronics ST9+ 8/16 bit microcontroller")
  } else if (p.File.Machine == elf.EM_ST7) {
    fmt.Println("STMicroelectronics ST7 8-bit microcontroller")
  } else if (p.File.Machine == elf.EM_68HC16) {
    fmt.Println("Motorola MC68HC16 Microcontroller")
  } else if (p.File.Machine == elf.EM_68HC11) {
    fmt.Println("Motorola MC68HC11 Microcontroller")
  } else if (p.File.Machine == elf.EM_68HC08) {
    fmt.Println("Motorola MC68HC08 Microcontroller")
  } else if (p.File.Machine == elf.EM_68HC05) {
    fmt.Println("Motorola MC68HC05 Microcontroller")
  } else if (p.File.Machine == elf.EM_SVX) {
    fmt.Println("Silicon Graphics SVx")
  } else if (p.File.Machine == elf.EM_ST19) {
    fmt.Println("STMicroelectronics ST19 8-bit microcontroller")
  } else if (p.File.Machine == elf.EM_VAX) {
    fmt.Println("Digital VAX")
  } else if (p.File.Machine == elf.EM_CRIS) {
    fmt.Println("Axis Communications 32-bit embedded processor")
  } else if (p.File.Machine == elf.EM_JAVELIN) {
    fmt.Println("Infineon Technologies 32-bit embedded processor")
  } else if (p.File.Machine == elf.EM_FIREPATH) {
    fmt.Println("Element 14 64-bit DSP Processor")
  } else if (p.File.Machine == elf.EM_ZSP) {
    fmt.Println("LSI Logic 16-bit DSP Processor")
  } else if (p.File.Machine == elf.EM_MMIX) {
    fmt.Println("Donald Knuth's educational 64-bit processor")
  } else if (p.File.Machine == elf.EM_HUANY) {
    fmt.Println("Harvard University machine-independent object files")
  } else if (p.File.Machine == elf.EM_PRISM) {
    fmt.Println("SiTera Prism")
  } else if (p.File.Machine == elf.EM_AVR) {
    fmt.Println("Atmel AVR 8-bit microcontroller")
  } else if (p.File.Machine == elf.EM_FR30) {
    fmt.Println("Fujitsu FR30")
  } else if (p.File.Machine == elf.EM_D10V) {
    fmt.Println("Mitsubishi D10V")
  } else if (p.File.Machine == elf.EM_D30V) {
    fmt.Println("Mitsubishi D30V")
  } else if (p.File.Machine == elf.EM_V850) {
    fmt.Println("NEC v850")
  } else if (p.File.Machine == elf.EM_M32R) {
    fmt.Println("Mitsubishi M32R")
  } else if (p.File.Machine == elf.EM_MN10300) {
    fmt.Println("Matsushita MN10300")
  } else if (p.File.Machine == elf.EM_MN10200) {
    fmt.Println("Matsushita MN10200")
  } else if (p.File.Machine == elf.EM_PJ) {
    fmt.Println("picoJava")
  } else if (p.File.Machine == elf.EM_OPENRISC) {
    fmt.Println("OpenRISC 32-bit embedded processor")
  } else if (p.File.Machine == elf.EM_ARC_COMPACT) {
    fmt.Println("ARC International ARCompact processor (old spelling/synonym: EM_ARC_A5)")
  } else if (p.File.Machine == elf.EM_XTENSA) {
    fmt.Println("Tensilica Xtensa Architecture")
  } else if (p.File.Machine == elf.EM_VIDEOCORE) {
    fmt.Println("Alphamosaic VideoCore processor")
  } else if (p.File.Machine == elf.EM_TMM_GPP) {
    fmt.Println("Thompson Multimedia General Purpose Processor")
  } else if (p.File.Machine == elf.EM_NS32K) {
    fmt.Println("National Semiconductor 32000 series")
  } else if (p.File.Machine == elf.EM_TPC) {
    fmt.Println("Tenor Network TPC processor")
  } else if (p.File.Machine == elf.EM_SNP1K) {
    fmt.Println("Trebia SNP 1000 processor")
  } else if (p.File.Machine == elf.EM_ST200) {
    fmt.Println("ST200 microcontroller")
  } else if (p.File.Machine == elf.EM_IP2K) {
    fmt.Println("Ubicom IP2xxx microcontroller family")
  } else if (p.File.Machine == elf.EM_MAX) {
    fmt.Println("MAX Processor")
  } else if (p.File.Machine == elf.EM_CR) {
    fmt.Println("National Semiconductor CompactRISC microprocessor")
  } else if (p.File.Machine == elf.EM_F2MC16) {
    fmt.Println("Fujitsu F2MC16")
  } else if (p.File.Machine == elf.EM_MSP430) {
    fmt.Println("Texas Instruments embedded microcontroller msp430")
  } else if (p.File.Machine == elf.EM_BLACKFIN) {
    fmt.Println("Analog Devices Blackfin (DSP) processor")
  } else if (p.File.Machine == elf.EM_SE_C33) {
    fmt.Println("S1C33 Family of Seiko Epson processors")
  } else if (p.File.Machine == elf.EM_SEP) {
    fmt.Println("Sharp embedded microprocessor")
  } else if (p.File.Machine == elf.EM_ARCA) {
    fmt.Println("Arca RISC Microprocessor")
  } else if (p.File.Machine == elf.EM_UNICORE) {
    fmt.Println("Microprocessor series from PKU-Unity Ltd. and MPRC of Peking University")
  } else if (p.File.Machine == elf.EM_EXCESS) {
    fmt.Println("eXcess: 16/32/64-bit configurable embedded CPU")
  } else if (p.File.Machine == elf.EM_DXP) {
    fmt.Println("Icera Semiconductor Inc. Deep Execution Processor")
  } else if (p.File.Machine == elf.EM_ALTERA_NIOS2) {
    fmt.Println("Altera Nios II soft-core processor")
  } else if (p.File.Machine == elf.EM_CRX) {
    fmt.Println("National Semiconductor CompactRISC CRX microprocessor")
  } else if (p.File.Machine == elf.EM_XGATE) {
    fmt.Println("Motorola XGATE embedded processor")
  } else if (p.File.Machine == elf.EM_C166) {
    fmt.Println("Infineon C16x/XC16x processor")
  } else if (p.File.Machine == elf.EM_M16C) {
    fmt.Println("Renesas M16C series microprocessors")
  } else if (p.File.Machine == elf.EM_DSPIC30F) {
    fmt.Println("Microchip Technology dsPIC30F Digital Signal Controller")
  } else if (p.File.Machine == elf.EM_CE) {
    fmt.Println("Freescale Communication Engine RISC core")
  } else if (p.File.Machine == elf.EM_M32C) {
    fmt.Println("Renesas M32C series microprocessors")
  } else if (p.File.Machine == elf.EM_TSK3000) {
    fmt.Println("Altium TSK3000 core")
  } else if (p.File.Machine == elf.EM_RS08) {
    fmt.Println("Freescale RS08 embedded processor")
  } else if (p.File.Machine == elf.EM_SHARC) {
    fmt.Println("Analog Devices SHARC family of 32-bit DSP processors")
  } else if (p.File.Machine == elf.EM_ECOG2) {
    fmt.Println("Cyan Technology eCOG2 microprocessor")
  } else if (p.File.Machine == elf.EM_SCORE7) {
    fmt.Println("Sunplus S+core7 RISC processor")
  } else if (p.File.Machine == elf.EM_DSP24) {
    fmt.Println("24-bit DSP Processor")
  } else if (p.File.Machine == elf.EM_VIDEOCORE3) {
    fmt.Println("Broadcom VideoCore III processor")
  } else if (p.File.Machine == elf.EM_LATTICEMICO32) {
    fmt.Println("RISC processor for Lattice FPGA architecture")
  } else if (p.File.Machine == elf.EM_SE_C17) {
    fmt.Println("Seiko Epson C17 family")
  } else if (p.File.Machine == elf.EM_TI_C6000) {
    fmt.Println("The Texas Instruments TMS320C6000 DSP family")
  } else if (p.File.Machine == elf.EM_TI_C2000) {
    fmt.Println("The Texas Instruments TMS320C2000 DSP family")
  } else if (p.File.Machine == elf.EM_TI_C5500) {
    fmt.Println("The Texas Instruments TMS320C55x DSP family")
  } else if (p.File.Machine == elf.EM_TI_ARP32) {
    fmt.Println("Texas Instruments Application Specific RISC Processor, 32bit fetch")
  } else if (p.File.Machine == elf.EM_TI_PRU) {
    fmt.Println("Texas Instruments Programmable Realtime Unit")
  } else if (p.File.Machine == elf.EM_MMDSP_PLUS) {
    fmt.Println("STMicroelectronics 64bit VLIW Data Signal Processor")
  } else if (p.File.Machine == elf.EM_CYPRESS_M8C) {
    fmt.Println("Cypress M8C microprocessor")
  } else if (p.File.Machine == elf.EM_R32C) {
    fmt.Println("Renesas R32C series microprocessors")
  } else if (p.File.Machine == elf.EM_TRIMEDIA) {
    fmt.Println("NXP Semiconductors TriMedia architecture family")
  } else if (p.File.Machine == elf.EM_QDSP6) {
    fmt.Println("QUALCOMM DSP6 Processor")
  } else if (p.File.Machine == elf.EM_8051) {
    fmt.Println("Intel 8051 and variants")
  } else if (p.File.Machine == elf.EM_STXP7X) {
    fmt.Println("STMicroelectronics STxP7x family of configurable and extensible RISC processors")
  } else if (p.File.Machine == elf.EM_NDS32) {
    fmt.Println("Andes Technology compact code size embedded RISC processor family")
  } else if (p.File.Machine == elf.EM_ECOG1) {
    fmt.Println("Cyan Technology eCOG1X family")
  } else if (p.File.Machine == elf.EM_ECOG1X) {
    fmt.Println("Cyan Technology eCOG1X family")
  } else if (p.File.Machine == elf.EM_MAXQ30) {
    fmt.Println("Dallas Semiconductor MAXQ30 Core Micro-controllers")
  } else if (p.File.Machine == elf.EM_XIMO16) {
    fmt.Println("16-bit DSP Processor")
  } else if (p.File.Machine == elf.EM_MANIK) {
    fmt.Println("M2000 Reconfigurable RISC Microprocessor")
  } else if (p.File.Machine == elf.EM_CRAYNV2) {
    fmt.Println("Cray Inc. NV2 vector architecture")
  } else if (p.File.Machine == elf.EM_RX) {
    fmt.Println("Renesas RX family")
  } else if (p.File.Machine == elf.EM_METAG) {
    fmt.Println("Imagination Technologies META processor architecture")
  } else if (p.File.Machine == elf.EM_MCST_ELBRUS) {
    fmt.Println("MCST Elbrus general purpose hardware architecture")
  } else if (p.File.Machine == elf.EM_ECOG16) {
    fmt.Println("Cyan Technology eCOG16 family")
  } else if (p.File.Machine == elf.EM_CR16) {
    fmt.Println("National Semiconductor CompactRISC CR16 16-bit microprocessor")
  } else if (p.File.Machine == elf.EM_ETPU) {
    fmt.Println("Freescale Extended Time Processing Unit")
  } else if (p.File.Machine == elf.EM_SLE9X) {
    fmt.Println("Infineon Technologies SLE9X core")
  } else if (p.File.Machine == elf.EM_L10M) {
    fmt.Println("Intel L10M")
  } else if (p.File.Machine == elf.EM_K10M) {
    fmt.Println("Intel K10M")
  } else if (p.File.Machine == elf.EM_AARCH64) {
    fmt.Println("ARM 64-bit Architecture (AArch64)")
  } else if (p.File.Machine == elf.EM_AVR32) {
    fmt.Println("Atmel Corporation 32-bit microprocessor family")
  } else if (p.File.Machine == elf.EM_STM8) {
    fmt.Println("STMicroeletronics STM8 8-bit microcontroller")
  } else if (p.File.Machine == elf.EM_TILE64) {
    fmt.Println("Tilera TILE64 multicore architecture family")
  } else if (p.File.Machine == elf.EM_TILEPRO) {
    fmt.Println("Tilera TILEPro multicore architecture family")
  } else if (p.File.Machine == elf.EM_MICROBLAZE) {
    fmt.Println("Xilinx MicroBlaze 32-bit RISC soft processor core")
  } else if (p.File.Machine == elf.EM_CUDA) {
    fmt.Println("NVIDIA CUDA architecture")
  } else if (p.File.Machine == elf.EM_TILEGX) {
    fmt.Println("Tilera TILE-Gx multicore architecture family")
  } else if (p.File.Machine == elf.EM_CLOUDSHIELD) {
    fmt.Println("CloudShield architecture family")
  } else if (p.File.Machine == elf.EM_COREA_1ST) {
    fmt.Println("KIPO-KAIST Core-A 1st generation processor family")
  } else if (p.File.Machine == elf.EM_COREA_2ND) {
    fmt.Println("KIPO-KAIST Core-A 2nd generation processor family")
  } else if (p.File.Machine == elf.EM_ARC_COMPACT2) {
    fmt.Println("Synopsys ARCompact V2")
  } else if (p.File.Machine == elf.EM_OPEN8) {
    fmt.Println("Open8 8-bit RISC soft processor core")
  } else if (p.File.Machine == elf.EM_RL78) {
    fmt.Println("Renesas RL78 family")
  } else if (p.File.Machine == elf.EM_VIDEOCORE5) {
    fmt.Println("Broadcom VideoCore V processor")
  } else if (p.File.Machine == elf.EM_78KOR) {
    fmt.Println("Renesas 78KOR family")
  } else if (p.File.Machine == elf.EM_56800EX) {
    fmt.Println("Freescale 56800EX Digital Signal Controller (DSC)")
  } else if (p.File.Machine == elf.EM_BA1) {
    fmt.Println("Beyond BA1 CPU architecture")
  } else if (p.File.Machine == elf.EM_BA2) {
    fmt.Println("Beyond BA2 CPU architecture")
  } else if (p.File.Machine == elf.EM_XCORE) {
    fmt.Println("XMOS xCORE processor family")
  } else if (p.File.Machine == elf.EM_MCHP_PIC) {
    fmt.Println("Microchip 8-bit PIC(r) family")
  } else if (p.File.Machine == elf.EM_INTEL205) {
    fmt.Println("Reserved by Intel")
  } else if (p.File.Machine == elf.EM_INTEL206) {
    fmt.Println("Reserved by Intel")
  } else if (p.File.Machine == elf.EM_INTEL207) {
    fmt.Println("Reserved by Intel")
  } else if (p.File.Machine == elf.EM_INTEL208) {
    fmt.Println("Reserved by Intel")
  } else if (p.File.Machine == elf.EM_INTEL209) {
    fmt.Println("Reserved by Intel")
  } else if (p.File.Machine == elf.EM_KM32) {
    fmt.Println("KM211 KM32 32-bit processor")
  } else if (p.File.Machine == elf.EM_KMX32) {
    fmt.Println("KM211 KMX32 32-bit processor")
  } else if (p.File.Machine == elf.EM_KMX16) {
    fmt.Println("KM211 KMX16 16-bit processor")
  } else if (p.File.Machine == elf.EM_KMX8) {
    fmt.Println("KM211 KMX8 8-bit processor")
  } else if (p.File.Machine == elf.EM_KVARC) {
    fmt.Println("KM211 KVARC processor")
  } else if (p.File.Machine == elf.EM_CDP) {
    fmt.Println("Paneve CDP architecture family")
  } else if (p.File.Machine == elf.EM_COGE) {
    fmt.Println("Cognitive Smart Memory Processor")
  } else if (p.File.Machine == elf.EM_COOL) {
    fmt.Println("Bluechip Systems CoolEngine")
  } else if (p.File.Machine == elf.EM_NORC) {
    fmt.Println("Nanoradio Optimized RISC")
  } else if (p.File.Machine == elf.EM_CSR_KALIMBA) {
    fmt.Println("CSR Kalimba architecture family")
  } else if (p.File.Machine == elf.EM_Z80) {
    fmt.Println("Zilog Z80")
  } else if (p.File.Machine == elf.EM_VISIUM) {
    fmt.Println("Controls and Data Services VISIUMcore processor")
  } else if (p.File.Machine == elf.EM_FT32) {
    fmt.Println("FTDI Chip FT32 high performance 32-bit RISC architecture")
  } else if (p.File.Machine == elf.EM_MOXIE) {
    fmt.Println("Moxie processor family")
  } else if (p.File.Machine == elf.EM_AMDGPU) {
    fmt.Println("AMD GPU architecture")
  } else if (p.File.Machine == elf.EM_RISCV) {
    fmt.Println("RISC-V")
  } else if (p.File.Machine == elf.EM_LANAI) {
    fmt.Println("Lanai 32-bit processor")
  } else if (p.File.Machine == elf.EM_BPF) {
    fmt.Println("Linux BPF – in-kernel virtual machine")
  } else {
    fmt.Println("Unknown")
  }
}

func (p *Information) ShowSymbols() {
  for _, symbol := range p.Symbols {
    if len(symbol.Name) > 0 {
      fmt.Printf(
        "%32s [0x%08x, %d]\n",
        symbol.Name, symbol.Value, symbol.Size)
    }
  }

  for _, symbol := range p.DynamicSymbols {
    if len(symbol.Name) > 0 {
      fmt.Printf(
        "%32s [0x%08x, %d]\n",
        symbol.Name, symbol.Value, symbol.Size)
    }
  }

  if len(p.Symbols) == 0 {
    fmt.Println("no symbols found")
  }
}

func (p *Information) ShowSections() {
  fmt.Printf("Entry Addr: 0x%08x\n", p.File.Entry)

  for _, section := range p.Sections {
    if len(section.Name) > 0 {
      fmt.Printf(
        "%32s [addr:0x%08x off:0x%08x size:0x%08x]\n",
        section.Name, section.Addr, section.Offset, section.Size)
    }
  }

  if len(p.Sections) == 0 {
    fmt.Println("no sections found")
  }
}

func (p *Information) GetStringFromAddress(addr uint64) (string, error) {
  r, _ := regexp.Compile(`[\d\w\s,.!?@#$%^&*()-_=+{}\[\];:'"<>~?/\\]+`)
  s := r.FindAllString(string(p.Data[addr:]), -1)

  if len(s) > 0 {
    return "\"" + s[0] + "\"", nil
  }

  return "", err.NoStringFound
}

func (p *Information) GetSymbolFromAddress(addr uint64) (string, error) {
  for _, symbol := range p.Symbols {
    if len(symbol.Name) > 0 {
      if symbol.Value == addr {
        return symbol.Name, nil
      }
    }
  }

  return "", err.NoSymbolFound
}

func (p *Information) GetAddressFromLea(addr uint64, code string) (uint64, error) {
  r, _ := regexp.Compile(`lea .*, \[([-+]?.*)\]`)
  s := r.FindAllSubmatch([]byte(code), -1)

  if len(s) > 0 {
    offsetStr := string(s[0][1])

    offset, err := strconv.ParseInt(strings.Replace(offsetStr, "0x", "", -1), 16, 64)

    if err == nil {
      return uint64(int64(addr) + offset), nil
    }
  }

  return 0, err.AddressNotFound
}

func (p *Information) GetAddressFromCall(addr uint64, code string) (uint64, error) {
  r, _ := regexp.Compile(`call \[([-+]?.*)\]`)
  s := r.FindAllSubmatch([]byte(code), -1)

  if len(s) > 0 {
    offsetStr := string(s[0][1])

    offset, err := strconv.ParseInt(strings.Replace(offsetStr, "0x", "", -1), 16, 64)

    if err == nil {
      return uint64(int64(addr) + offset), nil
    }
  }

  return 0, err.AddressNotFound
}

func (p *Information) GetAddressContent(addr uint64, code string) string {
  caddr, err := p.GetAddressFromLea(addr, code)

  if err != nil {
    caddr, err = p.GetAddressFromCall(addr, code)

    if err != nil {
      return ""
    }
  }

  str, err := p.GetSymbolFromAddress(caddr)

  if err == nil {
    return fmt.Sprintf("[0x%08x] %s", caddr, str)
  }

  str, err = p.GetStringFromAddress(caddr)

  if err == nil {
    return fmt.Sprintf("[0x%08x] %s", caddr, str)
  }

  return fmt.Sprintf("[0x%08x]", caddr)
}

func (p *Information) ShowAssemble(addr uint64, lines int) error {
  if p.Analyzed == false {
    fmt.Println("Call 'analyze' before")

    return nil
  }

  if addr > (uint64)(len(p.Data)) {
    fmt.Printf("Address:[0x%08x] is greater than data code:[0x%08x]\n", addr, len(p.Data))

    return nil
  }

  data := p.Data[addr:]

  for i:=0; i<lines; i++ {
    ins, err := x86asm.Decode(data, 32)

    if err != nil {
      // return err
    }

    instruction := strings.ToLower(ins.String())
    content := p.GetAddressContent(uint64(int(addr) + ins.Len), instruction)

    if len(content) > 0 {
      content = "; " + content
    }

    fmt.Printf("0x%08x:  %32v\t%-32v%s\n", addr, instruction, hex.EncodeToString(data[0:ins.Len]), content)

    data = data[ins.Len:]

    addr = addr + (uint64)(ins.Len)
  }

  return nil
}

