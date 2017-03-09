package config

import (
	"fmt"
	"io"
	"os"

	"github.com/tochka/flag"
)

var defaultConfig Config
var configPath string

func init() {
	flag.StringVar(&defaultConfig.Admin.Login, "admin-login", "admin", "User name for login to admin panel")
	flag.StringVar(&defaultConfig.Admin.Password, "admin-password", "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918", "SHA256 hash of admin password")
	flag.BoolVar(&defaultConfig.Admin.Enabled, "admin-enabled", false, "Enebled admin panel")
	flag.StringVar(&defaultConfig.Auth.Mode, "auth-mode", "no", "Authentication mode")
	flag.StringVar(&defaultConfig.Auth.Params.Host, "auth-address", "", "Remote authorization service address")
	flag.StringVar(&defaultConfig.Auth.TokenType, "auth-token-type", "VIRGIL", "Authorization type")

	flag.IntVar(&defaultConfig.Cards.Cache.Duration, "cache-duration", 3600, "Cache duration")
	flag.IntVar(&defaultConfig.Cards.Cache.SizeMb, "cache-size", 1024, "cache will not allocate more memory than this limit, value in MB. if value is reached then the oldest entries can be overridden for the new ones  0 value means no size limit")

	flag.StringVar(&defaultConfig.Cards.Mode, "mode", "cache", "VirgilD service mode")
	flag.StringVar(&defaultConfig.Cards.Remote.Authority.CardID, "authority-card-id", "3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853", "Authority card id")
	flag.StringVar(&defaultConfig.Cards.Remote.Authority.PublicKey, "authority-pubkey", "MCowBQYDK2VwAyEAYR501kV1tUne2uOdkw4kErRRbJrc2Syaz5V1fuG+rVs=", "Authority public key")
	flag.IntVar(&defaultConfig.Cards.Remote.Cache, "cache", 3600, "Caching duration for global cards (in secondes)")
	flag.StringVar(&defaultConfig.Cards.Remote.Services.Cards, "cards-service", "https://cards.virgilsecurity.com", "Address of Cards service")
	flag.StringVar(&defaultConfig.Cards.Remote.Services.CardsRO, "cards-ro-service", "https://cards-ro.virgilsecurity.com", "Address of Read only cards service")
	flag.StringVar(&defaultConfig.Cards.Remote.Services.Identity, "identity-service", "https://identity.virgilsecurity.com", "Address of identity service")
	flag.StringVar(&defaultConfig.Cards.Remote.Services.VRA, "ra-service", "https://ra.virgilsecurity.com", "Address of registration authority service")
	flag.StringVar(&defaultConfig.Cards.Remote.Token, "remote-token", "", "Token for get access to Virgil cloud")
	flag.StringVar(&defaultConfig.Cards.Signer.CardID, "vd-card-id", "", "VirgilD card id")
	flag.StringVar(&defaultConfig.Cards.Signer.PrivateKey, "vd-key", "", "VirgilD private key")
	flag.StringVar(&defaultConfig.Cards.Signer.PrivateKeyPassword, "vd-key-password", "", "Password for Virgild private key")
	flag.StringVar(&defaultConfig.Cards.VRA.CardID, "ra-card-id", "", "Registration Authority card id")
	flag.StringVar(&defaultConfig.Cards.VRA.PublicKey, "ra-pubkey", "", "Registration Authority public key")

	flag.StringVar(&defaultConfig.DB, "db", "sqlite3:virgild.db", "Database connection string {driver}:{connection}. Supported drivers: sqlite3, mysql, pq, mssql")
	flag.StringVar(&defaultConfig.LogFile, "log", "console", "Path to file log. 'console' is special value for print to stdout")
	flag.StringVar(&defaultConfig.Address, "address", ":8080", "VirgilD address")

	flag.StringVar(&configPath, flag.DefaultConfigFlagname, "virgild.conf", "path to config file")
}

type SignerConfig struct {
	CardID             string `json:"card_id"`
	PrivateKey         string `json:"private_key"`
	PrivateKeyPassword string `json:"private_key_password"`
}

type AuthorityConfig struct {
	CardID    string `json:"card_id"`
	PublicKey string `json:"public_key"`
}

type ServicesConfig struct {
	Cards    string `json:"cards"`
	CardsRO  string `json:"cards_ro"`
	Identity string `json:"identity"`
	VRA      string `json:"vra"`
}

type RemoteConfig struct {
	Token     string          `json:"token,omitempty"`
	Cache     int             `json:"cache,omitempty"`
	Services  ServicesConfig  `json:"services,omitempty"`
	Authority AuthorityConfig `json:"authority,omitempty"`
}

type CardsCacheConfig struct {
	Duration int
	SizeMb   int
}

type CardsConfig struct {
	Mode   string           `json:"mode"`
	Signer SignerConfig     `json:"signer,omitempty"`
	VRA    AuthorityConfig  `json:"vra,omitempty"`
	Remote RemoteConfig     `json:"remote,omitempty"`
	Cache  CardsCacheConfig `json:"cache"`
}

type AdminConfig struct {
	Enabled  bool
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthConfig struct {
	Mode      string `json:"mode"`
	TokenType string `json:"token_type"`
	Params    struct {
		Host string `json:"host"`
	} `json:"params"`
}

type Config struct {
	Admin   AdminConfig `json:"admin"`
	DB      string      `json:"db"`
	LogFile string      `json:"log,omitempty"`
	Cards   CardsConfig `json:"cards"`
	Auth    AuthConfig  `json:"auth"`
	Address string
}

func saveConfigToFole(config Config, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	saveStrVal(f, "admin-login", config.Admin.Login)
	saveStrVal(f, "admin-password", config.Admin.Password)
	saveBoolVal(f, "admin-enabled", config.Admin.Enabled)
	saveStrVal(f, "auth-mode", config.Auth.Mode)
	saveStrVal(f, "auth-address", config.Auth.Params.Host)
	saveStrVal(f, "auth-token-type", config.Auth.TokenType)
	saveStrVal(f, "mode", config.Cards.Mode)
	saveIntVal(f, "cache-duration", config.Cards.Cache.Duration)
	saveIntVal(f, "cache-size", config.Cards.Cache.SizeMb)
	saveStrVal(f, "authority-card-id", config.Cards.Remote.Authority.CardID)
	saveStrVal(f, "authority-pubkey", config.Cards.Remote.Authority.PublicKey)
	saveIntVal(f, "cache", config.Cards.Remote.Cache)
	saveStrVal(f, "cards-service", config.Cards.Remote.Services.Cards)
	saveStrVal(f, "identity-service", config.Cards.Remote.Services.Identity)
	saveStrVal(f, "ra-service", config.Cards.Remote.Services.VRA)
	saveStrVal(f, "remote-token", config.Cards.Remote.Token)
	saveStrVal(f, "vd-card-id", config.Cards.Signer.CardID)
	saveStrVal(f, "vd-key", config.Cards.Signer.PrivateKey)
	saveStrVal(f, "vd-key-password", config.Cards.Signer.PrivateKeyPassword)
	saveStrVal(f, "ra-card-id", config.Cards.VRA.CardID)
	saveStrVal(f, "ra-pubkey", config.Cards.VRA.PublicKey)
	saveStrVal(f, "db", config.DB)
	saveStrVal(f, "log", config.LogFile)
	saveStrVal(f, "address", config.Address)

	f.Close()
	return err
}

func saveStrVal(w io.Writer, name, val string) {
	fmt.Fprintf(w, "%v %v\n", name, val)
}

func saveBoolVal(w io.Writer, name string, val bool) {
	fmt.Fprintf(w, "%v %t\n", name, val)
}

func saveIntVal(w io.Writer, name string, val int) {
	fmt.Fprintf(w, "%v %v\n", name, val)
}
