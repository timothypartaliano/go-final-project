package main

import (
	"context"
	"fmt"
	"log"
	// "os"
	// "time"

	pb "final_project-ftgo-h8/pb" // Import your generated proto package
	"google.golang.org/grpc"
)

func main() {
	// Establish a connection to the gRPC server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	// Implement client code to call gRPC methods here.

	// Example: Call CreateProduct
	createProduct(client)

	// Example: Call GetAllProduct
	getAllProduct(client)

	// Example: Call GetProduct
	getProduct(client, "1") // Provide the product ID as an argument

	// Example: Call UpdateProduct
	updateProduct(client, "1") // Provide the product ID as an argument

	// Example: Call DeleteProduct
	deleteProduct(client, "1") // Provide the product ID as an argument
}

func createProduct(client pb.ProductServiceClient) {
	// Prepare a CreateProductRequest with the product details.
	req := &pb.CreateProductRequest{
		Product: &pb.Product{
			Name:        "Example Product",
			Description: "This is an example product.",
			Price:       29.99,
			Stock:       100,
		},
	}

	// Call the CreateProduct gRPC method.
	resp, err := client.CreateProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("CreateProduct failed: %v", err)
	}

	// Display the response from the server.
	fmt.Printf("Created Product: %v\n", resp)
}

func getAllProduct(client pb.ProductServiceClient) {
	// Prepare an empty request for GetAllProduct.
	req := &pb.GetAllProductRequest{}

	// Call the GetAllProduct gRPC method.
	resp, err := client.GetAllProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("GetAllProduct failed: %v", err)
	}

	// Display the list of products received from the server.
	fmt.Println("List of Products:")
	for _, product := range resp.GetProducts() {
		fmt.Printf("ID: %s, Name: %s, Description: %s, Price: %.2f, Stock: %d\n", product.GetId(), product.GetName(), product.GetDescription(), product.GetPrice(), product.GetStock())
	}
}

func getProduct(client pb.ProductServiceClient, productID string) {
	// Prepare a GetProductRequest with the product ID.
	req := &pb.GetProductRequest{
		Id: productID,
	}

	// Call the GetProduct gRPC method.
	resp, err := client.GetProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("GetProduct failed: %v", err)
	}

	// Display the product received from the server.
	fmt.Printf("Product Details: ID: %s, Name: %s, Description: %s, Price: %.2f, Stock: %d\n", resp.GetId(), resp.GetName(), resp.GetDescription(), resp.GetPrice(), resp.GetStock())
}

func updateProduct(client pb.ProductServiceClient, productID string) {
	// Prepare an UpdateProductRequest with the updated product details.
	req := &pb.UpdateProductRequest{
		Product: &pb.Product{
			Id:          productID,
			Name:        "Updated Product",
			Description: "This is an updated product.",
			Price:       39.99,
			Stock:       150,
		},
	}

	// Call the UpdateProduct gRPC method.
	resp, err := client.UpdateProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("UpdateProduct failed: %v", err)
	}

	// Display the updated product received from the server.
	fmt.Printf("Updated Product: %v\n", resp)
}

func deleteProduct(client pb.ProductServiceClient, productID string) {
	// Prepare a DeleteProductRequest with the product ID.
	req := &pb.DeleteProductRequest{
		Id: productID,
	}

	// Call the DeleteProduct gRPC method.
	resp, err := client.DeleteProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("DeleteProduct failed: %v", err)
	}

	// Display the deleted product received from the server.
	fmt.Printf("Deleted Product: %v\n", resp)
}