package database

import (
	"log"

	"shopProject/internal/models"

	"gorm.io/gorm"
)

// SeedData inserts sample categories and products if the database is empty
func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	log.Println("Seeding initial store data...")

	// 1. Categories
	catMen := models.Category{
		Name:        "پوشاک مردانه",
		ImageURL:    "https://images.unsplash.com/photo-1490578474895-699cd4e2cf59",
		Description: "انواع لباس، شلوار و کت مردانه شیک",
	}
	catWomen := models.Category{
		Name:        "پوشاک زنانه",
		ImageURL:    "https://images.unsplash.com/photo-1483985988355-763728e1935b",
		Description: "جدیدترین مدل‌های پوشاک و اکسسوری زنانه",
	}
	catDigital := models.Category{
		Name:        "لوازم دیجیتال",
		ImageURL:    "https://images.unsplash.com/photo-1505740420928-5e560c06d30e",
		Description: "هدفون، گجت و ساعت‌های هوشمند",
	}

	db.Create(&catMen)
	db.Create(&catWomen)
	db.Create(&catDigital)

	// 2. Products
	products := []models.Product{
		{
			CategoryID:      catMen.ID,
			Name:            "پیراهن نخی آستین بلند مردانه",
			ImageURL:        "https://images.unsplash.com/photo-1602810318383-e386cc2a3ccf",
			Description:     "تهیه شده از پارچه ۱۰۰٪ پنبه، بسیار خنک و شیک مناسب فصول گرم و مهمانی‌ها.",
			Price:           580000,
			DiscountPercent: 20,
			IsSpecialOffer:  true, // پیشنهاد شگفت انگیز
			Weight:          "250 گرم",
			Color:           "سرمه‌ای",
		},
		{
			CategoryID:      catMen.ID,
			Name:            "شلوار کتان اسپرت مردانه",
			ImageURL:        "https://images.unsplash.com/photo-1624378439575-d8705ad7ae80",
			Description:     "شلوار دمپا راسته با تن‌خور فوق‌العاده و دوخت تمیز صنعتی.",
			Price:           750000,
			DiscountPercent: 0,
			IsSpecialOffer:  false,
			Weight:          "450 گرم",
			Color:           "خاکی",
		},
		{
			CategoryID:      catWomen.ID,
			Name:            "مانتو کتی شیک زنانه",
			ImageURL:        "https://images.unsplash.com/photo-1515886657613-9f3515b0c78f",
			Description:     "طراحی مدرن و لایه‌دوزی شده، مناسب جلسات رسمی و استفاده روزمره.",
			Price:           1200000,
			DiscountPercent: 35,
			IsSpecialOffer:  true, // پیشنهاد شگفت انگیز
			Weight:          "600 گرم",
			Color:           "مشکی",
		},
		{
			CategoryID:      catDigital.ID,
			Name:            "هدفون بلوتوثی نویزکنسلینگ",
			ImageURL:        "https://images.unsplash.com/photo-1505740420928-5e560c06d30e",
			Description:     "شارژدهی تا ۳۰ ساعت، کیفیت صدای استثنایی و بیس عمیق.",
			Price:           2450000,
			DiscountPercent: 15,
			IsSpecialOffer:  true, // پیشنهاد شگفت انگیز
			Weight:          "280 گرم",
			Color:           "طوسی مات",
		},
	}

	for i := range products {
		db.Create(&products[i])
	}

	log.Println("Database seeded with sample categories and products.")
}
