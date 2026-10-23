# داکیومنت کامل APIهای پروژه فروشگاه (Shop Backend API)

این راهنما برای توسعه‌دهنده کلاینت اندروید (یا تست با Postman) تهیه شده است.

- **Base URL:** `http://localhost:8080/api` (یا آدرس سرور اصلی)
- **فرمت داده‌ها:** `application/json`

---

## ۱. نحوه احراز هویت و امنیت (Authentication)

برای بخش‌های خصوصی (سبد خرید و ثبت نظر)، باید توکن JWT را در هدر درخواست ارسال کنید:
```http
Authorization: Bearer <YOUR_TOKEN_HERE>
```

### ۱.۱. درخواست ارسال کد تأیید (Send OTP)
کد یک‌بار مصرف ۶ رقمی به ایمیل ارسال می‌شود (در حالت لوکال در ترمینال سرور لاگ می‌شود).

- **متد:** `POST`
- **آدرس:** `/api/auth/send-otp`
- **Body (JSON):**
```json
{
  "email": "user@gmail.com"
}
```
- **پاسخ (200 OK):**
```json
{
  "message": "Verification code sent to your email"
}
```

---

### ۱.۲. بررسی کد تأیید و دریافت توکن (Verify OTP & Login)
- **متد:** `POST`
- **آدرس:** `/api/auth/verify-otp`
- **Body (JSON):**
```json
{
  "email": "user@gmail.com",
  "code": "729410"
}
```
- **پاسخ (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@gmail.com"
  }
}
```
> **نکته برای اندروید:** فیلد `token` را در `EncryptedSharedPreferences` یا `DataStore` ذخیره کنید تا در درخواست‌های بعدی ارسال شود.

---

## ۲. دسته‌بندی‌ها و محصولات (Categories & Products)

### ۲.۱. دریافت لیست تمام دسته‌بندی‌ها
- **متد:** `GET`
- **آدرس:** `/api/categories`
- **پاسخ (200 OK):**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "پوشاک مردانه",
      "image_url": "https://images.unsplash.com/photo-1490578474895-699cd4e2cf59",
      "description": "انواع لباس، شلوار و کت مردانه شیک"
    },
    {
      "id": 2,
      "name": "پوشاک زنانه",
      "image_url": "https://images.unsplash.com/photo-1483985988355-763728e1935b",
      "description": "جدیدترین مدل‌های پوشاک زنانه"
    }
  ]
}
```

---

### ۲.۲. دریافت محصولات یک دسته‌بندی (با نام یا شناسه)
می‌توانید دقیقاً نام دسته (یا بخشی از آن) یا شناسه عددی را بفرستید:
- **متد:** `GET`
- **آدرس با نام:** `/api/categories/مردانه/products`
- **آدرس با شناسه:** `/api/categories/1/products`
- **پاسخ (200 OK):**
```json
{
  "category": "پوشاک مردانه",
  "products": [
    {
      "id": 1,
      "category_id": 1,
      "name": "پیراهن نخی آستین بلند مردانه",
      "image_url": "https://images.unsplash.com/photo-1602810318383-e386cc2a3ccf",
      "description": "تهیه شده از پارچه ۱۰۰٪ پنبه، بسیار خنک و شیک...",
      "price": 580000,
      "discount_percent": 20,
      "is_special_offer": true,
      "weight": "250 گرم",
      "color": "سرمه‌ای"
    }
  ]
}
```

---

### ۲.۳. دریافت پیشنهادهای شگفت‌انگیز (تخفیف‌دار)
- **متد:** `GET`
- **آدرس:** `/api/products/special-offers`
- **پاسخ (200 OK):**
```json
{
  "special_offers": [
    {
      "id": 1,
      "name": "پیراهن نخی آستین بلند مردانه",
      "price": 580000,
      "discount_percent": 20,
      "is_special_offer": true,
      "image_url": "https://images.unsplash.com/...",
      "color": "سرمه‌ای",
      "weight": "250 گرم"
    }
  ]
}
```

---

### ۲.۴. دریافت جزئیات یک محصول
- **متد:** `GET`
- **آدرس:** `/api/products/1`
- **پاسخ (200 OK):**
```json
{
  "product": {
    "id": 1,
    "name": "پیراهن نخی آستین بلند مردانه",
    "image_url": "https://...",
    "description": "توضیحات کامل محصول...",
    "price": 580000,
    "discount_percent": 20,
    "is_special_offer": true,
    "weight": "250 گرم",
    "color": "سرمه‌ای",
    "category": {
      "id": 1,
      "name": "پوشاک مردانه"
    },
    "comments": []
  }
}
```

---

## ۳. نظرات و کامنت‌ها (Comments)

### ۳.۱. دریافت نظرات یک محصول
- **متد:** `GET`
- **آدرس:** `/api/products/1/comments`
- **پاسخ (200 OK):**
```json
{
  "comments": [
    {
      "id": 1,
      "product_id": 1,
      "content": "جنس پارچه بسیار خنک و سبک بود.",
      "rating": 5,
      "user": {
        "id": 1,
        "email": "user@gmail.com"
      },
      "created_at": "2026-10-22T21:40:00Z"
    }
  ]
}
```

---

### ۳.۲. ثبت نظر جدید (نیاز به توکن JWT)
- **متد:** `POST`
- **آدرس:** `/api/products/1/comments`
- **هدر:** `Authorization: Bearer <TOKEN>`
- **Body (JSON):**
```json
{
  "content": "کیفیت دوختش عالی بود و سریع ارسال شد.",
  "rating": 5
}
```
- **پاسخ (201 Created):**
```json
{
  "message": "Comment submitted successfully",
  "comment": {
    "id": 2,
    "product_id": 1,
    "content": "کیفیت دوختش عالی بود و سریع ارسال شد.",
    "rating": 5,
    "user_id": 1
  }
}
```

---

## ۴. سبد خرید (Cart) - نیازمند لاگین با توکن JWT

همه درخواست‌های این بخش نیاز به هدر احراز هویت دارند:
`Authorization: Bearer <TOKEN>`

### ۴.۱. مشاهده محتویات سبد خرید
- **متد:** `GET`
- **آدرس:** `/api/cart`
- **پاسخ (200 OK):**
```json
{
  "count": 1,
  "total_price": 464000,
  "items": [
    {
      "id": 1,
      "user_id": 1,
      "product_id": 1,
      "quantity": 1,
      "product": {
        "id": 1,
        "name": "پیراهن نخی آستین بلند مردانه",
        "price": 580000,
        "discount_percent": 20,
        "image_url": "https://..."
      }
    }
  ]
}
```

---

### ۴.۲. افزودن محصول به سبد خرید
اگر محصول از قبل در سبد موجود باشد، تعداد آن به صورت خودکار افزایش پیدا می‌کند.

- **متد:** `POST`
- **آدرس:** `/api/cart/items`
- **Body (JSON):**
```json
{
  "product_id": 1,
  "quantity": 1
}
```
- **پاسخ (201 Created یا 200 OK):**
```json
{
  "message": "Item added to cart",
  "cart_item": {
    "id": 1,
    "user_id": 1,
    "product_id": 1,
    "quantity": 1
  }
}
```

---

### ۴.۳. حذف آیتم از سبد خرید
- **متد:** `DELETE`
- **آدرس:** `/api/cart/items/1` (عدد آخر شناسه ردیف سبد خرید `id` است)
- **پاسخ (200 OK):**
```json
{
  "message": "Item removed from cart successfully"
}
```

---

## نمونه پیاده‌سازی سرویس در اندروید (Retrofit / Kotlin)

```kotlin
interface ShopApiService {

    // Auth
    @POST("api/auth/send-otp")
    suspend fun sendOtp(@Body request: SendOtpRequest): Response<MessageResponse>

    @POST("api/auth/verify-otp")
    suspend fun verifyOtp(@Body request: VerifyOtpRequest): Response<AuthResponse>

    // Categories & Products
    @GET("api/categories")
    suspend fun getCategories(): Response<CategoryListResponse>

    @GET("api/categories/{identifier}/products")
    suspend fun getProductsByCategory(@Path("identifier") identifier: String): Response<CategoryProductsResponse>

    @GET("api/products/special-offers")
    suspend fun getSpecialOffers(): Response<SpecialOffersResponse>

    @GET("api/products/{id}")
    suspend fun getProductDetails(@Path("id") id: Int): Response<ProductDetailsResponse>

    // Comments
    @GET("api/products/{id}/comments")
    suspend fun getComments(@Path("id") productId: Int): Response<CommentsResponse>

    @POST("api/products/{id}/comments")
    suspend fun addComment(
        @Header("Authorization") token: String,
        @Path("id") productId: Int,
        @Body request: AddCommentRequest
    ): Response<CommentResponse>

    // Cart
    @GET("api/cart")
    suspend fun getCart(@Header("Authorization") token: String): Response<CartResponse>

    @POST("api/cart/items")
    suspend fun addToCart(
        @Header("Authorization") token: String,
        @Body request: AddToCartRequest
    ): Response<CartItemResponse>

    @DELETE("api/cart/items/{id}")
    suspend fun removeFromCart(
        @Header("Authorization") token: String,
        @Path("id") itemId: Int
    ): Response<MessageResponse>
}
```
