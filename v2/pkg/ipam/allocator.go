package ipam

import (
	"fmt"
	"net"

	"github.com/cybozu-go/coil/v2/pkg/util"
	"github.com/willf/bitset"
)

type allocator struct {
	ipv4  *net.IPNet
	ipv6  *net.IPNet
	usage *bitset.BitSet
}

func newAllocator(ipv4, ipv6 *string) (a allocator) {
	if ipv4 != nil {
		ip, n, _ := net.ParseCIDR(*ipv4)
		if ip.To4() == nil {
			panic("ipv4 must be an IPv4 subnet")
		}
		a.ipv4 = n
		ones, bits := n.Mask.Size()
		a.usage = bitset.New(uint(1) << (bits - ones))
	}
	if ipv6 != nil {
		_, n, _ := net.ParseCIDR(*ipv6)
		a.ipv6 = n
		if a.usage == nil {
			ones, bits := n.Mask.Size()
			a.usage = bitset.New(uint(1) << (bits - ones))
		}
	}
	return
}

func (a allocator) isFull() bool {
	return a.usage.All()
}

func (a allocator) isEmpty() bool {
	return a.usage.None()
}

func (a allocator) register(ipv4, ipv6 net.IP) (uint, bool) {
	if a.ipv4 != nil && a.ipv4.Contains(ipv4) {
		offset := util.IPDiff(a.ipv4.IP, ipv4)
		if offset < 0 {
			panic(fmt.Sprintf("ip: %v, base: %v, offset: %v", ipv4, a.ipv4.IP, offset))
		}
		a.usage.Set(uint(offset))
		return uint(offset), true
	}
	if a.ipv6 != nil && a.ipv6.Contains(ipv6) {
		offset := util.IPDiff(a.ipv6.IP, ipv6)
		if offset < 0 {
			panic(fmt.Sprintf("ip: %v, base: %v, offset: %v", ipv6, a.ipv6.IP, offset))
		}
		a.usage.Set(uint(offset))
		return uint(offset), true
	}
	return 0, false
}

func (a allocator) allocate() (ipv4, ipv6 net.IP, idx uint, ok bool) {
	idx, ok = a.usage.NextClear(0)
	if !ok {
		return nil, nil, 0, false
	}

	if a.ipv4 != nil {
		ipv4 = util.IPAdd(a.ipv4.IP, int64(idx))
	}
	if a.ipv6 != nil {
		ipv6 = util.IPAdd(a.ipv6.IP, int64(idx))
	}
	a.usage.Set(idx)
	return
}

func (a allocator) free(idx uint) {
	a.usage.Clear(idx)
}
