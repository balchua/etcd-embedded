package etcd

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/balchua/etcd-embedded/pkg/util"
	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	DefaultPeerPort   int    = 2380
	DefaultClientPort int    = 2379
	DefaultScheme     string = "http"
)

type EtcdConfig struct {
	// Human-readable name for this member.
	Name string `yaml:"name,omitempty"`

	// Path to the data directory.
	DataDir string `yaml:"data-dir,omitempty"`

	// Path to the dedicated wal directory.
	WalDir string `yaml:"wal-dir,omitempty"`

	// List of comma separated URLs to listen on for peer traffic.
	ListenPeerUrls string `yaml:"listen-peer-urls,omitempty"`

	// List of comma separated URLs to listen on for client traffic.
	ListenClientUrls string `yaml:"listen-client-urls,omitempty"`

	// List of this member's peer URLs to advertise to the rest of the cluster.
	// The URLs needed to be a comma-separated list.
	InitialAdvertisePeerUrls string `yaml:"initial-advertise-peer-urls,omitempty"`

	// List of this member's client URLs to advertise to the public.
	// The URLs needed to be a comma-separated list.
	AdvertiseClientUrls string `yaml:"advertise-client-urls,omitempty"`

	// Initial cluster configuration for bootstrapping.
	InitialCluster string `yaml:"initial-cluster,omitempty"`

	// Initial cluster token for the etcd cluster during bootstrap.
	InitialClusterToken string `yaml:"initial-cluster-token,omitempty"`

	// Initial cluster state ('new' or 'existing').
	InitialClusterState string `yaml:"initial-cluster-state,omitempty"`

	ClientTransportSecurity ClientTransportSecurityInfo `yaml:"client-transport-security,omitempty"`

	PeerTransportSecurity PeerTransportSecurityInfo `yaml:"peer-transport-security,omitempty"`

	// Enable debug-level logging for etcd.
	Debug bool `yaml:"debug,omitempty"`

	// Maximum number of snapshot files to retain (0 is unlimited).
	MaxSnapshots int `yaml:"max-snapshots,omitempty"`

	// Maximum number of wal files to retain (0 is unlimited).
	MaxWals int `yaml:"max-wals,omitempty"`
}

type ClientTransportSecurityInfo struct {

	// Path to the client server TLS cert file.
	CertFile string `yaml:"cert-file,omitempty"`

	// Path to the client server TLS key file.
	KeyFile string `yaml:"key-file,omitempty"`

	// Enable client cert authentication.
	ClientCertAuth bool `yaml:"client-cert-auth,omitempty"`

	// Path to the client server TLS trusted CA cert file.
	TrustedCaFile string `yaml:"trusted-ca-file,omitempty"`

	// Client TLS using generated certificates
	AutoTls bool `yaml:"auto-tls,omitempty"`
}

type PeerTransportSecurityInfo struct {
	// Path to the peer server TLS cert file.
	CertFile string `yaml:"cert-file,omitempty"`

	// Path to the peer client server TLS key file.
	KeyFile string `yaml:"key-file,omitempty"`

	// Enable peer client cert authentication.
	ClientCertAuth bool `yaml:"client-cert-auth,omitempty"`

	// Path to the peer client server TLS trusted CA cert file.
	TrustedCaFile string `yaml:"trusted-ca-file,omitempty"`

	// Peer TLS using generated certificates
	AutoTls bool `yaml:"auto-tls,omitempty"`
}

func LoadEtcdConfig(configPath string) *EtcdConfig {
	var lg *zap.Logger

	lg, err := zap.NewProduction()

	if err != nil {
		lg.Warn("Unable to read etcd config file.", zap.Error(err))
	}

	etcdConfig := &EtcdConfig{}

	yamlFile, ioerr := ioutil.ReadFile(configPath)

	if ioerr != nil {
		lg.Fatal("Unable to read etcd config file.", zap.Error(ioerr))
	}

	yamlErr := yaml.Unmarshal(yamlFile, etcdConfig)

	if yamlErr != nil {
		lg.Fatal("Unable to read etcd config file.", zap.Error(yamlErr))
	}

	etcdConfig.init()
	return etcdConfig
}

func (e *EtcdConfig) init() {
	if e.Name == "" {
		e.Name, _ = os.Hostname()
	}

	defaultIp := getDefaultIP()

	if e.ListenPeerUrls == "" {
		e.ListenPeerUrls = DefaultScheme + "://" + defaultIp + ":" + fmt.Sprint(DefaultPeerPort)
	}

	if e.ListenClientUrls == "" {
		e.ListenClientUrls = DefaultScheme + "://" + defaultIp + ":" + fmt.Sprint(DefaultClientPort)
	}

	if e.InitialCluster == "" {
		e.InitialCluster = e.Name + "=" + e.ListenPeerUrls
	}

	if e.InitialAdvertisePeerUrls == "" {
		e.InitialAdvertisePeerUrls = DefaultScheme + "://" + defaultIp + ":" + fmt.Sprint(DefaultPeerPort)
	}

	if e.AdvertiseClientUrls == "" {
		e.AdvertiseClientUrls = DefaultScheme + "://" + defaultIp + ":" + fmt.Sprint(DefaultClientPort)
	}

	if e.InitialClusterState == "" {
		e.InitialClusterState = "new"
	}

}

func (e *EtcdConfig) ToEmbedEtcdConfig() *embed.Config {
	var embedConfig = embed.NewConfig()

	embedConfig.Name = e.Name

	embedConfig.Dir = e.DataDir
	embedConfig.WalDir = e.WalDir

	embedConfig.LPUrls = e.toUrl(e.ListenPeerUrls)

	embedConfig.LCUrls = e.toUrl(e.ListenClientUrls)

	embedConfig.InitialCluster = e.InitialCluster
	embedConfig.InitialClusterToken = e.InitialClusterToken

	embedConfig.APUrls = e.toUrl(e.InitialAdvertisePeerUrls)

	embedConfig.ACUrls = e.toUrl(e.AdvertiseClientUrls)

	embedConfig.ClusterState = e.InitialClusterState

	embedConfig.StrictReconfigCheck = false
	embedConfig.ClientAutoTLS = e.ClientTransportSecurity.AutoTls
	embedConfig.ClientTLSInfo.CertFile = e.ClientTransportSecurity.CertFile
	embedConfig.ClientTLSInfo.KeyFile = e.ClientTransportSecurity.KeyFile
	embedConfig.ClientTLSInfo.ClientCertAuth = e.ClientTransportSecurity.ClientCertAuth
	embedConfig.ClientTLSInfo.TrustedCAFile = e.ClientTransportSecurity.TrustedCaFile

	embedConfig.PeerAutoTLS = e.PeerTransportSecurity.AutoTls
	embedConfig.PeerTLSInfo.CertFile = e.PeerTransportSecurity.CertFile
	embedConfig.PeerTLSInfo.KeyFile = e.PeerTransportSecurity.KeyFile
	embedConfig.PeerTLSInfo.ClientCertAuth = e.PeerTransportSecurity.ClientCertAuth
	embedConfig.PeerTLSInfo.TrustedCAFile = e.PeerTransportSecurity.TrustedCaFile
	embedConfig.MaxSnapFiles = 0
	embedConfig.MaxWalFiles = 0
	embedConfig.Logger = "zap"
	if e.Debug {
		embedConfig.LogLevel = "debug"
	} else {
		embedConfig.LogLevel = "info"
	}

	embedConfig.LogOutputs = []string{"stderr"}
	return embedConfig
}

func (e *EtcdConfig) toUrl(commaSeparatedUrl string) []url.URL {
	lg, _ := zap.NewProduction()

	list := strings.Split(commaSeparatedUrl, ",")
	var urls []url.URL
	//urls := make([]url.URL, len(list))

	for _, item := range list {
		url, err := url.Parse(item)
		if err != nil {
			lg.Fatal("Unable to convert to URL", zap.String("urlInString", item))
		}
		urls = append(urls, *url)
	}
	return urls
}

func (e *EtcdConfig) ToFile(config string) {
	lg, _ := zap.NewProduction()
	data, err := yaml.Marshal(e)

	if err != nil {
		lg.Fatal("", zap.Error(err))
	}

	err2 := ioutil.WriteFile(config, data, 0)

	if err2 != nil {
		lg.Fatal("", zap.Error(err2))
	}

}

func getDefaultIP() string {
	ip, _ := util.GetDefaultIPV4()
	return ip
}
