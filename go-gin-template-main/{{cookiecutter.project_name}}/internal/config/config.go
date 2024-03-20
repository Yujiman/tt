package config

import (
	"context"
	"fmt"
	"strings"
	"sync"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Service       *Service
		Worker        *Worker
		Database      *Database
		Authorization *Authorization
		Nats          *Nats
		Integration   *Integration
		Sentry        *Sentry
		Vault         *Vault
	}

	Service struct {
		AppName         string         `envconfig:"APP_NAME" required:"true"`
		Environment     AppEnvironment `envconfig:"ENVIRONMENT" default:"local"`
		Port            string         `envconfig:"PORT" default:"8080"`
		Domain          string         `envconfig:"DOMAIN" default:"localhost"`
		Limiter         string         `envconfig:"LIMITER_SETTINGS"`
		OpenapiEndpoint string         `envconfig:"OPENAPI_ENDPOINT" required:"false"`
	}

	Authorization struct {
		BaseURL string `envconfig:"AUTHORIZATION_BASE_URL" required:"false"`
	}

	Vault struct {
		URL           string `envconfig:"VAULT_URL" required:"true"`
		Namespace     string `envconfig:"NAMESPACE" required:"true"`
		CommonKV      string `envconfig:"VAULT_COMMON_KV" required:"true"`
		VaultUser     string `envconfig:"VAULT_USER" default:""`
		VaultPassword string `envconfig:"VAULT_PASSWORD" default:""`
	}

	Worker struct {
		RunEchoWorker bool `envconfig:"RUN_ECHO_WORKER" default:"true"`
	}

	Database struct {
		ReadDSNs  []string `envconfig:"DATABASE_READ_DSN"`
		WriteDSNs []string `envconfig:"DATABASE_WRITE_DSN"`
		Schema    string   `envconfig:"DATABASE_SCHEMA"`

		MongoDSN string `envconfig:"MONGO_DSN"`
	}

	Nats struct {
		DSN string `envconfig:"NATS_DSN"`
	}

	Integration struct {
		ExampleIntegrationBaseUrl string `envconfig:"EXAMPLE_INTEGRATION_BASE_URL"`
		ExampleIntegrationToken   string `envconfig:"EXAMPLE_INTEGRATION_TOKEN"`
	}

	Sentry struct {
		DSN string `envconfig:"SENTRY_DSN"`
	}
)

var (
	once   sync.Once
	config *Config
)

// GetConfig Загружает конфиг из .env файла и возвращает объект конфигурации
// В случае, если не передать параметр `envfiles`, берется `.env` файл из корня проекта
func GetConfig(envfiles ...string) (*Config, error) {
	var err error
	once.Do(func() {
		_ = godotenv.Load(envfiles...)

		var c Config
		err = envconfig.Process("", &c)
		if err != nil {
			err = fmt.Errorf("error parse config from env variables: %w\n", err)
			return
		}

		vaultErr := updateCfgFromVault(&c)
		if vaultErr != nil {
			// Есть возможность добить в env на случай если в vault едоступен
			fmt.Println(vaultErr.Error())
		}

		if e := c.Service.Environment.Validate(); e != nil {
			err = fmt.Errorf("error parse config from env variables: %w\n", e)
			return
		}

		config = &c
	})

	return config, err
}

func updateCfgFromVault(cfg *Config) error {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = cfg.Vault.URL
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return fmt.Errorf("can't create valut client: %w", err)
	}

	if cfg.Vault.VaultUser == "" {
		k8sAuth, k8sErr := auth.NewKubernetesAuth(
			cfg.Vault.Namespace,
			auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
		)
		if k8sErr != nil {
			return fmt.Errorf("k8s auth error: %w", k8sErr)
		}

		authInfo, k8sErr := client.Auth().Login(context.Background(), k8sAuth)
		if k8sErr != nil {
			return fmt.Errorf("unable to log in with Kubernetes auth: %w", k8sErr)
		}
		if authInfo == nil {
			return fmt.Errorf("no auth info was returned after login to vault")
		}
	} else {
		path := fmt.Sprintf("auth/userpass/login/%s", cfg.Vault.VaultUser)
		secret, err := client.Logical().Write(path, map[string]interface{}{
			"password": cfg.Vault.VaultPassword,
		})
		if err != nil {
			return fmt.Errorf("can't auth in vault: %w", err)
		}

		client.SetToken(secret.Auth.ClientToken)
	}

	commonSecret, err := client.KVv2(cfg.Vault.Namespace).Get(context.Background(), cfg.Vault.CommonKV)
	if err != nil {
		return fmt.Errorf(
			"can't get commonSecrets in namespace '%s', commonkv '%s': %w",
			cfg.Vault.Namespace,
			cfg.Vault.CommonKV,
			err,
		)
	}

	// Расскомментить, если для сервиса используются какие-то уникальные секреты
	// serviceSecret, err := client.KVv2(cfg.Vault.Namespace).Get(context.Background(), cfg.Service.AppName)
	// if err != nil {
	// 	return fmt.Errorf(
	// 		"can't get serviceSecret in namespace '%s' for service '%s': %w",
	// 		cfg.Vault.Namespace,
	// 		cfg.Service.AppName,
	// 		err,
	// 	)
	// }

	writeDSNs := commonSecret.Data["DATABASE_WRITE_DSN"].(string)
	cfg.Database.WriteDSNs = strings.Split(writeDSNs, ",")

	if commonSecret.Data["DATABASE_READ_DSN"] != nil {
		readDSNs := commonSecret.Data["DATABASE_READ_DSN"].(string)
		cfg.Database.ReadDSNs = strings.Split(readDSNs, ",")
	}

	cfg.Database.MongoDSN = commonSecret.Data["MONGO_DSN"].(string)
	cfg.Nats.DSN = commonSecret.Data["NATS_DSN"].(string)
	cfg.Service.OpenapiEndpoint = commonSecret.Data["OPENAPI_ENDPOINT"].(string)

	if !cfg.Service.Environment.IsLocal() {
		cfg.Authorization.BaseURL = commonSecret.Data["AUTHORIZATION_BASE_URL"].(string)
	}

	return nil
}

type AppEnvironment string

const (
	PRODUCTION  AppEnvironment = "prod"
	STAGE       AppEnvironment = "stage"
	DEVELOPMENT AppEnvironment = "dev"
	LOCAL       AppEnvironment = "local"
)

func (e AppEnvironment) IsProduction() bool {
	return e == PRODUCTION
}

func (e AppEnvironment) IsStage() bool {
	return e == STAGE
}

func (e AppEnvironment) IsDevelopment() bool {
	return e == DEVELOPMENT
}

func (e AppEnvironment) IsLocal() bool {
	return e == LOCAL
}

func (e AppEnvironment) String() string {
	return string(e)
}

func (e AppEnvironment) Validate() error {
	switch e {
	case LOCAL, DEVELOPMENT, STAGE, PRODUCTION:
		return nil
	default:
		return fmt.Errorf("unexpected ENVIRONMENT in .env: %s", e)
	}
}
