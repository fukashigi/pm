package paramset

import (
	"testing"
)

func TestSetOps(t *testing.T) {
	param := Param{Path: "/cfg/alpha", Val: "the first", Typ: "String", Ver: "1"}
	oneParam := []Param{param}
	leftParams := []Param{
		{Path: "/cfg/alpha", Val: "the first", Typ: "String", Ver: "1"},
		{Path: "/cfg/beta", Val: "the second", Typ: "String", Ver: "2"},
		{Path: "/cfg/beta/beth", Val: "the 2nd 2nd", Typ: "String", Ver: "1"},
	}
	rightParams := []Param{
		{Path: "/cfg/aleph", Val: "the first", Typ: "String", Ver: "6"},
		{Path: "/cfg/beta", Val: "the second", Typ: "String", Ver: "4"},
		{Path: "/cfg/beta/bzzt", Val: "the 2nd second", Typ: "String", Ver: "9"},
		{Path: "/cfg/beta/beth", Val: "the 2nd 2nd", Typ: "String", Ver: "11"},
		{Path: "/cfg/alpha", Val: "the 2nd first", Typ: "String", Ver: "3"},
	}
	unionParams := []Param{
		{Path: "/cfg/aleph", Val: "the first", Typ: "String", Ver: "6"},
		{Path: "/cfg/alpha", Val: "the 2nd first", Typ: "String", Ver: "3"},
		{Path: "/cfg/alpha", Val: "the first", Typ: "String", Ver: "1"},
		{Path: "/cfg/beta", Val: "the second", Typ: "String", Ver: "?"},
		{Path: "/cfg/beta/beth", Val: "the 2nd 2nd", Typ: "String", Ver: "?"},
		{Path: "/cfg/beta/bzzt", Val: "the 2nd second", Typ: "String", Ver: "9"},
	}
	interParams := []Param{
		{Path: "/cfg/beta", Val: "the second", Typ: "String", Ver: "2"},
		{Path: "/cfg/beta/beth", Val: "the 2nd 2nd", Typ: "String", Ver: "1"},
	}
	diffParams := []Param{param}
	symDiffParams := []Param{
		{Path: "/cfg/aleph", Val: "the first", Typ: "String", Ver: "6"},
		{Path: "/cfg/alpha", Val: "the 2nd first", Typ: "String", Ver: "3"},
		{Path: "/cfg/alpha", Val: "the first", Typ: "String", Ver: "1"},
		{Path: "/cfg/beta/bzzt", Val: "the 2nd second", Typ: "String", Ver: "9"},
	}

	one := ParamSet{pp: oneParam}
	left := ParamSet{pp: leftParams}
	right := ParamSet{pp: rightParams}

	if !one.IsSubset(left) {
		t.Errorf("failed one.IsSubset(left)")
	}
	if left.IsSubset(right) {
		t.Errorf("failed left.IsSubset(right)")
	}
	if !left.IsSuperset(one) {
		t.Errorf("failed IsSuperset")
	}
	if right.Contains(param) {
		t.Errorf("failed right.Contains(param)")
	}
	u := left.Union(right)
	if !u.Equals(ParamSet{pp: unionParams}) {
		t.Errorf("failed left.Union(right)")
	}
	i := left.Intersection(right)
	if !i.Equals(ParamSet{pp: interParams}) {
		t.Errorf("failed left.Intersection(right)")
	}
	d := left.Difference(right)
	if !d.Equals(ParamSet{pp: diffParams}) {
		t.Errorf("failed left.Difference(right)")
	}
	sd := left.SymmetricDiff(right)
	if !sd.Equals(ParamSet{pp: symDiffParams}) {
		t.Errorf("failed left.SymmetricDiff(right)")
	}
}
