package library

import (
	"strings"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-3-aybukeertekin/domain/model"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-3-aybukeertekin/domain/repository"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-3-aybukeertekin/domain/utils"
)

type Library struct {
	bookRepository   *repository.BookRepository
	authorRepository *repository.AuthorRepository
}

func NewLibrary(bookRepository *repository.BookRepository, authorRepository *repository.AuthorRepository) *Library {
	return &Library{
		bookRepository:   bookRepository,
		authorRepository: authorRepository,
	}
}

//This function checks if the book exists,
//then check stock count
//then decreases stock
func (library *Library) BuyBooks(id int, count int) (error, model.Book) {
	err, book := library.bookRepository.GetById(id)
	if err != nil {
		if err == utils.ErrBookNotFoundWithId {
			return utils.ErrBookNotFoundWithId, model.Book{}
		} else {
			return err, model.Book{}
		}
	} else {
		err := book.Buy(count)
		if err != nil {
			return err, model.Book{}
		} else {
			return library.bookRepository.Update(book)
		}
	}
}

func (library *Library) CreateBook(book *model.Book) error {
	return library.bookRepository.Create(book)
}

//This functions prints book by checking their availability
func (library *Library) ListBooks() []model.Book {
	return library.bookRepository.FindAll()
}

func (library *Library) DeleteBook(bookId int) error {
	return library.bookRepository.DeleteById(bookId)
}

//This function prints book list.
func PrintBooks(bookList []model.Book) {
	for index := 0; index < len(bookList); index++ {
		bookList[index].Print()
	}
}

//This function searchs book list and look if a book exists
//with the given book name, isbn, or stock code and prints them.
//Prints error message if no book exists.
func (library *Library) SearchBooks(searchWord string) {
	searchWord = strings.ToLower(searchWord)
	bookList := library.bookRepository.FindBooksByNameOrAuthorNameOrStockCodeOrISBN(searchWord)
	if len(bookList) == 0 {
		utils.PrintLineMessage("\tNo books are found with given argument.")
	} else {
		utils.PrintLineMessage("Found Books:")
		PrintBooks(bookList)
	}
}
