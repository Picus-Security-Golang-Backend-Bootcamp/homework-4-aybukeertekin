package model

import (
	"fmt"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/utils"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name       string
	StockCode  string
	ISBN       string
	PageCount  int
	Price      float64
	StockCount int
	IsDeleted  bool
	AuthorID   int
}

//constructor
func (b *Book) Init(bookName string, isbn string, pageCount int, price float64, stockCount int, stockCode string) {
	b.Name = bookName
	b.StockCode = stockCode
	b.ISBN = isbn
	b.PageCount = pageCount
	b.Price = price
	b.StockCount = stockCount
	b.IsDeleted = false
}

func (book *Book) Buy(count int) error {
	if book.StockCount < count {
		return utils.ErrStockIsNotEnough
	} else {
		book.StockCount -= count
		return nil
	}
}

func (book *Book) Print() {
	fmt.Printf("\tID:\t\t%v\n\tISBN:\t\t%v\n\tBook Name: \t%s\n\tAuthor:\t\t%s\n\tPage Count:\t%v\n\tPrice:\t\t%v\n\tStock Code:\t%s\n\tStock Count:\t%v\t\n\n", book.ID, book.ISBN, book.Name, "ss", book.PageCount, book.Price, book.StockCode, book.StockCount)
}
