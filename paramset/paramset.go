package paramset

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface"
)

// Services ...
type Services struct {
	SSM ssmiface.ClientAPI
}

// ParamSet ...
type ParamSet struct {
	Svc Services
	pp  []Param
}

// Param ...
type Param struct {
	Path string
	Val  string
	Typ  string
	Ver  string
}

// Len returns the length (cardinality) of the set.
func (p ParamSet) Len() int {
	return len(p.pp)
}

// Contains tests other membership in p.
func (p ParamSet) Contains(other Param) bool {
	for _, this := range p.pp {
		if other.Path == this.Path && other.Value == this.Value {
			return true
		}
	}
	return false
}

// IsSubset tests whether every element in p is in ps.
func (p ParamSet) IsSubset(ps ParamSet) bool {
	return errors.New("not implemented")
}

// IsSuperset tests whether every element in ps is in p.
func (p ParamSet) IsSuperset(ps ParamSet) bool {
	return errors.New("not implemented")
}

// Union returns a new ParamSet with elements from both p and ps.
func (p ParamSet) Union(ps ParamSet) ParamSet {
	return errors.New("not implemented")
}

// Intersection returns a new ParamSet with elements common to p and ps.
func (p ParamSet) Intersection(ps ParamSet) ParamSet {
	return errors.New("not implemented")
}

// Difference returns a new ParamSet with elements in p but not in ps.
func (p ParamSet) Difference(ps ParamSet) ParamSet {
	return errors.New("not implemented")
}

// SymmetricDiff returns a new ParamSet in either p or ps but not both.
func (p ParamSet) SymmetricDiff(ps ParamSet) ParamSet {
	return errors.New("not implemented")
}
