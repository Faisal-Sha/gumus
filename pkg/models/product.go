package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name                 string     `json:"name"`
	CategoryPath        string     `json:"category_path"`
	Images              []string   `json:"images" gorm:"type:text[]"`
	VideoURL            *string    `json:"video_url"`
	SellerName          string     `json:"seller_name"`
	BrandName           string     `json:"brand_name"`
	RatingScore         float32    `json:"rating_score"`
	FavoriteCount       int        `json:"favorite_count"`
	CommentCount        int        `json:"comment_count"`
	AddToCartCount      int        `json:"add_to_cart_count"`
	ViewCount           int        `json:"view_count"`
	OrderCount          int        `json:"order_count"`
	TopReviews          []string   `json:"top_reviews" gorm:"type:text[]"`
	SizeRecommendation  string     `json:"size_recommendation"`
	EstimatedDelivery   *time.Time `json:"estimated_delivery"`
	DeliveryTimeWindow  string     `json:"delivery_time_window"`
	IsActive            bool       `json:"is_active" gorm:"default:true"`
	Price               float64    `json:"price"`
	InstallmentOptions  []string   `json:"installment_options" gorm:"type:text[]"`
	SimilarProducts     []string   `json:"similar_products" gorm:"type:text[]"`
	ComplementaryProducts []string `json:"complementary_products" gorm:"type:text[]"`
}

type ProductVariant struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Size     string  `json:"size"`
	Color    string  `json:"color"`
	Stock    int     `json:"stock"`
	Price    float64 `json:"price"`
	IsActive bool    `json:"is_active" gorm:"default:true"`
}

type ProductAttribute struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type PriceHistory struct {
	gorm.Model
	ProductID    uint    `json:"product_id"`
	VariantID    uint    `json:"variant_id"`
	OldPrice     float64 `json:"old_price"`
	NewPrice     float64 `json:"new_price"`
	ChangedAt    time.Time
}
