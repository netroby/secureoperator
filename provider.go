package secureoperator

import (
	"github.com/miekg/dns"
	"net"
)

type DNSQuestion struct {
	Name string `json:"name,omitempty"`
	Type int32  `json:"type,omitempty"`
}

type DNSRR struct {
	Name string `json:"name,omitempty"`
	Type int32  `json:"type,omitempty"`
	TTL  int32  `json:"TTL,omitempty"`
	Data string `json:"data,omitempty"`
}

func (r DNSRR) RR() dns.RR {
	var rr dns.RR

	// Build an RR header
	rrhdr := dns.RR_Header{
		Name:     r.Name,
		Rrtype:   uint16(r.Type),
		Class:    dns.ClassINET,
		Ttl:      uint32(r.TTL),
		Rdlength: uint16(len(r.Data)),
	}

	constructor, ok := dns.TypeToRR[uint16(r.Type)]
	if ok {
		// Construct a new RR
		rr = constructor()
		*(rr.Header()) = rrhdr
		switch v := rr.(type) {
		case *dns.A:
			v.A = net.ParseIP(r.Data)
		case *dns.AAAA:
			v.AAAA = net.ParseIP(r.Data)
		case *dns.CNAME:
			v.Target = r.Data
		case *dns.MX:
			v.Mx = r.Data
		}
	} else {
		rr = dns.RR(&dns.RFC3597{
			Hdr:   rrhdr,
			Rdata: r.Data,
		})
	}
	return rr
}

type DNSResponse struct {
	Question           []DNSQuestion
	Answer             []DNSRR
	Authority          []DNSRR
	Extra              []DNSRR
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	AuthenticatedData  bool
	CheckingDisabled   bool
	ResponseCode       int
}

type Provider interface {
	Query(DNSQuestion) (*DNSResponse, error)
}