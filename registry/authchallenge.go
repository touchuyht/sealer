// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	bearerRegex = regexp.MustCompile(
		`^\s*Bearer\s+(.*)$`)
	basicRegex = regexp.MustCompile(`^\s*Basic\s+.*$`)

	// ErrBasicAuth indicates that the repository requires basic rather than token authentication.
	ErrBasicAuth = errors.New("basic auth required")
)

func parseAuthHeader(header http.Header) (*authService, error) {
	ch, err := parseChallenge(header.Get("www-authenticate"))
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func parseChallenge(challengeHeader string) (*authService, error) {
	if basicRegex.MatchString(challengeHeader) {
		return nil, ErrBasicAuth
	}

	match := bearerRegex.FindAllStringSubmatch(challengeHeader, -1)
	if d := len(match); d != 1 {
		return nil, fmt.Errorf("malformed auth challenge header: '%s', %d", challengeHeader, d)
	}
	parts := strings.SplitN(strings.TrimSpace(match[0][1]), ",", 3)

	var realm, service string
	var scope []string
	for _, s := range parts {
		p := strings.SplitN(s, "=", 2)
		if len(p) != 2 {
			return nil, fmt.Errorf("malformed auth challenge header: '%s'", challengeHeader)
		}
		key := p[0]
		value := strings.TrimSuffix(strings.TrimPrefix(p[1], `"`), `"`)
		switch key {
		case "realm":
			realm = value
		case "service":
			service = value
		case "scope":
			scope = strings.Fields(value)
		default:
			return nil, fmt.Errorf("unknown field in challenge header %s: %v", key, challengeHeader)
		}
	}
	parsedRealm, err := url.Parse(realm)
	if err != nil {
		return nil, err
	}

	a := &authService{
		Realm:   parsedRealm,
		Service: service,
		Scope:   scope,
	}

	return a, nil
}