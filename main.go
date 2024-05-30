package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
)

type Contact struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	ZipCode   string
	Address   string
}

type ContactMap struct {
	contacts map[string]Contact
}

func NewContactMap(contacts []Contact) (*ContactMap, error) {
	// Time complexity: O(n), where n is the number of contacts
	// Space complexity: O(n), where n is the number of contacts
	cm := &ContactMap{contacts: make(map[string]Contact)}
	for _, contact := range contacts {
		if _, exists := cm.contacts[fmt.Sprintf("%s-%s-%s-%s-%s", contact.FirstName, contact.LastName, contact.Email, contact.ZipCode, contact.Address)]; exists {
			return nil, fmt.Errorf("contact with ID %d already exists", contact.ID)
		}
		cm.contacts[fmt.Sprintf("%s-%s-%s-%s-%s", contact.FirstName, contact.LastName, contact.Email, contact.ZipCode, contact.Address)] = contact
	}
	return cm, nil
}

func (cm *ContactMap) GetContact(key string) (Contact, bool) {
	// Time complexity: O(1), where n is the number of contacts
	// Space complexity: O(1), where n is the number of contacts
	if contact, exists := cm.contacts[key]; exists {
		return contact, true
	}
	return Contact{}, false
}

// We can optimize the code by using goroutines to calculate the similarity scores concurrently.
// This can help reduce the overall time taken to calculate the similarity scores for all pairs of contacts.
func (cm *ContactMap) GetSimilarityScores() map[string]int {
	// Time complexity: O(n^2), where n is the number of contacts
	// Space complexity: O(n), where n is the number of contact pairs (n*(n-1)/2)
	contactList := make([]Contact, 0, len(cm.contacts))
	for _, contact := range cm.contacts {
		contactList = append(contactList, contact)
	}

	// Using a channel to collect similarity scores
	scores := make(map[string]int)
	scoreChan := make(chan map[string]int)
	var wg sync.WaitGroup

	for i := 0; i < len(contactList); i++ {
		for j := i + 1; j < len(contactList); j++ {
			wg.Add(1)
			go func(contact1, contact2 Contact) {
				defer wg.Done()
				score := calculateMatchScore(contact1, contact2)
				key := fmt.Sprintf("%d-%d", contact1.ID, contact2.ID)
				scoreChan <- map[string]int{key: score}
			}(contactList[i], contactList[j])
		}
	}

	go func() {
		wg.Wait()
		close(scoreChan)
	}()

	for score := range scoreChan {
		for k, v := range score {
			scores[k] = v
		}
	}

	return scores
}

func calculateMatchScore(contact1, contact2 Contact) int {
	// Time complexity: O(1), where n is the number of fields in a contact
	// Space complexity: O(1), where n is the number of fields in a contact
	score := 0.0
	if contact1.FirstName == contact2.FirstName {
		score += 2
	} else {
		score += similarity(contact1.FirstName, contact2.FirstName) * 2
	}

	if contact1.LastName == contact2.LastName {
		score += 2
	} else {
		score += similarity(contact1.LastName, contact2.LastName) * 2
	}

	if contact1.Email == contact2.Email {
		score += 3
	} else {
		score += similarity(contact1.Email, contact2.Email) * 3
	}

	if contact1.ZipCode == contact2.ZipCode {
		score += 2
	} else {
		score += similarity(contact1.ZipCode, contact2.ZipCode) * 2
	}

	if contact1.Address == contact2.Address {
		score += 2
	} else {
		score += similarity(contact1.Address, contact2.Address) * 2
	}

	if contact1.ID == 1001 && contact2.ID == 1000 {
		fmt.Println("Contact1: ", contact1)
		fmt.Println("Contact2: ", contact2)
		fmt.Println("Score: ", score)
	}

	return int(score * 1000 / 11)
}

func similarity(s1, s2 string) float64 {
	// Time complexity: O(mn), where m and n are the lengths of the input strings
	// Space complexity: O(mn), where m and n are the lengths of the input strings
	return 1 - float64(levenshteinDistance(s1, s2))/math.Max(float64(len(s1)), float64(len(s2)))
}

func levenshteinDistance(s1, s2 string) int {
	// Time complexity: O(mn), where m and n are the lengths of the input strings
	// Space complexity: O(mn), where m and n are the lengths of the input strings
	d := make([][]int, len(s1)+1)
	for i := range d {
		d[i] = make([]int, len(s2)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(s2); j++ {
		for i := 1; i <= len(s1); i++ {
			if s1[i-1] == s2[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				d[i][j] = 1 + minCalc(d[i-1][j], d[i][j-1], d[i-1][j-1])
			}
		}
	}
	return d[len(s1)][len(s2)]
}

func minCalc(a, b, c int) int {
	// Time complexity: O(1), where n is the number of inputs
	// Space complexity: O(1), where n is the number of inputs
	return int(math.Min(float64(min(a, b)), float64(c)))
}

func readContactsFromXLSX(filePath string) ([]Contact, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	for _, row := range rows[1:] { // Skipping the header row
		id, _ := strconv.Atoi(row[0])

		contact := Contact{
			ID:        id,
			FirstName: "",
			LastName:  "",
			Email:     "",
			ZipCode:   "",
			Address:   "",
		}

		if len(row) > 1 {
			contact.FirstName = row[1]
		}
		if len(row) > 2 {
			contact.LastName = row[2]
		}
		if len(row) > 3 {
			contact.Email = row[3]
		}
		if len(row) > 4 {
			contact.ZipCode = row[4]
		}
		if len(row) > 5 {
			contact.Address = row[5]
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func writeScoresToXLSX(filePath string, scores map[string]int) error {
	f := excelize.NewFile()
	_, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	err = f.SetCellValue("Sheet1", "A1", "ContactID1")
	if err != nil {
		return err
	}
	err = f.SetCellValue("Sheet1", "B1", "ContactID2")
	if err != nil {
		return err
	}
	err = f.SetCellValue("Sheet1", "C1", "SimilarityScore")
	if err != nil {
		return err
	}

	row := 2
	for key, score := range scores {
		parts := strings.Split(key, "-")
		err = f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), parts[0])
		if err != nil {
			return err
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), parts[1])
		if err != nil {
			return err
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), score)
		if err != nil {
			return err
		}
		row++
	}

	if err := f.SaveAs(filePath); err != nil {
		return err
	}
	return nil
}

func calculateScore(inputFile, outputFile string) error {
	contacts, err := readContactsFromXLSX(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read contacts from XLSX: %w", err)
	}

	cm, err := NewContactMap(contacts)
	if err != nil {
		return fmt.Errorf("failed to create contact map: %w", err)
	}

	scores := cm.GetSimilarityScores()

	if err := writeScoresToXLSX(outputFile, scores); err != nil {
		return fmt.Errorf("failed to write scores to XLSX: %w", err)
	}

	fmt.Println("Similarity scores written to", outputFile)
	return nil
}

func main() {
	// if I am in container utility env, I take the env vars from the system
	inputFile := os.Getenv("INPUT_FILE")
	outputFile := os.Getenv("OUTPUT_FILE")
	if inputFile != "" && outputFile != "" {
		err := calculateScore(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	app := &cli.App{
		Name:  "Contact Similarity",
		Usage: "Calculate similarity scores for contacts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Input XLSX file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output XLSX file",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			inputFile := c.String("file")
			outputFile := c.String("output")
			err := calculateScore(inputFile, outputFile)
			if err != nil {
				return err
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "help",
				Usage: "Shows help",
				Action: func(c *cli.Context) error {
					err := cli.ShowAppHelp(c)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
