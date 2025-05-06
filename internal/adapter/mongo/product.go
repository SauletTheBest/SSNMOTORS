package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"AP2Assignment2/inventory-service/internal/adapter/mongo/dao"
	"AP2Assignment2/inventory-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepo represents the adapter layer for products
type ProductRepo struct {
	conn       *mongo.Database
	collection string
}

// NewProductRepo initializes the product adapter
func NewProductRepo(conn *mongo.Database) *ProductRepo {
	return &ProductRepo{
		conn:       conn,
		collection: CollectionProducts,
	}
}

// Create inserts a new product into the database
func (p *ProductRepo) Create(ctx context.Context, product domain.Product) error {
	productDoc := dao.FromProduct(product)
	_, err := p.conn.Collection(p.collection).InsertOne(ctx, productDoc)
	if err != nil {
		return fmt.Errorf("product with ID %d has not been created: %w", product.ID, err)
	}

	return nil
}

// Update modifies an existing product based on a filter
func (p *ProductRepo) Update(ctx context.Context, filter domain.ProductFilter, update domain.ProductUpdateData) error {
	res, err := p.conn.Collection(p.collection).UpdateOne(
		ctx,
		dao.FromProductFilter(filter),
		dao.FromProductUpdateData(update),
	)
	if err != nil {
		return fmt.Errorf("product has not been updated with filter: %v, err: %w", filter, err)
	}

	if res.ModifiedCount == 0 {
		return fmt.Errorf("product has not been updated with filter: %v", filter)
	}

	return nil
}

// GetWithFilter retrieves a single product matching the filter
func (p *ProductRepo) GetWithFilter(ctx context.Context, filter domain.ProductFilter) (domain.Product, error) {
	var daoProduct dao.Product
	err := p.conn.Collection(p.collection).FindOne(
		ctx,
		dao.FromProductFilter(filter),
	).Decode(&daoProduct)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Product{}, domain.ErrProductNotFound // No product found
		}
		return domain.Product{}, fmt.Errorf("failed to find product: %w", err)
	}

	product := dao.ToProduct(daoProduct)
	return product, nil
}

// GetListWithFilter retrieves multiple products based on a filter
func (p *ProductRepo) GetListWithFilter(ctx context.Context, filter domain.ProductFilter, page, limit int64) ([]domain.Product, int, error) {
	// Create the filter for the query
	findFilter := dao.FromProductFilter(filter)

	// Calculate pagination parameters
	skip := (page - 1) * limit

	// Set up find options for pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Get total count
	totalCount, err := p.conn.Collection(p.collection).CountDocuments(ctx, findFilter)
	if err != nil {
		return nil, 0, err
	}

	// Execute the query with pagination
	cursor, err := p.conn.Collection(p.collection).Find(ctx, findFilter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}(cursor, ctx)

	var daoProducts []dao.Product
	if err := cursor.All(ctx, &daoProducts); err != nil {
		return nil, 0, fmt.Errorf("failed to decode products: %w", err)
	}
	products := dao.ToProductList(daoProducts)
	return products, int(totalCount), nil
}

// Delete permanently deletes a product from the database
func (p *ProductRepo) Delete(ctx context.Context, filter domain.ProductFilter) error {
	res, err := p.conn.Collection(p.collection).DeleteOne(
		ctx,
		dao.FromProductFilter(filter),
	)
	if err != nil {
		return fmt.Errorf("product has not been deleted with filter: %v, err: %w", filter, err)
	}

	if res.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}
