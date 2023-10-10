package main

import (
	"context"
	pb "final_project-ftgo-h8/pb" // Import your generated proto package
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float32
	Stock       int32
}

type ProductServiceImpl struct {
	db *gorm.DB
	pb.ProductServiceServer
}

func main() {
	// Initialize a database connection to PostgreSQL using GORM v2.
	dsn := "host=localhost user=postgres password=timothy dbname=final port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Product{}) // Auto-migrate the "products" table

	// Create a gRPC server.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &ProductServiceImpl{db: db}
	pb.RegisterProductServiceServer(s, srv)

	// Register reflection service on the gRPC server.
    reflection.Register(s)

	fmt.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateProduct implements the CreateProduct RPC method.
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	// Extract product details from the request.
	newProduct := req.GetProduct()

	// Validate input data (e.g., check for required fields).
	if newProduct.GetName() == "" || newProduct.GetPrice() <= 0 || newProduct.GetStock() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid product data")
	}

	// Create a new Product instance and set its fields.
	product := &Product{
		Name:        newProduct.GetName(),
		Description: newProduct.GetDescription(),
		Price:       newProduct.GetPrice(),
		Stock:       newProduct.GetStock(),
	}

	// Use GORM v2 to create the product in the database.
	if err := s.db.Create(product).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create product: %v", err)
	}

	// Return the created product with its generated ID.
	createdProduct := &pb.Product{
		Id:          fmt.Sprintf("%d", product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}

	return createdProduct, nil
}

// GetAllProduct implements the GetAllProduct RPC method.
func (s *ProductServiceImpl) GetAllProduct(ctx context.Context, req *pb.GetAllProductRequest) (*pb.GetAllProductResponse, error) {
    var products []*Product

    // Use GORM to retrieve all products from the database.
    if err := s.db.Find(&products).Error; err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to retrieve products: %v", err)
    }

    // Convert the database results to the protobuf response format.
    productResponses := make([]*pb.Product, len(products))
    for i, product := range products {
        productResponses[i] = &pb.Product{
            Id:          fmt.Sprintf("%d", product.ID),
            Name:        product.Name,
            Description: product.Description,
            Price:       product.Price,
            Stock:       product.Stock,
        }
    }

    // Create and return the response with the list of products.
    response := &pb.GetAllProductResponse{
        Products: productResponses,
    }

    return response, nil
}

// GetProduct implements the GetProduct RPC method.
func (s *ProductServiceImpl) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
    // Extract the product ID from the request.
    productID := req.GetId()

    // Initialize a Product instance to store the retrieved product.
    var product Product

    // Use GORM to retrieve the product from the database based on the provided ID.
    if err := s.db.First(&product, productID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, status.Errorf(codes.NotFound, "Product with ID %s not found", productID)
        }
        return nil, status.Errorf(codes.Internal, "Failed to retrieve product: %v", err)
    }

    // Convert the retrieved product to the protobuf response format.
    response := &pb.Product{
        Id:          fmt.Sprintf("%d", product.ID),
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Stock:       product.Stock,
    }

    // Return the retrieved product in the response.
    return response, nil
}

// UpdateProduct implements the UpdateProduct RPC method.
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
    // Extract product details from the request.
    updatedProduct := req.GetProduct()

    // Validate input data (e.g., check for required fields).
    if updatedProduct.GetId() == "" || updatedProduct.GetName() == "" || updatedProduct.GetPrice() <= 0 || updatedProduct.GetStock() < 0 {
        return nil, status.Errorf(codes.InvalidArgument, "Invalid product data")
    }

    // Convert the product ID from string to an integer.
    productID, err := strconv.Atoi(updatedProduct.GetId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID")
    }

    // Check if the product with the given ID exists in the database.
    var existingProduct Product
    if err := s.db.First(&existingProduct, productID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, status.Errorf(codes.NotFound, "Product with ID %s not found", updatedProduct.GetId())
        }
        return nil, status.Errorf(codes.Internal, "Failed to retrieve product: %v", err)
    }

    // Update the existing product's fields with the new data.
    existingProduct.Name = updatedProduct.GetName()
    existingProduct.Description = updatedProduct.GetDescription()
    existingProduct.Price = updatedProduct.GetPrice()
    existingProduct.Stock = updatedProduct.GetStock()

    // Use GORM to update the product in the database.
    if err := s.db.Save(&existingProduct).Error; err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update product: %v", err)
    }

    // Convert the updated product to the protobuf response format.
    updatedResponse := &pb.Product{
        Id:          fmt.Sprintf("%d", existingProduct.ID),
        Name:        existingProduct.Name,
        Description: existingProduct.Description,
        Price:       existingProduct.Price,
        Stock:       existingProduct.Stock,
    }

    // Return the updated product in the response.
    return updatedResponse, nil
}

// DeleteProduct implements the DeleteProduct RPC method.
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Product, error) {
    // Extract the product ID from the request.
    productID := req.GetId()

    // Convert the product ID from string to an integer.
    productIDInt, err := strconv.Atoi(productID)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID")
    }

    // Check if the product with the given ID exists in the database.
    var existingProduct Product
    if err := s.db.First(&existingProduct, productIDInt).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, status.Errorf(codes.NotFound, "Product with ID %s not found", productID)
        }
        return nil, status.Errorf(codes.Internal, "Failed to retrieve product: %v", err)
    }

    // Use GORM to delete the product from the database.
    if err := s.db.Delete(&existingProduct).Error; err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to delete product: %v", err)
    }

    // Convert the deleted product to the protobuf response format.
    deletedResponse := &pb.Product{
        Id:          fmt.Sprintf("%d", existingProduct.ID),
        Name:        existingProduct.Name,
        Description: existingProduct.Description,
        Price:       existingProduct.Price,
        Stock:       existingProduct.Stock,
    }

    // Return the deleted product in the response.
    return deletedResponse, nil
}