package paramset

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface"
)

// Services ...
type Services struct {
	SSM ssmiface.ClientAPI
}

// ParamSet ...
type ParamSet struct {
	S  Services
	pp []Param
}

// Param ...
type Param struct {
	Path string
	Val  string
	Typ  string
	Ver  string
}

// New returns a new, empty ParamSet with a default Services struct.
// If anything goes wrong setting up the Services struct, the empty ParamSet
// is returned _without_ a Services struct and the error is dropped.
func New() ParamSet {
	p := ParamSet{pp: []Param{}}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return p
	}
	p.S = Services{SSM: ssm.New(cfg)}
	return p
}

// Len returns the length (cardinality) of the set.
func (p ParamSet) Len() int {
	return len(p.pp)
}

// Equals returns true if the sets are equal.
func (p ParamSet) Equals(ps ParamSet) bool {
	return p.IsSubset(ps) && p.IsSuperset(ps)
}

// Contains tests other membership in p.
func (p ParamSet) Contains(other Param) bool {
	for _, this := range p.pp {
		if other.Path == this.Path && other.Val == this.Val {
			return true
		}
	}
	return false
}

// IsSubset tests whether every element in p is in ps.
func (p ParamSet) IsSubset(ps ParamSet) bool {
	for _, this := range p.pp {
		if !ps.Contains(this) {
			return false
		}
	}
	return true
}

// IsSuperset tests whether every element in ps is in p.
func (p ParamSet) IsSuperset(ps ParamSet) bool {
	return ps.IsSubset(p)
}

// Union returns a new ParamSet with elements from both p and ps.
func (p ParamSet) Union(ps ParamSet) ParamSet {
	for _, this := range p.pp {
		if !ps.Contains(this) {
			ps.pp = append(ps.pp, this)
		}
	}
	return ps
}

// Intersection returns a new ParamSet with elements common to p and ps.
func (p ParamSet) Intersection(ps ParamSet) ParamSet {
	np := ParamSet{pp: []Param{}}
	for _, this := range p.pp {
		if ps.Contains(this) {
			np.pp = append(np.pp, this)
		}
	}
	return np
}

// Difference returns a new ParamSet with elements in p but not in ps.
func (p ParamSet) Difference(ps ParamSet) ParamSet {
	np := ParamSet{pp: []Param{}}
	for _, this := range p.pp {
		if !ps.Contains(this) {
			np.pp = append(np.pp, this)
		}
	}
	return np
}

// SymmetricDiff returns a new ParamSet in either p or ps but not both.
func (p ParamSet) SymmetricDiff(ps ParamSet) ParamSet {
	np := ParamSet{pp: []Param{}}
	for _, this := range p.pp {
		if ps.Contains(this) {
			continue
		}
		np.pp = append(np.pp, this)
	}
	for _, other := range ps.pp {
		if p.Contains(other) {
			continue
		}
		np.pp = append(np.pp, other)
	}
	return np
}
