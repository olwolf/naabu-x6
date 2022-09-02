package result

import (
	"sync"

	"github.com/olwolf/naabu-x6/pkg/utils"
)

type HostResult struct {
	Host  string
	IP    string
	Ports []int
}

// Result of the scan
type Result struct {
	sync.RWMutex
	IpPorts map[string]map[int]struct{}
	ips     map[string]struct{}
}

// NewResult structure
func NewResult() *Result {
	ipPorts := make(map[string]map[int]struct{})
	ips := make(map[string]struct{})
	return &Result{IpPorts: ipPorts, ips: ips}
}

// AddPort to a specific ip
func (r *Result) GetIPs() chan string {
	r.Lock()

	out := make(chan string)

	go func() {
		defer close(out)
		defer r.Unlock()

		for ip := range r.ips {
			out <- ip
		}
	}()

	return out
}

func (r *Result) HasIPS() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.ips) > 0
}

// AddPort to a specific ip
func (r *Result) GetIPsPorts() chan *HostResult {
	r.RLock()

	out := make(chan *HostResult)

	go func() {
		defer close(out)
		defer r.RUnlock()

		for ip, ports := range r.IpPorts {
			out <- &HostResult{IP: ip, Ports: utils.MapKeysToSliceInt(ports)}
		}
	}()

	return out
}

func (r *Result) HasIPsPorts() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.IpPorts) > 0
}

// AddPort to a specific ip
func (r *Result) AddPort(k string, v int) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.IpPorts[k]; !ok {
		r.IpPorts[k] = make(map[int]struct{})
	}

	r.IpPorts[k][v] = struct{}{}
	r.ips[k] = struct{}{}
}

// SetPorts for a specific ip
func (r *Result) SetPorts(ip string, ports []int) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.IpPorts[ip]; !ok {
		r.IpPorts[ip] = make(map[int]struct{})
	}

	for _, port := range ports {
		r.IpPorts[ip][port] = struct{}{}
	}
	r.ips[ip] = struct{}{}
}

// IPHasPort checks if an ip has a specific port
func (r *Result) IPHasPort(k string, v int) bool {
	r.RLock()
	defer r.RUnlock()

	vv, hasports := r.IpPorts[k]
	if !hasports {
		return false
	}
	_, hasport := vv[v]

	return hasport
}

// AddIp adds an ip to the results
func (r *Result) AddIp(ip string) {
	r.Lock()
	defer r.Unlock()

	r.ips[ip] = struct{}{}
}

// HasIP checks if an ip has been seen
func (r *Result) HasIP(ip string) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.ips[ip]
	return ok
}

func (r *Result) IsEmpty() bool {
	return r.Len() == 0
}

func (r *Result) Len() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.ips)
}
