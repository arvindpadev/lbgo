package tables

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var TestMockAwsCfg *aws.Config = nil

type Ctxt interface {
	Ctx() context.Context
	Cfg() *aws.Config
}

type ctxt struct {
	ctx context.Context
	cfg *aws.Config
}

func (c ctxt) Ctx() context.Context {
	return c.ctx
}

func (c ctxt) Cfg() *aws.Config {
	return c.cfg
}

func Context() (Ctxt, error) {
	ctx := context.TODO()
	cfg := TestMockAwsCfg
	if cfg == nil {
		c, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("failed to load configuration, %v", err)
			return nil, err
		}

		cfg = &c
	}

	return ctxt {
			ctx: ctx,
			cfg: cfg,
		}, nil
}
