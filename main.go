package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/debug/elf"
	elf2 "github.com/sitano/goelf/elf"
	"reflect"
	"strings"
)

var filename = flag.StringP("filename", "f", "", "Path to the elf binary")
var all = flag.BoolP("all", "a", false, "Print all available information")
var header = flag.Bool("header", false, "Print header")
var sections = flag.Bool("sections", false, "Print sections")
var symbols = flag.Bool("symbols", false, "Print symbols")
var imports = flag.Bool("imports", false, "Print imports")
var progs = flag.Bool("progs", false, "Print progs")
var notes = flag.Bool("notes", false, "Print notes")
var note_prstatus = flag.Bool("note_prstatus", false, "Print prstatus note")
var note_prpsinfo = flag.Bool("note_prpsinfo", false, "Print prpsinfo note")

func main() {
	flag.Parse()

	if *filename == "" {
		fmt.Fprintln(os.Stderr, "Filename is required")
		os.Exit(1)
	}

	p, err := New(*filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file", err)
		os.Exit(1)
	}

	if *all || *header {
		p.PrintHeader()
	}

	if *all || *sections {
		p.PrintSections()
	}

	if *all || *progs {
		p.PrintProgs()
	}

	if *all || *imports {
		p.PrintImports()
	}

	if *all || *notes {
		p.PrintNotes()
	}

	if *all || *note_prstatus {
		p.PrintPRStatus()
	}

	if *all || *note_prpsinfo {
		p.PrintPRPSInfo()
	}

	if *all || *symbols {
		p.PrintSymbols()
	}
}

func (p *Process) PrintHeader() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Class", "Data", "Version",
		"OSABI", "ABIVersion", "ByteOrder",
		"Type", "Machine", "Entry",
	})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Append([]string{
		fmt.Sprintf("%v", p.efd.Class),
		fmt.Sprintf("%v", p.efd.Data),
		fmt.Sprintf("%v", p.efd.Version),
		fmt.Sprintf("%v", p.efd.OSABI),
		fmt.Sprintf("0x%x", p.efd.ABIVersion),
		fmt.Sprintf("%v", p.efd.ByteOrder),
		fmt.Sprintf("%v", p.efd.Type),
		fmt.Sprintf("%v", p.efd.Machine),
		fmt.Sprintf("0x%x", p.efd.Entry),

	})
	table.Render()
	fmt.Println()
}

func (p *Process) PrintSections() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Id", "Sections", "Type", "Flags",
		"Addr", "Offset", "Size",
		"Link", "Info", "Addralign",
		"Entsize",
	})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for id, s := range p.efd.Sections {
		table.Append([]string{
			fmt.Sprintf("%d", id),
			s.Name,
			fmt.Sprintf("%v", s.Type),
			fmt.Sprintf("%v", s.Flags),
			fmt.Sprintf("0x%x", s.Addr),
			fmt.Sprintf("0x%x", s.Offset),
			fmt.Sprintf("0x%x", s.Size),
			fmt.Sprintf("0x%x", s.Link),
			fmt.Sprintf("0x%x", s.Info),
			fmt.Sprintf("%d", s.Addralign),
			fmt.Sprintf("%d", s.Entsize),
		})
	}
	table.Render()
	fmt.Println()
}

func (p *Process) PrintProgs() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Progs", "Flags", "Off",
		"Vaddr", "Paddr", "Filesz",
		"Memsz", "Align",
	})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, p := range p.efd.Progs {
		table.Append([]string{
			fmt.Sprintf("%v", p.Type),
			fmt.Sprintf("%v", p.Flags),
			fmt.Sprintf("0x%x", p.Off),
			fmt.Sprintf("0x%x", p.Vaddr),
			fmt.Sprintf("0x%x", p.Paddr),
			fmt.Sprintf("0x%x", p.Filesz),
			fmt.Sprintf("0x%x", p.Memsz),
			fmt.Sprintf("0x%x", p.Align),
		})
	}
	table.Render()
	fmt.Println()
}

func (p *Process) PrintSymbols() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Sym", "Info", "Other", "Section", "Offset", "Size",
	})
	table.SetBorder(false)
	table.SetAutoWrapText(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	sym, err := p.efd.Symbols()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading .symtab", err)
	}

	for _, s := range sym {
		section := ""
		if s.Section == elf.SHN_UNDEF || s.Section >= elf.SHN_LORESERVE {
			section = fmt.Sprintf("%v", s.Section)
		} else {
			section = fmt.Sprintf("%d", int(s.Section))
		}

		table.Append([]string{
			s.Name,
			fmt.Sprintf("0x%x", s.Info),
			fmt.Sprintf("0x%x", s.Other),
			fmt.Sprintf("%v", section),
			fmt.Sprintf("0x%x", s.Value),
			fmt.Sprintf("%d", s.Size),
		})
	}

	table.Render()
	fmt.Println()
}

func (p *Process) PrintImports() {
	isym, err := p.efd.ImportedSymbols()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading .dynsym", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Imported Symbols", "Version", "Library",
	})
	table.SetBorder(false)
	table.SetAutoWrapText(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, s := range isym {
		table.Append([]string{
			s.Name,
			s.Version,
			s.Library,
		})
	}

	table.Render()
	fmt.Println()

	libs, err := p.efd.ImportedLibraries()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading .needed", err)
	}

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Library"})
	table.SetBorder(false)
	table.SetAutoWrapText(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, l := range libs {
		table.Append([]string{l})
	}

	table.Render()
	fmt.Println()
}

func (p *Process) PrintNotes() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Note", "Type", "Size", "Data",
	})
	table.SetBorder(false)
	table.SetAutoWrapText(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	s := p.efd.SectionByType(elf.SHT_NOTE)
	if s == nil {
		fmt.Fprintln(os.Stderr, "Error searching for note (7) section")
	}

	notes, err := elf2.ReadNotes(s, p.efd.ByteOrder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading", s.Name, ":", err)
	}

	for _, n := range notes {
		typeString := fmt.Sprintf("%v", n.Type)

		data := "..."

		if n.Name == "Go" && n.Type == elf2.NT_GO_BUILD {
			typeString = "NT_GOBUILDID"
			data = string(n.Data)
		}

		table.Append([]string{
			n.Name,
			typeString,
			fmt.Sprintf("0x%x", len(n.Data)),
			data,
		})
	}

	table.Render()
	fmt.Println()
}

func (p *Process) PrintPRStatus() {
	s := p.efd.SectionByType(elf.SHT_NOTE)
	if s == nil {
		fmt.Fprintln(os.Stderr, "Error searching for note (7) section")
	}

	note, err := elf2.ReadNoteByType(s, p.efd.ByteOrder, elf2.NT_PRSTATUS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading NT_PRSTATUS:", err)
		return
	}

	prs, err := elf2.ReadPRStatus(note, p.efd.ByteOrder, p.efd.Class)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading NT_PRSTATUS:", err)
		return
	}

	PrintStruct(*prs, 1)
	fmt.Println()

	PrintStruct(elf2.GetUserRegs(prs.Regs), 1)
	fmt.Println()
}

func (p *Process) PrintPRPSInfo() {
	s := p.efd.SectionByType(elf.SHT_NOTE)
	if s == nil {
		fmt.Fprintln(os.Stderr, "Error searching for note (7) section")
	}

	note, err := elf2.ReadNoteByType(s, p.efd.ByteOrder, elf2.NT_PRPSINFO)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading NT_PRPSINFO:", err)
		return
	}

	prps, err := elf2.ReadPRPSInfo(note, p.efd.ByteOrder, p.efd.Class)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading NT_PRPSINFO:", err)
		return
	}

	PrintStruct(*prps, 1)
	fmt.Println()
}

func PrintStruct(s interface{}, indent int) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	ind := strings.Repeat("\t", indent)

	fmt.Println("Struct", t.Name())
	for i := 0; i < v.NumField(); i ++ {
		tf := t.Field(i)
		vf := v.Field(i)
		fn := tf.Name
		ft := tf.Type

		if ft.Kind() == reflect.Struct {
			fmt.Printf("%s%s = ", ind, fn)
			PrintStruct(v.Field(i).Interface(), indent + 1)
		} else {
			fmt.Printf("%s%s = 0x%x (%v)\n", ind, fn, vf.Interface(), vf.Interface())
		}
	}
}