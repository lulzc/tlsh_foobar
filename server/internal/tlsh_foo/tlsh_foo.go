package tlsh_foo

import (
	"encoding/csv"
	"github.com/glaslos/tlsh"
	"io"
	"log"
	"os"
	"strings"
)

type Item struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

// Result todo: review if store result in Item is better?
type Result struct {
	Distance  int    `json:"distance"`
	Signature string `json:"signature"`
}

type Service struct {
	tlshitem []Item
}

func NewService() *Service {
	return &Service{
		tlshitem: make([]Item, 0),
	}

}

func (svc *Service) Add(name string, hash string) {

	svc.tlshitem = append(svc.tlshitem, Item{
		Name: name,
		Hash: hash,
	})
}

func (svc *Service) Search(query string) []string {
	var results []string
	for _, item := range svc.tlshitem {
		if strings.Contains(item.Hash, query) {
			results = append(results, item.Hash)
		}
	}
	return results
}

func (svc *Service) Distance(query string) []Result {
	var results []Result

	// generate TLSH hash from the query string
	q, err := tlsh.HashBytes([]byte(query))
	if err != nil {
		log.Println("error converting from query to bytes", err)
		return results
	}

	for _, item := range svc.tlshitem {

		// generate TLSH hash for each item
		compareHash, err := tlsh.HashBytes([]byte(item.Hash))
		if err != nil {
			log.Println("error generating hash for item:", err)
			continue
		}

		distance := q.Diff(compareHash)
		// only store results with distance smaller than 150
		if distance < 150 {
			results = append(results, Result{Distance: distance, Signature: item.Signature})
		}

	}

	return results
}

func (svc *Service) GetAll() []Item {
	return svc.tlshitem
}

// ReadCsv reads data from MalwareBazaar csv file and populates the Service's tlshitem
func (s *Service) ReadCsv(csvFile string) {

	file, err := os.Open(csvFile)
	if err != nil {
		log.Println("Error open csv file:", err)
		return
	}
	// reading - defer file.Close() should be fine
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Println("Reading CSV finished...")
				break
			}
			log.Println("CSV Error:", err)
			continue
		}

		// sanitize and remove double quotes to make the json response more pretty
		r := Item{
			Hash:      strings.ReplaceAll(record[13], "\"", ""),
			Name:      strings.ReplaceAll(record[5], "\"", ""),
			Signature: strings.ReplaceAll(record[8], "\"", ""),
		}

		s.tlshitem = append(s.tlshitem, r)
	}
	// for debug print all items from the csv
	//log.Printf("%v\n", s.tlshitem)

}
