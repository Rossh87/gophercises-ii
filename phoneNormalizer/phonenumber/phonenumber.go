package phonenumber

import (
	"bytes"
	"log"
)

type db interface {
	GetRecords() ([]PhoneRecord, error)
	Update(id int, value string) error
	Delete(id int) error
}

type PhoneRecord struct {
	Id int
	Pn string
}

type seenPn map[string]struct{}

type formatInstruction struct {
	id             int
	formattedValue string
	shouldDelete   bool
	shouldUpdate   bool
}

func formatOne(pn string) string {
	var buf bytes.Buffer

	for _, rn := range pn {
		if rn >= '0' && rn <= '9' {
			buf.WriteRune(rn)
		}
	}

	return buf.String()
}

func getInstructions(Pns []PhoneRecord) []formatInstruction {
	seen := make(seenPn)

	var ret []formatInstruction

	for _, record := range Pns {
		formatted := formatOne(record.Pn)

		if _, ok := seen[formatted]; ok {
			// we've already seen this, so mark it for deletion
			ret = append(ret, formatInstruction{id: record.Id, formattedValue: formatted, shouldDelete: true, shouldUpdate: false})
			continue
		}

		seen[formatted] = struct{}{}

		if formatted == record.Pn {
			// no db ops needed--string is correctly formatted and not yet seen,
			// so just add it to 'seen' for checking future entries
			continue
		}

		// otherwise formatting resulted in a change, so add an update instruction to list
		ret = append(ret, formatInstruction{id: record.Id, formattedValue: formatted, shouldDelete: false, shouldUpdate: true})
	}

	return ret
}

func Format(db db) {
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
