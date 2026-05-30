package parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Alaxay8/routeflux/internal/domain"
)

// ParseSocks parses SOCKS5 share links in socks:// formats.
func ParseSocks(raw, provider string) (domain.Node, error) {
	trimmed := strings.TrimSpace(raw)
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return domain.Node{}, fmt.Errorf("parse socks link: %w", err)
	}

	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		return domain.Node{}, fmt.Errorf("parse socks port: %w", err)
	}

	var username, password string
	if parsed.User != nil {
		u := parsed.User.Username()
		p, hasPass := parsed.User.Password()

		// SOCKS URLs often base64-encode the credentials (username:password)
		// e.g. YWxheGF5OmFsYXhheQ== -> alaxay:alaxay
		if !hasPass && looksLikeBase64Payload(u) {
			decoded, err := base64.StdEncoding.DecodeString(u)
			if err == nil {
				parts := strings.SplitN(string(decoded), ":", 2)
				if len(parts) == 2 {
					username = parts[0]
					password = parts[1]
				} else {
					username = string(decoded)
				}
			} else {
				username = u
			}
		} else {
			username = u
			password = p
		}
	}

	node := domain.Node{
		Name:     parsed.Fragment,
		Remark:   parsed.Fragment,
		Protocol: domain.ProtocolSocks,
		Address:  parsed.Hostname(),
		Port:     port,
		UUID:     username,
		Password: password,
	}

	return normalizeNode(node, provider)
}
