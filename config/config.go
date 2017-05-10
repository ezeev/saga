package config

import (
	"reflect"
	"os"
	"strconv"
	"github.com/ezeev/saga/metrics"
)


var c SagaConfig

type SagaConfig struct {
	AppName string `yaml:"APP_NAME"`
	AppDomain string `yaml:"APP_DOMAIN"`
	TestDomain string `yaml:"TEST_DOMAIN"`
	MetricsProvider string `yaml:"METRICS_PROVIDER"`
	MetricsURI string `yaml:"METRICS_URI"`
	TemplateGlobPattern string `yaml:"TEMPLATE_GLOB_PATTERN"`
	StripeCardRedirect string `yaml:"STRIPE_CARD_REDIRECT"`
	StripeTestPublicKey string `yaml:"STRIPE_TEST_PK"`
	StripeTestSecretKey string `yaml:"STRIPE_TEST_SK"`
	StripeLivePublicKey string `yaml:"STRIPE_LIVE_PK"`
	StripeLiveSecretKey string `yaml:"STRIPE_LIVE_SK"`
	StripeAppFilter string `yaml:"STRIPE_APP_FILTER"`
	CloudSqlConnectionName string `yaml:"CLOUDSQL_CONNECTION_NAME"`
	CloudSQLDBName string `yaml:"CLOUDSQL_DB_NAME"`
	CloudSQLUser string `yaml:"CLOUDSQL_USER"`
	CloudSQLPassword string `yaml:"CLOUDSQL_PASSWORD"`
	CloudSQLDevConnStr string `yaml:"CLOUDSQL_DEV_CONN_STR"`
	Auth0Domain string `yaml:"AUTH0_DOMAIN"`
	Auth0ClientID string `yaml:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `yaml:"AUTH0_CLIENT_SECRET"`
	Auth0CallbackURI string `yaml:"AUTH0_CALLBACK_URI"`
	Auth0SignoutURI string `yaml:"AUTH0_SIGNOUT_URI"`
	Auth0CallbackHostDev string `yaml:"AUTH0_CALLBACK_HOST_DEV"`
	Auth0CallbackHostLive string `yaml:"AUTH0_CALLBACK_HOST_LIVE"`
	OAuthSuccessRedirect string `yaml:"OAUTH_SUCCESS_REDIRECT"`
	ApiRateLimitPerMin int `yaml:"API_RATE_LIMIT_PER_MIN"`
}

func Config() (*SagaConfig, error) {
	// Lazy load config when accessed
	if c == (SagaConfig{}) {
		err := Load()
		if err != nil {
			return nil, err
		}
		return &c, nil
	} else {
		return &c, nil
	}
}

func Load() error {

	ps := reflect.ValueOf(&c)
	s := ps.Elem()
	t := reflect.TypeOf(c)
	if s.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldName := field.Name
			tagVal := field.Tag.Get("yaml")
			envvar := os.Getenv(tagVal)

			//reflect.ValueOf(&c).Elem().FieldByName(fieldName).SetString(envvar)
			//use reflection to set it
			f := s.FieldByName(fieldName)
			if f.Kind() == reflect.String {
				f.SetString(envvar)
			} else if f.Kind() == reflect.Int {
				intVal, err := strconv.ParseInt(envvar,10,0)
				if err != nil {
					return err
				}
				f.SetInt(intVal)
			}

		}
	}
	metrics.Registry().IncConfigLoads()
	return nil
}


