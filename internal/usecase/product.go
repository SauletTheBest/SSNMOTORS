package usecase

import (
	"AP2Assignment2/inventory-service/internal/adapter/mongo"
	"AP2Assignment2/inventory-service/internal/domain"
	"context"
)

type Product struct {
	aiRepo auto_inc_Repo
	repo   product_Repo
}

func NewProduct(aiRepo auto_inc_Repo, repo product_Repo) *Product {
	return &Product{
		aiRepo: aiRepo,
		repo:   repo,
	}
}

func (p *Product) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	id, err := p.aiRepo.Next(ctx, mongo.CollectionProducts)
	if err != nil {
		return domain.Product{}, err
	}
	product.ID = id
	err = p.repo.Create(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}
	return domain.Product{
		ID:   id,
		Name: product.Name,
	}, nil
}

func (p *Product) Get(ctx context.Context, pf domain.ProductFilter) (domain.Product, error) {
	product, err := p.repo.GetWithFilter(ctx, pf)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (p *Product) GetAll(ctx context.Context, pf domain.ProductFilter, page, limit int64) ([]domain.Product, int, error) {
	products, totalCount, err := p.repo.GetListWithFilter(ctx, pf, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}

func (p *Product) Update(ctx context.Context, filter domain.ProductFilter, updated domain.ProductUpdateData) error {
	err := p.repo.Update(ctx, filter, updated)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) Delete(ctx context.Context, filter domain.ProductFilter) error {
	err := p.repo.Delete(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
