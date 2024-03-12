package bootstrap

import (
	"context"
	"strings"

	"github.com/lengocson131002/go-clean-core/config"
	"github.com/lengocson131002/go-clean-core/es"
	"github.com/lengocson131002/go-clean-core/logger"
)

func NewElasticSearchClient(cfg config.Configure, logger logger.Logger) es.ElasticSearchClient {
	var (
		addresses = strings.Split(cfg.GetString("ELASTIC_SEARCH_ADDRESSES"), ",")
		username  = cfg.GetString("ELASTIC_SEARCH_USERNAME")
		password  = cfg.GetString("ELASTIC_SEARCH_PASSWORD")
	)

	client, err := es.NewElasticSearchClient(
		es.WithEsAddresses(addresses),
		es.WithEsUsername(username),
		es.WithEsPassword(password),
	)

	if err != nil {
		panic(err)
	}

	logger.Info(context.TODO(), "Connected to Elasticsearch")
	return client
}
