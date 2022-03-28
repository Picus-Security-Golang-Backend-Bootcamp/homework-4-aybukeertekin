package repository

import (
	"errors"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/model"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/utils"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) FindAll() []model.Book {
	var books []model.Book
	r.db.Find(&books)
	return books
}

func (r *BookRepository) GetById(id int) (error, model.Book) {
	var book model.Book
	result := r.db.First(&book, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.ErrBookNotFoundWithId, model.Book{}
	}
	return result.Error, book
}

func (r *BookRepository) Create(book *model.Book) error {
	result := r.db.Create(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BookRepository) Update(book model.Book) (error, model.Book) {
	result := r.db.Save(book)
	if result.Error != nil {
		return result.Error, model.Book{}
	}
	return nil, book
}

func (r *BookRepository) Delete(book model.Book) error {
	result := r.db.Delete(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BookRepository) DeleteById(id int) error {
	result := r.db.Delete(&model.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BookRepository) FindBooksByNameOrAuthorNameOrStockCodeOrISBN(searchWord string) []model.Book {
	var books []model.Book
	r.db.Where("Name LIKE ?", "%"+searchWord+"%").Or("StockCode LIKE ?", "%"+searchWord+"%").Or("ISBN LIKE ?", "%"+searchWord+"%").Or("AuthorName LIKE ?", "%"+searchWord+"%").Association("Author").Find(&books)
	return books
}

func (r *BookRepository) Migration() {
	r.db.AutoMigrate(&model.Book{})
}

func (r *AuthorRepository) Migration() {
	r.db.AutoMigrate(&model.Author{})
}

// func (r *CityRepository) InsertSampleData() {
// 	cities := []City{
// 		{Name: "Adana", Code: "01", CountryCode: "TR"},
// 		{Name: "Adıyaman", Code: "02", CountryCode: "TR"},
// 		{Name: "Ankara", Code: "06", CountryCode: "TR"},
// 		{Name: "İstanbul", Code: "34", CountryCode: "TR"},
// 		{Name: "İzmir", Code: "35", CountryCode: "TR"},
// 	}

// 	for _, c := range cities {
// 		r.db.Where(City{Code: c.Code}).Attrs(City{Code: c.Code, Name: c.Name}).FirstOrCreate(&c)
// 	}
// }
