package books_test

import (
	//	. "github.com/crywolf/goexperiments/books"

	"github.com/crywolf/goexperiments/books"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Book", func() {
	var (
		longBook  books.Book
		shortBook books.Book
	)

	BeforeEach(func() {
		longBook = books.Book{
			Title:  "Bídníci",
			Author: "Victor Hugo",
			Pages:  1488,
		}

		shortBook = books.Book{
			Title:  "Fox In Socks",
			Author: "Dr. Seuss",
			Pages:  24,
		}
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})

		Context("With less than 300 pages", func() {
			It("should be a short story", func() {
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})
	})
})
