package tlsh_foo

import (
	"encoding/csv"
	"github.com/glaslos/tlsh"
	"io"
	"log"
	"os"
	"strings"
	"tlsh_foobar/server/internal/models"
)

type Service struct {
	tlshitem []models.Item
}

func NewService() *Service {
	return &Service{
		tlshitem: make([]models.Item, 0),
	}

}

func (svc *Service) Add(name string, hash string, signature string) {

	svc.tlshitem = append(svc.tlshitem, models.Item{
		Name:      name,
		Hash:      hash,
		Signature: signature,
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

func (svc *Service) Distance(query string) []models.Result {
	var results []models.Result

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
			results = append(results, models.Result{Distance: distance, Signature: item.Signature})
		}

	}
	return results
}

func (svc *Service) GetAll() []models.Item {
	return svc.tlshitem
}

// reads the data from MalwareBazaar csv file
func (svc *Service) ReadCsv(csvFile string) {

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

		r := models.Item{
			Name:      strings.ReplaceAll(record[5], "\"", ""),
			Signature: strings.ReplaceAll(record[8], "\"", ""),
			Hash:      strings.ReplaceAll(record[13], "\"", ""),
		}

		svc.tlshitem = append(svc.tlshitem, r)

	}

	// for debug print all items from the csv in the console
	// log.Printf("%v\n", s.tlshitem)

}
