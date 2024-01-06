package main

import (
	"log"
	"net"
	"os"

	"github.com/akshay0074700747/proto-files-for-microservices/pb"
	"github.com/akshay0074700747/wishlist-service/db"
	"github.com/akshay0074700747/wishlist-service/initializer"
	"github.com/akshay0074700747/wishlist-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	if godotenv.Load(".env") != nil {
		log.Fatal("couldnt load env files")
	}

	addr := os.Getenv("DATABASE_ADDR")
	if addr == "" {
		log.Fatal("address not found")
	}

	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	servicee := initializer.InitService(DB)
	server := grpc.NewServer()

	pb.RegisterWishlistServiceServer(server, servicee)

	productConn, err := grpc.Dial("product-service:50004", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}

	service.InitClients(pb.NewProductServiceClient(productConn))

	listener, err := net.Listen("tcp", ":50007")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Wishlist-server is llistening or port 50007")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}
