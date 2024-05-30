package main

import (
	"testing"
)

func TestReadContactsFromXLSX(t *testing.T) {
	filePath := "./test_contacts.xlsx" // Update with the path to your test XLSX file
	contacts, err := readContactsFromXLSX(filePath)
	if err != nil {
		t.Errorf("Error reading contacts from XLSX: %v", err)
	}

	// Check if the number of contacts read is correct
	expectedNumContacts := 20 // Update with the expected number of contacts in your test XLSX file
	if len(contacts) != expectedNumContacts {
		t.Errorf("Expected %d contacts, got %d", expectedNumContacts, len(contacts))
	}

	// Add more assertions to validate the contents of the contacts if needed
}

func TestCalculateMatchScore(t *testing.T) {
	// Test cases for the calculateMatchScore function
	testCases := []struct {
		name          string
		contact1      Contact
		contact2      Contact
		expectedScore int
	}{
		// Define test cases here
		{
			name: "Same contacts",
			contact1: Contact{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				ZipCode:   "12345",
				Address:   "123 Main St",
			},
			contact2: Contact{
				ID:        2,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				ZipCode:   "12345",
				Address:   "123 Main St",
			},
			expectedScore: 1000,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := calculateMatchScore(tc.contact1, tc.contact2)
			if score != tc.expectedScore {
				t.Errorf("Expected score %d, got %d", tc.expectedScore, score)
			}
		})
	}
}
