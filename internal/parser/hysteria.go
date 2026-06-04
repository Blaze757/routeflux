package parser

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/Alaxay8/routeflux/internal/domain"
)

// ParseHysteria parses a hysteria:// link.
func ParseHysteria(raw, provider string) (domain.Node, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return domain.Node{}, fmt.Errorf("parse hysteria link: %w", err)
	}

	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		return domain.Node{}, fmt.Errorf("parse hysteria port: %w", err)
	}

	query := parsed.Query()
	auth := query.Get("auth")
	if auth == "" && parsed.User != nil {
		auth = parsed.User.Username()
	}

	node := domain.Node{
		Name:       parsed.Fragment,
		Remark:     parsed.Fragment,
		Protocol:   domain.ProtocolHysteria,
		Address:    parsed.Hostname(),
		Port:       port,
		UUID:       auth,
		Password:   auth,
		ServerName: firstNonEmpty(query.Get("sni"), query.Get("peer")),
		ALPN:       splitCSV(firstNonEmpty(query.Get("alpn"), "h3")),
		Extras: map[string]string{
			"insecure":      query.Get("insecure"),
			"allowInsecure": query.Get("allowInsecure"),
			"obfs":          query.Get("obfs"),
			"obfs-password": firstNonEmpty(query.Get("obfs-password"), query.Get("obfsParam")),
		},
	}

	return normalizeNode(node, provider)
}

// ParseHysteria2 parses a hysteria2:// or hy2:// link.
func ParseHysteria2(raw, provider string) (domain.Node, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return domain.Node{}, fmt.Errorf("parse hysteria2 link: %w", err)
	}

	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		return domain.Node{}, fmt.Errorf("parse hysteria2 port: %w", err)
	}

	query := parsed.Query()
	auth := ""
	if parsed.User != nil {
		auth = parsed.User.Username()
	}

	node := domain.Node{
		Name:       parsed.Fragment,
		Remark:     parsed.Fragment,
		Protocol:   domain.ProtocolHysteria2,
		Address:    parsed.Hostname(),
		Port:       port,
		UUID:       auth,
		Password:   auth,
		ServerName: firstNonEmpty(query.Get("sni"), query.Get("peer")),
		ALPN:       splitCSV(firstNonEmpty(query.Get("alpn"), "h3")),
		Extras: map[string]string{
			"insecure":      query.Get("insecure"),
			"allowInsecure": query.Get("allowInsecure"),
			"obfs":          query.Get("obfs"),
			"obfs-password": firstNonEmpty(query.Get("obfs-password"), query.Get("obfs-param")),
		},
	}

	return normalizeNode(node, provider)
}
