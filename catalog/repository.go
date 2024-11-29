package catalog

import (
	"context"
	"encoding/json"
	"errors"

	elastic "gopkg.in/olivere/elastic.v5"
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
	return &elasticRepository{
		client: client,
	}, nil
}

func (repository *elasticRepository) Close() {

}
func (repository *elasticRepository) PostCatalog(ctx context.Context, catalog Catalog) error {
	_, err := repository.client.
		Index().Index("catalog").
		Type("product").
		Id(catalog.ID).
		BodyJson(
			productDocumnet{
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
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(taken)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	catalogs := []Catalog{}
	for _, hit := range res.Hits.Hits {
		product := productDocumnet{}
		if err := json.Unmarshal(*hit.Source, &product); err == nil {
			catalogs = append(catalogs, Catalog{
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				ID:          hit.Id,
			})
		}
	}

	return catalogs, err
}

func (repository *elasticRepository) GetProductsByID(ctx context.Context, id string) (*Catalog, error) {
	res, err := repository.client.
		Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}

	product := productDocumnet{}

	if err := json.Unmarshal(*res.Source, &product); err != nil {
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

	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(
			items,
			elastic.NewMultiGetItem().Index("catalog").
				Type("product").Id(id),
		)
	}

	res, err := repository.client.
		MultiGet().
		Add(items...).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	catalogs := []Catalog{}
	for _, doc := range res.Docs {
		product := productDocumnet{}
		if err := json.Unmarshal(*doc.Source, &product); err == nil {
			catalogs = append(catalogs, Catalog{
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				ID:          doc.Id,
			})
		}
	}
	return catalogs, err
}

func (repository *elasticRepository) SearchProduct(ctx context.Context, query string, taken uint64, skip uint64) ([]Catalog, error) {
	res, err := repository.client.
		Search().
		Index("catalog").Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(taken)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	catalogs := []Catalog{}
	for _, hit := range res.Hits.Hits {
		product := productDocumnet{}
		if err := json.Unmarshal(*hit.Source, &product); err == nil {
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
type productDocumnet struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// error
var (
	ErrNotFound = errors.New("entity not found ")
)
