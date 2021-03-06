package repository

import (
	server "allincecup-server"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user server.User, role string) (int, int, error)
	GetUser(email string, password string) (server.User, error)
	NewSession(session server.Session) (*server.Session, error)
	GetSessionByRefresh(refresh string) (*server.Session, error)
	DeleteSessionByRefresh(refresh string) error
}

type Category interface {
	GetAll() ([]server.Category, error)
	Update(category server.Category) (int, error)
	Create(category server.Category) (int, error)
	Delete(id int, title string) error
}

type Products interface {
	GetProductById(id int) (server.ProductInfoDescription, error)
	AddProduct(product server.Product, info []server.ProductInfo) (int, error)
	Update(product server.ProductInfoDescription) (int, error)
	Delete(productId int) error
}

type Shopping interface {
	AddToCart(userId int, info server.CartProduct) (float64, error)
	PriceValidation(productId, quantity int) (float64, error)
	GetProductsInCart(userId int) ([]server.CartProduct, error)
	AddToFavourites(userId, productId int) error
	GetFavourites(userId int) ([]server.Product, error)
}

type Repository struct {
	Authorization
	Category
	Products
	Shopping
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Products:      NewProductsPostgres(db),
		Category:      NewCategoryPostgres(db),
		Shopping:      NewShoppingPostgres(db),
	}
}
