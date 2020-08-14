package tcpproxy

import (
	"bufio"
	"log"
)

// AddressResolver defines the interface for resolving proxy target address
type AddressResolver interface {
	Resolve() (string, error)
}

// AddResolverRoute appends a targetResolver matching route to the ipPort listener,
// directing any connection to dest.
//
// This is generally used as either the only rule (for simple TCP
// proxies), or as the final fallback rule for an ipPort.
//
// The ipPort is any valid net.Listen TCP address.
func (p *Proxy) AddResolverRoute(ipPort string, r AddressResolver) {
	p.addRoute(ipPort, resolverMatch{resolver: r})
}

type resolverMatch struct {
	target   Target
	resolver AddressResolver
}

func (tr resolverMatch) match(*bufio.Reader) (Target, string) {
	// resolve the proxy address
	addr, err := tr.resolver.Resolve()
	if err != nil {
		// log
		log.Printf("tcpproxy: could not resolve address for target %v", err)
		return nil, ""
	}

	tr.target = &DialProxy{Addr: addr}
	return tr.target, ""
}
