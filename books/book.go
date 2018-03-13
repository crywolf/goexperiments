package books

// Book type
type Book struct {
	Title  string
	Author string
	Pages  int
}

// CategoryByLength determines category by length of the book
func (b Book) CategoryByLength() string {
	if b.Pages > 300 {
		return "NOVEL"
	}
	return "SHORT STORY"

}

func someFunc() {
}
