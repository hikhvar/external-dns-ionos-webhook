package inwx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hikhvar/external-dns-inwx-webhook/pkg/endpoint"
	"github.com/hikhvar/external-dns-inwx-webhook/pkg/plan"
	"github.com/hikhvar/external-dns-inwx-webhook/pkg/provider"
	"github.com/nrdcg/goinwx"
	"moul.io/http2curl"
)

type ProviderConfig struct {
	BaseURL  string `env:"BASE_URL" envDefault:"https://api.domrobot.com/xmlrpc/"`
	Username string `env:"USER" envDefault:"foo"`
	Password string `env:"PASSWORD" envDefault:""`
	OTPKey   string `env:"OTP_KEY" envDefault:""`
	Debug    bool   `env:"DEBUG" envDefault:"false"`
}

type Provider struct {
	*provider.BaseProvider
	client *goinwx.Client
}

var _ http.RoundTripper = &DebugTransport{}

type DebugTransport struct {
	upstream http.RoundTripper
}

func (d DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	command, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(command)
	return d.upstream.RoundTrip(req)
}

func NewProvider(base *provider.BaseProvider, cfg ProviderConfig) (*Provider, error) {
	baseUrl, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base url: %w", err)
	}
	if cfg.Debug {
		http.DefaultTransport = &DebugTransport{upstream: http.DefaultTransport}
	}
	client := goinwx.NewClient(cfg.Username, cfg.Password, &goinwx.ClientOptions{BaseURL: baseUrl})

	_, err = client.Account.Login()
	if err != nil {
		return nil, fmt.Errorf("could not login: %w", err)
	}

	return &Provider{BaseProvider: base, client: client}, nil
}

func (p *Provider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	resp, err := p.client.Nameservers.List("")
	if err != nil {
		return nil, fmt.Errorf("could not list nameservers: %w", err)
	}

	entries := []*endpoint.Endpoint{}
	for _, domain := range resp.Domains {
		infoResp, err := p.client.Nameservers.Info(&goinwx.NameserverInfoRequest{
			Domain: domain.Domain,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get info for domain %s: %w", domain.Domain, err)
		}

		for _, record := range infoResp.Records {
			entries = append(entries, &endpoint.Endpoint{
				DNSName:       record.Name,
				Targets:       endpoint.NewTargets(record.Content),
				RecordType:    record.Type,
				SetIdentifier: "",
				RecordTTL:     endpoint.TTL(record.TTL),
				ProviderSpecific: endpoint.ProviderSpecific{
					{
						Name:  "recordID",
						Value: strconv.Itoa(record.ID),
					},
				},
			})
		}
	}

	return entries, nil
}

func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	return nil
}
