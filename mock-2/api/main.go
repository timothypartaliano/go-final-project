// main.go

package main

import (
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"

    "final_project-ftgo-h8/api/handlers"
    "final_project-ftgo-h8/api/config"

    pb "final_project-ftgo-h8/pb"
)

func main() {
    db, err := config.InitDB()
    if err != nil {
        log.Fatal(err)
    }

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    productHandler := handlers.NewProductHandler(db)
    pb.RegisterProductServiceServer(s, productHandler)

    fmt.Println("Server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}