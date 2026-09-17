package customer

import (
	"regexp"
	"strings"
	"testing"
)

func TestCustomerAgingColumnOrder(t *testing.T) {
	upper := strings.ToUpper(CustomerAgingQuery)
	selIdx := strings.Index(upper, "SELECT")
	fromIdx := strings.Index(upper, "FROM OCRD")
	if selIdx < 0 || fromIdx < 0 || fromIdx <= selIdx {
		t.Fatalf("cannot locate SELECT..FROM OCRD in query")
	}
	selectClause := CustomerAgingQuery[selIdx+len("SELECT") : fromIdx]

	aliasRe := regexp.MustCompile(`(?i)AS\s+"([^"]+)"`)
	matches := aliasRe.FindAllStringSubmatch(selectClause, -1)

	var aliases []string
	for _, m := range matches {
		aliases = append(aliases, m[1])
	}

	want := []string{
		"DocEntry",
		"Invoice Number",
		"Invoice Date",
		"Customer Code",
		"Customer Name",
		"Balance",
		"Tax Number",
		"Due Date",
		"Payment Terms",
		"Sales Employee",
		"Open Orders Amount",
		"Payment Term Group",
		"Outstanding Days",
		"Invoice Amount",
		"Paid Amount",
		"Outstanding Amount",
		"0-30 Days",
		"31-60 Days",
		"61-90 Days",
		"91+ Days",
	}

	if len(aliases) != len(want) {
		t.Fatalf("expected %d select aliases, got %d: %v", len(want), len(aliases), aliases)
	}
	for i := range want {
		if aliases[i] != want[i] {
			t.Fatalf("column %d: expected %q, got %q (full order: %v)", i, want[i], aliases[i], aliases)
		}
	}
}
