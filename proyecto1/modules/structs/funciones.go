package structs

import (
	"fmt"
	"os"
	"text/tabwriter"

	seccomp "github.com/seccomp/libseccomp-golang"
)

type SyscallCounter []int

const maxSyscalls = 303

func (s SyscallCounter) Init() SyscallCounter {
	s = make(SyscallCounter, maxSyscalls)
	return s
}

func (s SyscallCounter) Inc(syscallID uint64) error {
	if syscallID > maxSyscalls {
		return fmt.Errorf("invalid syscall ID (%x)", syscallID)
	}

	s[syscallID]++
	return nil
}

func (s SyscallCounter) Print() []Strace {
	var list []Strace
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for k, v := range s {
		if v > 0 {
			name, _ := seccomp.ScmpSyscall(k).GetName()
			fmt.Fprintf(w, "%d\t%s\n", v, name)
			var strace Strace
			strace.Name = name
			strace.Recurrence = v
			list = append(list, strace)
		}
	}
	w.Flush()
	return list
}

func (s SyscallCounter) GetName(syscallID uint64) string {
	name, _ := seccomp.ScmpSyscall(syscallID).GetName()
	return name
}