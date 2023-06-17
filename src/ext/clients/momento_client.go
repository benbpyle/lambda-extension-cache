package clients

import (
	"cache-layer/models"
	"context"
	"encoding/json"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type CacheRepository struct {
	client momento.CacheClient
}

func NewCacheRepository(client momento.CacheClient) *CacheRepository {
	return &CacheRepository{
		client: client,
	}
}

func NewMomentoClient(token string) momento.CacheClient {
	// Initializes Momento
	credentialProvider, err := auth.FromString(token)

	if err != nil {
		return nil
	}

	// Initializes Momento
	client, err := momento.NewCacheClient(
		config.LaptopLatest(),
		credentialProvider,
		600*time.Second)

	if err != nil {
		return nil
	}

	return client
}

func (c *CacheRepository) WriteCache(ctx context.Context, model *models.Model) error {
	b, err := json.Marshal(model)
	if err != nil {
		return err
	}

	_, err = c.client.Set(ctx, &momento.SetRequest{
		CacheName: "SampleCache",
		Key:       momento.String(model.Id),
		Value:     momento.Bytes(b),
	})

	return err
}

func (c *CacheRepository) ReadCache(ctx context.Context, id string) (*models.Model, error) {
	resp, err := c.client.Get(ctx, &momento.GetRequest{
		CacheName: "SampleCache",
		Key:       momento.String(id),
	})

	if err != nil {
		return nil, err
	}

	// See if you got a cache hit or cache miss
	switch r := resp.(type) {
	case *responses.GetHit:
		var m models.Model
		err = json.Unmarshal(r.ValueByte(), &m)
		if err != nil {
			return nil, err
		}

		return &m, nil
	}

	return nil, nil

}

func GetSecretString() (*string, error) {
	cache, err := secretcache.New()
	if err != nil {
		return nil, err
	}

	secretString, err := cache.GetSecretString("mo-data-flow-router-cache-token")

	if err != nil {
		return nil, err
	}

	ss := struct {
		Token string `json:"token"`
	}{}

	err = json.Unmarshal([]byte(secretString), &ss)

	if err != nil {
		return nil, err
	}

	return &ss.Token, nil
}
