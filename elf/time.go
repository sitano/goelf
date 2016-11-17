package elf

type KernelTime int64
type KernelSUSeconds int64

// http://lxr.free-electrons.com/source/include/uapi/linux/time.h#L15
type TimeVal struct {
	Sec KernelTime       /* seconds */
	USec KernelSUSeconds /* microseconds */
}
