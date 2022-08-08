package phonenumber

import (
	"log"
	"strconv"
)

type DB interface {
	GetRecords() ([]PhoneRecord, error)
	Update(id int, value string) error
	Delete(id int) error
}

type PhoneRecord struct {
	id int
	pn string
}

type seenPn map[string]struct{}

type formatInstruction struct {
	id             int
	formattedValue string
	shouldDelete   bool
	shouldUpdate   bool
}

func formatOne(pn string) string {
	var ret []int32

	for _, rn := range pn {
		_, err := strconv.ParseInt(string(rn), 10, 8)

		if err != nil {
			continue
		}

		ret = append(ret, rn)
	}

	return string(ret)
}

func getInstructions(pns []PhoneRecord) []formatInstruction {
	seen := make(seenPn)

	var ret []formatInstruction

	for _, record := range pns {
		formatted := formatOne(record.pn)

		if _, ok := seen[formatted]; ok {
			// we've already seen this, so mark it for deletion
			ret = append(ret, formatInstruction{id: record.id, formattedValue: formatted, shouldDelete: true, shouldUpdate: false})
			continue
		}

		seen[formatted] = struct{}{}

		if formatted == record.pn {
			// no db ops needed--string is correctly formatted and not yet seen,
			// so just add it to 'seen' for checking future entries
			continue
		}

		// otherwise formatting resulted in a change, so add an update instruction to list
		ret = append(ret, formatInstruction{id: record.id, formattedValue: formatted, shouldDelete: false, shouldUpdate: true})
	}

	return ret
}

func Format(db DB) {
	records, err := db.GetRecords()

	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	instructions := getInstructions(records)

	for _, instruction := range instructions {
		if instruction.shouldUpdate {
			err := db.Update(instruction.id, instruction.formattedValue)

			if err != nil {
				log.Fatalf("%+v\n", err)
			}

			continue
		}

		if instruction.shouldDelete {
			err := db.Delete(instruction.id)

			if err != nil {
				log.Fatalf("%+v\n", err)
			}

			continue
		}
	}
}
