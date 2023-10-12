package controller

import (
	"final_project-ftgo-h8/pb"
	"final_project-ftgo-h8/api/dto"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

// CreateProduct handles the creation of a new product.
func (c *ProductController) CreateProduct(ctx echo.Context) error {
	// Bind the request body
	reqBody := dto.ReqBodyCreateProduct{}
	if err := ctx.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(400, echo.Map{
			"message": "error to bind",
			"detail":  err.Error(),
		})
	}

	// Create a gRPC request
	req := pb.CreateProductRequest{
		Product: &pb.Product{
			Name:        reqBody.Name,
			Description: reqBody.Description,
			Price:       reqBody.Price,
			Stock:       int32(reqBody.Stock),
		},
	}

	// Call the gRPC service method
	newProduct, err := c.Service.CreateProduct(ctx.Request().Context(), &req)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.JSON(201, echo.Map{
		"message": "success create",
		"detail":  newProduct,
	})
}

// GetAllProducts retrieves all products.
func (c *ProductController) GetAllProducts(ctx echo.Context) error {
    // Call the gRPC service method to get all products
    allProducts, err := c.Service.GetAllProduct(ctx.Request().Context(), &pb.GetAllProductRequest{})
    if err != nil {
        return handleError(ctx, err)
    }

    return ctx.JSON(200, allProducts)
}

// GetProduct retrieves a product by ID.
func (c *ProductController) GetProduct(ctx echo.Context) error {
	// Extract product ID from URL param
	productID := ctx.Param("id")

	// Call the gRPC service method
	product, err := c.Service.GetProduct(ctx.Request().Context(), &pb.GetProductRequest{Id: productID})
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.JSON(200, product)
}

// UpdateProduct updates a product by ID.
func (c *ProductController) UpdateProduct(ctx echo.Context) error {
	// Extract product ID from URL param
	productID := ctx.Param("id")

	// Bind the request body
	reqBody := dto.ReqBodyUpdateProduct{}
	if err := ctx.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(400, echo.Map{
			"message": "error to bind",
			"detail":  err.Error(),
		})
	}

	// Create a gRPC request
	req := pb.UpdateProductRequest{
		Product: &pb.Product{
			Id:          productID,
			Name:        reqBody.Name,
			Description: reqBody.Description,
			Price:       reqBody.Price,
			Stock:       int32(reqBody.Stock),
		},
	}

	// Call the gRPC service method
	updatedProduct, err := c.Service.UpdateProduct(ctx.Request().Context(), &req)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.JSON(200, updatedProduct)
}

// DeleteProduct deletes a product by ID.
func (c *ProductController) DeleteProduct(ctx echo.Context) error {
	// Extract product ID from URL param
	productID := ctx.Param("id")

	// Call the gRPC service method
	_, err := c.Service.DeleteProduct(ctx.Request().Context(), &pb.DeleteProductRequest{Id: productID})
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.NoContent(204)
}

// handleError handles gRPC errors and converts them to Echo HTTP errors.
func handleError(ctx echo.Context, err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return echo.NewHTTPError(404, echo.Map{"message": "not found", "detail": st.Message()})
		case codes.InvalidArgument:
			return echo.NewHTTPError(400, echo.Map{"message": "invalid argument", "detail": st.Message()})
		default:
			return echo.NewHTTPError(500, echo.Map{"message": "internal server error", "detail": st.Message()})
		}
	}
	return err
}