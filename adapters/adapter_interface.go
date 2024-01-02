package adapters

import "github.com/akshay0074700747/wishlist-service/entities"

type WishlistAdapterInterface interface {
	CreateWishlist(req entities.Wishlist) (entities.Wishlist, error)
	InsertIntoWishlist(req entities.WishlistItems,user_id uint) (entities.WishlistItems, error)
	GetWishlistItems(user_id uint) ([]entities.WishlistItems, error)
	DeleteWishlistItem(req entities.WishlistItems,user_id uint) error
	TruncateWishlistItems(user_id uint) (error)
}
