package infrastructure

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-3-aybukeertekin/domain/model"
)

func Read(path string) {
	jobs := make(chan []string)
	results := make(chan model.Book)

	wg := new(sync.WaitGroup)

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go convertToBookStruct(jobs, results, wg)
	}

	go func() {
		file, _ := os.Open(path)
		defer file.Close()

		lines, _ := csv.NewReader(file).ReadAll()

		for _, line := range lines[1:] {
			jobs <- line
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()
}

func convertToBookStruct(jobs <-chan []string, results chan<- model.Book, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		book := model.Book{}
		pageCount, _ := strconv.Atoi(j[2])
		price, _ := strconv.ParseFloat(j[2], 32)
		stockCount, _ := strconv.Atoi(j[3])
		book.Init(j[0], j[1], pageCount, price, stockCount, j[4])
		book.Print()
		results <- book
	}
}
