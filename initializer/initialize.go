package initializer

import (
	"github.com/akshay0074700747/wishlist-service/adapters"
	"github.com/akshay0074700747/wishlist-service/service"
	"gorm.io/gorm"
)


func InitService(db *gorm.DB) *service.WishlistService {
	return service.NewWishlistService(adapters.NewWishlistAdapter(db))
}