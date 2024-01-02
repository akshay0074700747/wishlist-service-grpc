package service

import (
	"context"
	"errors"
	"log"

	"github.com/akshay0074700747/proto-files-for-microservices/pb"
	"github.com/akshay0074700747/wishlist-service/adapters"
	"github.com/akshay0074700747/wishlist-service/entities"
)

type WishlistService struct {
	Adapter adapters.WishlistAdapterInterface
	pb.UnimplementedWishlistServiceServer
}

var (
	productClient pb.ProductServiceClient
)

func InitClients(product pb.ProductServiceClient) {
	productClient = product
}

func NewWishlistService(adapter adapters.WishlistAdapterInterface) *WishlistService {
	return &WishlistService{
		Adapter: adapter,
	}
}

func (wishlist *WishlistService) CreateWishlist(ctx context.Context, req *pb.WishlistRequest) (*pb.WishlistResponce, error) {

	res, err := wishlist.Adapter.CreateWishlist(entities.Wishlist{UserID: uint(req.UserId)})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.WishlistResponce{WishlistId: uint32(res.ID), UserId: uint32(res.UserID), IsEmpty: true}, nil
}

func (wishlist *WishlistService) GetWishlist(ctx context.Context, req *pb.WishlistRequest) (*pb.GetWishlistResponce, error) {

	var ids []uint32

	res, err := wishlist.Adapter.GetWishlistItems(uint(req.UserId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if len(res) == 0 {
		return &pb.GetWishlistResponce{}, nil
	}

	for _, prod := range res {
		ids = append(ids, uint32(prod.ProductID))
	}

	productRes, err := productClient.GetArrayofProducts(ctx, &pb.ArrayofProductsRequest{Id: ids})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.GetWishlistResponce{WishlistId: uint32(res[0].WishlistID), Products: productRes.Products}, nil
}

func (wishlist *WishlistService) AddtoWishlist(ctx context.Context, req *pb.AddtoWishlistRequest) (*pb.AddProductResponce, error) {

	productRes, err := productClient.GetProduct(context.TODO(), &pb.GetProductByID{Id: req.ProductId})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if productRes.Name == "" {
		return nil, errors.New("Sorry the product doesnt exist")
	}

	_, err = wishlist.Adapter.InsertIntoWishlist(entities.WishlistItems{
		ProductID: uint(req.ProductId)}, uint(req.UserId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return productRes, nil
}

func (wishlist *WishlistService) DeleteWishlistItem(ctx context.Context, req *pb.AddtoWishlistRequest) (*pb.GetWishlistResponce, error) {

	productRes, err := productClient.GetProduct(context.TODO(), &pb.GetProductByID{Id: req.ProductId})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if productRes.Name == "" {
		return nil, errors.New("A product with the given ID doesnt exist")
	}

	err = wishlist.Adapter.DeleteWishlistItem(entities.WishlistItems{ProductID: uint(req.ProductId)}, uint(req.UserId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.GetWishlistResponce{}, nil
}

func (wishlist *WishlistService) TruncateWishlist(ctx context.Context, req *pb.WishlistRequest) (*pb.WishlistResponce, error) {

	if err := wishlist.Adapter.TruncateWishlistItems(uint(req.UserId)); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.WishlistResponce{UserId: req.UserId, IsEmpty: true}, nil
}
