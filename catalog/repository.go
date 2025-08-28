package catalog

import (
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/olivere/elastic/v7"
)

type Repository interface {
	Close()
	PostCatalog(ctx context.Context, catalog Catalog) error
	GetProducts(ctx context.Context, taken uint64, skip uint64) ([]Catalog, error)
	GetProductsByID(ctx context.Context, id string) (*Catalog, error)
	GetProductsWithIds(ctx context.Context, ids []string) ([]Catalog, error)
	SearchProduct(ctx context.Context, query string, taken uint64, skip uint64) ([]Catalog, error)
}

type elasticRepository struct {
	client *elastic.Client
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	exists, err := client.IndexExists("catalog").Do(context.Background())
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err = client.CreateIndex("catalog").BodyString(`{
			"mappings": {
				"properties": {
					"name":        { "type": "text" },
					"description": { "type": "text" },
					"price":       { "type": "float" }
				}
			}
		}`).Do(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return &elasticRepository{
		client: client,
	}, nil
}

func (repository *elasticRepository) Close() {
	// elastic client doesn’t need explicit close in v7
}

func (repository *elasticRepository) PostCatalog(ctx context.Context, catalog Catalog) error {
	_, err := repository.client.
		Index().
		Index("catalog").
		Id(catalog.ID).
		BodyJson(productDocument{
			Name:        catalog.Name,
			Description: catalog.Description,
			Price:       catalog.Price,
		}).
		Do(ctx)

	return err
}

func (repository *elasticRepository) GetProducts(ctx context.Context, taken uint64, skip uint64) ([]Catalog, error) {
	res, err := repository.client.Search().
		Index("catalog").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(taken)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	catalogs := []Catalog{}
	for _, hit := range res.Hits.Hits {
		product := productDocument{}
		if err := json.Unmarshal(hit.Source, &product); err == nil {
			catalogs = append(catalogs, Catalog{
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				ID:          hit.Id,
			})
		}
	}

	return catalogs, nil
}

func (repository *elasticRepository) GetProductsByID(ctx context.Context, id string) (*Catalog, error) {
	res, err := repository.client.
		Get().
		Index("catalog").
		Id(id).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}

	product := productDocument{}

	if err := json.Unmarshal(res.Source, &product); err != nil {
		return nil, err
	}

	return &Catalog{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ID:          id,
	}, nil
}

func (repository *elasticRepository) GetProductsWithIds(ctx context.Context, ids []string) ([]Catalog, error) {
	multi := repository.client.MultiGet()
	for _, id := range ids {
		multi = multi.Add(elastic.NewMultiGetItem().Index("catalog").Id(id))
	}

	res, err := multi.Do(ctx)
	if err != nil {
		return nil, err
	}

	catalogs := []Catalog{}
	for _, doc := range res.Docs {
		if doc.Found {
			product := productDocument{}
			if err := json.Unmarshal(doc.Source, &product); err == nil {
				catalogs = append(catalogs, Catalog{
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
					ID:          doc.Id,
				})
			}
		}
	}
	return catalogs, nil
}

func (repository *elasticRepository) SearchProduct(ctx context.Context, query string, taken uint64, skip uint64) ([]Catalog, error) {
	res, err := repository.client.
		Search().
		Index("catalog").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(taken)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	catalogs := []Catalog{}
	for _, hit := range res.Hits.Hits {
		product := productDocument{}
		if err := json.Unmarshal(hit.Source, &product); err == nil {
			catalogs = append(catalogs, Catalog{
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				ID:          hit.Id,
			})
		}
	}
	return catalogs, nil
}

// elastic model for request and response
type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// error
var (
	ErrNotFound = errors.New("entity not found")
)
