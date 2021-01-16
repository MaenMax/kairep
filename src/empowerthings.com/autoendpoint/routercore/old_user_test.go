package routercore

import(
//	"fmt"
//	"strings"
	"testing"
	"time"
)

var NOW time.Time = time.Unix(31536000, 0)   //  Friday, January 1, 1971 00:00:00

func Test_CheckIfOld(t *testing.T){
	var test_cases = []struct {
		Now              time.Time
		Year             int
		Month            int
		Expected         bool
	}{
		{
			Year:            1971,
			Month:           01,
			Now:             NOW,
			Expected:        false,
		},
		{
			Year:            1970,
			Month:           12,
			Now:             NOW,
			Expected:        false,
		},
		{
			Year:            1970,
			Month:           11,
			Now:             NOW,
			Expected:        false,
		},
		{
			Year:            1970,
			Month:           10,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           9,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           8,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           7,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           6,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           5,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           4,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           3,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           2,
			Now:             NOW,
			Expected:        true,
		},
		{
			Year:            1970,
			Month:           1,
			Now:             NOW,
			Expected:        true,
		},
	}

	for i,tc:= range(test_cases){
		result:=CheckIfOld(tc.Year, tc.Month, tc.Now)
		if result!=tc.Expected {
			t.Errorf("Test case# %v failed. Expected: %v Got: %v", i+1, tc.Expected, result)
			return
		}	
	}	
}

