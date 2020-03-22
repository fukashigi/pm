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

type param struct {
	path      string
	value     string
	encrypted bool
}

// Params ...
type Params struct {
	pp []param
}

// New ...
func New() (Handle, error) {
	h := Handle{}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return h, fmt.Errorf("cant get aws cfg: %w", err)
	}
	h.SSM = ssm.New(cfg)
	return h, nil
}

// Get ...
func (h Handle) Get(p string) (Params, error) {
	i := &ssm.GetParametersByPathInput{
		Path:      aws.String(p),
		Recursive: aws.Bool(true),
	}
	r := h.SSM.GetParametersByPathRequest(i)
	pg := ssm.NewGetParametersByPathPaginator(r)
	pp := []param{}
	for pg.Next(context.Background()) {
		o := pg.CurrentPage()
		for _, p := range o.Parameters {
			np := param{path: *p.Name, value: *p.Value}
			if p.Type == "SecureString" {
				np.encrypted = true
			}
			pp = append(pp, np)
		}
	}
	if err := pg.Err(); err != nil {
		return Params{}, fmt.Errorf("cant get params: %w", err)
	}
	return Params{pp: pp}, nil
}

// Save ...
func (h Handle) Save(pp Params) error {
	for _, p := range pp.pp { // TODO is there a putparameters?
		typ := "String"
		if p.encrypted {
			typ = "SecureString"
		}
		i := &ssm.PutParameterInput{
			Name:      aws.String(p.path),
			Value:     aws.String(p.value),
			Type:      ssm.ParameterType(typ),
			Overwrite: aws.Bool(true),
		}
		r := h.SSM.PutParameterRequest(i)
		_, err := r.Send(context.Background())
		if err != nil {
			// TODO collect errors
			return fmt.Errorf("cant write param: %w", err)
		}
	}
	return nil
}

// Delete ...
func (h Handle) Delete(pp Params) error {
	for _, p := range pp.pp {
		r := h.SSM.DeleteParameterRequest(&ssm.DeleteParameterInput{
			Name: aws.String(p.path),
		})
		_, err := r.Send(context.Background())
		if err != nil {
			// TODO collect errors
			return fmt.Errorf("cant delete param: %w", err)
		}
	}
	return nil
}
