package tables

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func TransactDelete(ctx context.Context, ddb *dynamodb.Client, shop *ShopType, instanceRecord *InstanceType) error {
	instance := shop.Instance
	newVersion := uuid.New().String()
	instancePut, err := putNewInstanceRecord(instance, instanceRecord.Streams - 1, newVersion, instanceRecord.Version)
	if err != nil {
		return err
	}

	instancePortDelete, err := deleteInstancePort(instance, shop.Port)
	if err != nil {
		return err
	}

	streamNameDelete, err := deleteStreamName(shop.Stream)
	if err != nil {
		return err
	}

	shopDelete, err := deleteShopRecord(shop.ShopId, shop.Version)
	if err != nil {
		return err
	}

	transactItems := []types.TransactWriteItem {
		types.TransactWriteItem { Put: instancePut },
		types.TransactWriteItem { Delete: instancePortDelete },
		types.TransactWriteItem { Delete: streamNameDelete },
		types.TransactWriteItem { Delete: shopDelete },
	}

	input := dynamodb.TransactWriteItemsInput {
		TransactItems: transactItems,
		ClientRequestToken: &newVersion,
	}

	_, err = ddb.TransactWriteItems(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}

func deleteShopRecord(shopId string, version string) (*types.Delete, error) {
	vexpr := expression.Equal(
		expression.Name(*Shops.Version.AttributeName),
		expression.Value(version))
	expr, err := expression.NewBuilder().WithCondition(vexpr).Build()
	if err != nil {
		return nil, fmt.Errorf("Unable to create expression for instance key [%v]", err)
	}

	object := struct {
		ShopId string
	}{
		ShopId: shopId,
	}

	key, err := attributevalue.MarshalMap(object)
	if err != nil {
		return nil, err
	}

	delete := types.Delete {
		TableName: Shops.TableName,
		Key: key,
		ConditionExpression: expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	return &delete, nil
}

func deleteInstancePort(instance string, port uint16) (*types.Delete, error) {
	instancePortObj := InstancePortType {
		Instance: instance,
		Port: port,
	}

	return deleteItem(instancePortObj, InstancePorts.TableName)
}

func deleteStreamName(stream string) (*types.Delete, error) {
	streamObj := StreamType { Stream: stream }
	return deleteItem(streamObj, StreamNames.TableName)
}

func deleteItem(object interface{}, table *string) (*types.Delete, error) {
	key, err := attributevalue.MarshalMap(object)
	if err != nil {
		return nil, err
	}

	delete := types.Delete {
		TableName: table,
		Key: key,
	}

	return &delete, nil
}
