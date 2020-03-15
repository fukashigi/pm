// Package store is an abstraction for AWS Parameter Store.
package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface"
)

// Handle is a service handle for store. Use it to do things.
type Handle struct {
	SSM ssmiface.ClientAPI
}

// Param is a single AWS Parameter Store Parameter ;)
type Param struct {
	Path      string
	Value     string
	Encrypted bool
}

// New returns a functioning Handle and an error. If the error is not nil, the
// Handle probably doesn't work.
func New() (Handle, error) {
	h := Handle{}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return h, fmt.Errorf("cant get aws cfg: %w", err)
	}
	h.SSM = ssm.New(cfg)
	return h, nil
}

// Get takes a path and returns the parameter at that location as a Param and
// an error.
func (h Handle) Get(path string, decrypt bool) (Param, error) {
	r := h.SSM.GetParameterRequest(&ssm.GetParameterInput{
		Name: aws.String(path),
	})
	o, err := r.Send(context.Background())
	if err != nil {
		return Param{}, fmt.Errorf("cant get param: %w", err)
	}
	p := Param{
		Path:  path,
		Value: *o.Parameter.Value,
	}
	if o.Parameter.Type == "SecureString" {
		p.Encrypted = true
	}
	return p, nil
}

// Gets ...
func (h Handle) Gets(p string) ([]Param, error) {
	i := &ssm.GetParametersByPathInput{
		Path:      aws.String(p),
		Recursive: aws.Bool(true),
	}
	r := h.SSM.GetParametersByPathRequest(i)
	pg := ssm.NewGetParametersByPathPaginator(r)
	pp := []Param{}
	for pg.Next(context.Background()) {
		o := pg.CurrentPage()
		for _, p := range o.Parameters {
			np := Param{Path: *p.Name, Value: *p.Value}
			if p.Type == "SecureString" {
				np.Encrypted = true
			}
			pp = append(pp, np)
		}
	}
	if err := pg.Err(); err != nil {
		return pp, fmt.Errorf("cant get params: %w", err)
	}
	return pp, nil
}

// Put ...
func (h Handle) Put(p, v string, encrypt bool) error {
	typ := "String"
	if encrypt {
		typ = "SecureString"
	}
	i := &ssm.PutParameterInput{
		Name:      aws.String(p),
		Value:     aws.String(v),
		Type:      ssm.ParameterType(typ),
		Overwrite: aws.Bool(true),
	}
	r := h.SSM.PutParameterRequest(i)
	_, err := r.Send(context.Background())
	if err != nil {
		return fmt.Errorf("cant write param: %w", err)
	}
	return nil
}

// Del ...
func (h Handle) Del(p string) error {
	r := h.SSM.DeleteParameterRequest(&ssm.DeleteParameterInput{
		Name: aws.String(p),
	})
	_, err := r.Send(context.Background())
	if err != nil {
		return fmt.Errorf("cant delete param: %w", err)
	}
	return nil
}
