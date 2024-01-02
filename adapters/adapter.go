package adapters

import (
	"errors"

	"github.com/akshay0074700747/wishlist-service/entities"
	"gorm.io/gorm"
)

type WishlistAdapter struct {
	DB *gorm.DB
}

func NewWishlistAdapter(db *gorm.DB) *WishlistAdapter {
	return &WishlistAdapter{
		DB: db,
	}
}

func (wishlist *WishlistAdapter) CreateWishlist(req entities.Wishlist) (entities.Wishlist, error) {

	var res entities.Wishlist
	query := "INSERT INTO wishlists (user_id) VALUES($1) RETURNS id,user_id"

	tx := wishlist.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := wishlist.DB.Raw(query, req.UserID).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return entities.Wishlist{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (wishlist *WishlistAdapter) InsertIntoWishlist(req entities.WishlistItems, user_id uint) (entities.WishlistItems, error) {

	var res entities.WishlistItems
	query := "INSERT INTO wishlist_items (wishlist_id,product_id) SELECT w.id AS wishlist_id, $1 AS product_id FROM wishlist w WHERE user_id = $2 RETURNING id,wishlist_id,product_id"

	tx := wishlist.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := wishlist.DB.Raw(query, req.WishlistID, req.ProductID).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return entities.WishlistItems{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (wishlist *WishlistAdapter) GetWishlistItems(user_id uint) ([]entities.WishlistItems, error) {

	var res []entities.WishlistItems
	query := "SELECT * FROM wishlist_items WHERE wishlist_id = (SELECT id FROM wishlists WHERE user_id = $1)"

	return res, wishlist.DB.Raw(query, user_id).Scan(&res).Error
}

func (wishlist *WishlistAdapter) DeleteWishlistItem(req entities.WishlistItems, user_id uint) error {

	query := "DELETE FROM wishlist_items WHERE product_id = $1 AND wishlist_id = (SELECT id FROM wishlists WHERE user_id = $2)"
	res := wishlist.DB.Exec(query, req.ProductID, user_id)

	tx := wishlist.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("Wishlist Item not deleted")
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (wishlist *WishlistAdapter) TruncateWishlistItems(user_id uint) error {

	query := "DELETE FROM wishlist_items WHERE wishlist_id = (SELECT id FROM wishlists WHERE user_id = $1)"

	tx := wishlist.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := wishlist.DB.Raw(query, user_id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
