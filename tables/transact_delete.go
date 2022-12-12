package tables

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TransactDelete(ctx context.Context, ddb *dynamodb.Client, shop *ShopType, instanceRecord *InstanceType) error {
	return fmt.Errorf("Not Implemented")
}
