package txt

import (
	"github.com/766b/mobi"
	"testing"
)

func Test_downCover3(t *testing.T) {
	downCover()
}
func Test_downCover(t *testing.T) {
	//downCover()
	m, err := mobi.NewWriter("output-test.mobi")
	if err != nil {
		panic(err)
	}

	m.Title("Book Title")
	m.Compression(mobi.CompressionNone) // LZ77 compression is also possible using  mobi.CompressionPalmDoc

	// Add cover image
	path := downCover()
	m.AddCover(path, path)

	// Meta data
	m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi.EXTH_AUTHOR, "Book Author Name")
	// See exth.go for additional EXTH record IDs

	// Add chapters and subchapters
	ch1 := m.NewChapter("Chapter 1", []byte("Some text here <br style=\"line-height: 220%\">    sdf br <br style=\"line-height: 26px\">  sdf"))
	ch1.AddSubChapter("Chapter 1-1", []byte("&nbsp;&nbsp;&nbsp;&nbsp;Some text here<br>"+
		"&nbsp;&nbsp;&nbsp;&nbsp; ss sdf  \n"))

	// Output MOBI File
	m.Write()
}
