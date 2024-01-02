package entities

type Wishlist struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"unique"`
}

type WishlistItems struct {
	ID        uint `gorm:"primaryKey"`
	WishlistID    uint `gorm:"foreignKey:WishlistID;references:wishlists(id)"`
	ProductID uint
}