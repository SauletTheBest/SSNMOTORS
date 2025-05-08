package dao

import (
	domain "AP2Assignment2/inventory-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Product struct {
	ID        uint64    `bson:"_id"`
	Name      string    `bson:"name"`
	Category  string    `bson:"category"`
	Price     float64   `bson:"price"`
	Stock     uint64    `bson:"stock"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func ToProductList(daoProducts []Product) []domain.Product {
	products := make([]domain.Product, len(daoProducts))
	for i, p := range daoProducts {
		products[i] = domain.Product{
			ID:        p.ID,
			Name:      p.Name,
			Category:  p.Category,
			Price:     p.Price,
			Stock:     p.Stock,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return products
}

func ToProduct(product Product) domain.Product {
	return domain.Product{
		ID:        product.ID,
		Name:      product.Name,
		Category:  product.Category,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func FromProduct(product domain.Product) Product {
	return Product{
		ID:        product.ID,
		Name:      product.Name,
		Category:  product.Category,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func FromProductFilter(filter domain.ProductFilter) bson.M {
	query := bson.M{}

	if filter.ID != nil {
		query["_id"] = *filter.ID
	}

	if filter.Name != nil {
		query["name"] = *filter.Name
	}

	if filter.Category != nil {
		query["category"] = *filter.Category
	}

	if filter.Price != nil {
		query["price"] = *filter.Price
	}

	if filter.Stock != nil {
		query["stock"] = *filter.Stock
	}

	return query
}

func FromProductUpdateData(updateData domain.ProductUpdateData) bson.M {
	query := bson.M{}

	if updateData.Name != nil {
		query["name"] = *updateData.Name
	}

	if updateData.Category != nil {
		query["category"] = *updateData.Category
	}

	if updateData.Price != nil {
		query["price"] = *updateData.Price
	}

	if updateData.Stock != nil {
		query["stock"] = *updateData.Stock
	}

	if updateData.UpdatedAt != nil {
		query["updatedAt"] = updateData.UpdatedAt
	}

	return bson.M{"$set": query}
}
