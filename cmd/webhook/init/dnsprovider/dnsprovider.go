package dnsprovider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v8"
	"github.com/hikhvar/external-dns-inwx-webhook/cmd/webhook/init/configuration"
	"github.com/hikhvar/external-dns-inwx-webhook/internal/inwx"
	"github.com/hikhvar/external-dns-inwx-webhook/pkg/endpoint"
	"github.com/hikhvar/external-dns-inwx-webhook/pkg/provider"
	log "github.com/sirupsen/logrus"
)

func Init(config configuration.Config) (provider.Provider, error) {
	var domainFilter endpoint.DomainFilter
	createMsg := "Creating inwx provider with "

	if config.RegexDomainFilter != "" {
		createMsg += fmt.Sprintf("Regexp domain filter: '%s', ", config.RegexDomainFilter)
		if config.RegexDomainExclusion != "" {
			createMsg += fmt.Sprintf("with exclusion: '%s', ", config.RegexDomainExclusion)
		}
		domainFilter = endpoint.NewRegexDomainFilter(
			regexp.MustCompile(config.RegexDomainFilter),
			regexp.MustCompile(config.RegexDomainExclusion),
		)
	} else {
		if len(config.DomainFilter) > 0 {
			createMsg += fmt.Sprintf("zoneNode filter: '%s', ", strings.Join(config.DomainFilter, ","))
		}
		if len(config.ExcludeDomains) > 0 {
			createMsg += fmt.Sprintf("Exclude domain filter: '%s', ", strings.Join(config.ExcludeDomains, ","))
		}
		domainFilter = endpoint.NewDomainFilterWithExclusions(config.DomainFilter, config.ExcludeDomains)
	}

	createMsg = strings.TrimSuffix(createMsg, ", ")
	if strings.HasSuffix(createMsg, "with ") {
		createMsg += "no kind of domain filters"
	}
	log.Info(createMsg)
	/*ionosConfig := ionos.Configuration{}
	if err := env.Parse(&ionosConfig); err != nil {
		return nil, fmt.Errorf("reading ionos ionosConfig failed: %v", err)
	}

	ionosProvider := createProvider(baseProvider, &ionosConfig)
	return ionosProvider, nil

	*/

	cfg := inwx.ProviderConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse inwx config: %w", err)
	}
	baseProvider := provider.NewBaseProvider(domainFilter)
	return inwx.NewProvider(baseProvider, cfg)
}
