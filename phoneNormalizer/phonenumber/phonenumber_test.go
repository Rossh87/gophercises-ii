package phonenumber

import (
	"fmt"
	"reflect"
	"testing"
)

type tc struct {
	given []PhoneRecord
	want  []formatInstruction
}

func TestGetInstructions(t *testing.T) {
	cases := []tc{
		{
			[]PhoneRecord{{1, "(123) 456 7891"}},
			[]formatInstruction{{shouldDelete: false, shouldUpdate: true, formattedValue: "1234567891", id: 1}},
		},
		{
			[]PhoneRecord{{1, "(123) 456 7892"}, {2, "1234567892"}},
			[]formatInstruction{{shouldDelete: false, shouldUpdate: true, formattedValue: "1234567892", id: 1}, {shouldDelete: true, shouldUpdate: false, formattedValue: "1234567892", id: 2}},
		},
	}

	for cn, c := range cases {
		t.Run(fmt.Sprintf("case: %d", cn), func(t *testing.T) {
			got := getInstructions(c.given)

			for i, value := range got {
				if !reflect.DeepEqual(value, c.want[i]) {
					t.Fatalf("Input: %v\nWanted: %v\nGot: %v\n", c.given, c.want[i], value)
				}
			}
		})
	}
}
