package usecase

import (
    "context"
    "errors"
    "inventory-service/internal/model"
    "inventory-service/internal/repository"
)

type ProductUsecase struct {
    repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) *ProductUsecase {
    return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, p *model.Product) (string, error) {
    if p.Name == "" {
        return "", errors.New("name is required")
    }
    if p.Description == "" {
        return "", errors.New("description is required")
    }
    if p.Category == "" {
        return "", errors.New("category is required")
    }
    if p.Stock < 0 {
        return "", errors.New("stock cannot be negative")
    }
    if p.Price < 0 {
        return "", errors.New("price cannot be negative")
    }
    return u.repo.Create(ctx, p)
}

func (u *ProductUsecase) GetProduct(ctx context.Context, id string) (*model.Product, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }
    return u.repo.GetByID(ctx, id)
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, p *model.Product) error {
    if p.ID == "" {
        return errors.New("id is required")
    }
    // Все поля обязательны — нет частичного обновления
    if p.Name == "" || p.Description == "" || p.Category == "" {
        return errors.New("name, description and category are required")
    }
    if p.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    if p.Price < 0 {
        return errors.New("price cannot be negative")
    }
    return u.repo.Update(ctx, p)
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("id is required")
    }
    return u.repo.Delete(ctx, id)
}

func (u *ProductUsecase) ListProducts(ctx context.Context, category string, page, limit int32) ([]*model.Product, error) {
    if page < 1 || limit < 1 {
        return nil, errors.New("page and limit must be ≥ 1")
    }
    return u.repo.List(ctx, category, page, limit)
}

func (u *ProductUsecase) DecreaseStock(ctx context.Context, productID string, quantity int32) error {
    return u.repo.DecreaseStock(ctx, productID, quantity)
}

