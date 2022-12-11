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

func TransactAddStream(ctx context.Context, ddb *dynamodb.Client, shopId string, stream string, port uint16, instanceRecord *InstanceType) error {
	// NOTE: The same version is reused across tables. However, equality cannot
	// be assumed. Do not rely on equality. Its use for idempotency is also just
	// a convenience, and has no significance besides being a random string that
	// can reasonably be assumed to be ungeneratable again for a long time.
	newVersion := uuid.New().String()
	instancePut, err := putNewInstanceRecord(instanceRecord.Instance, instanceRecord.Streams + 1, newVersion, instanceRecord.Version)
	if err != nil {
		return err
	}

	instancePortPut, err := putNewInstancePort(instanceRecord.Instance, port)
	if err != nil {
		return err
	}

	streamNamePut, err := putNewStreamName(stream)
	if err != nil {
		return err
	}

	shopPut, err := putShopRecord(shopId, stream, port, instanceRecord.Instance, newVersion)
	if err != nil {
		return err
	}

	transactItems := []types.TransactWriteItem {
		types.TransactWriteItem { Put: instancePut },
		types.TransactWriteItem { Put: instancePortPut},
		types.TransactWriteItem { Put: streamNamePut},
		types.TransactWriteItem { Put: shopPut },
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

func putNewInstanceRecord(instance string, streams uint8, newVersion string, oldVersion string) (*types.Put, error) {
	vexpr := expression.Equal(
		expression.Name(*Instances.Version.AttributeName),
		expression.Value(oldVersion))
	expr, err := expression.NewBuilder().WithCondition(vexpr).Build()
	if err != nil {
		return nil, fmt.Errorf("Unable to create expression for instance key [%v]", err)
	}
	
	instanceObj := InstanceType {
		Instance: instance,
		Streams: streams,
		Version: newVersion,
	}

	item, err := attributevalue.MarshalMap(&instanceObj)
	if err != nil {
		return nil, err
	}

	put := types.Put {
		TableName: Instances.TableName,
		Item: item,
		ConditionExpression: expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	return &put, nil
}

func putNewInstancePort(instance string, port uint16) (*types.Put, error) {
	instancePortObj := InstancePortType {
		Instance: instance,
		Port: port,
	}

	return putItem(instancePortObj, InstancePorts.TableName)
}

func putNewStreamName(stream string) (*types.Put, error) {
	streamObj := StreamType { Stream: stream }
	return putItem(streamObj, StreamNames.TableName)
}

func putShopRecord(shopId string, stream string, port uint16, instance string, version string) (*types.Put, error) {
	shopObj := ShopType {
		ShopId: shopId,
		Stream: stream,
		Port: port,
		Instance: instance,
		Version: version,
	}

	return putItem(shopObj, Shops.TableName)
}

func putItem(object interface{}, table *string) (*types.Put, error) {
	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		return nil, err
	}

	put := types.Put {
		TableName: table,
		Item: item,
	}

	return &put, nil
}
