package tables

import (
	"context"
	"fmt"
	"log"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func TestShopIdPresence(ctx context.Context, ddb *dynamodb.Client, shopId string) (bool, error) {
	shopIdKeyMatch := struct {
		ShopId string
	} {
		ShopId: shopId,
	}
	shopIdKeyMatchMap, err := attributevalue.MarshalMap(shopIdKeyMatch)
	if err != nil {
		// true so that we error out even if the err is not considered
		return true, err
	}

	projectionExpression := *Shops.ShopId.AttributeName
	input := dynamodb.GetItemInput {
		TableName: Shops.TableName,
		Key: shopIdKeyMatchMap,
		ProjectionExpression: &projectionExpression,
	}

	output, err := ddb.GetItem(ctx, &input)
	if err != nil {
		log.Println(fmt.Sprintf("INFO: Error getting shop for shop id %v: [%v]", shopId, err))
		return true, err
	}

	if len(output.Item) == 0 {
		return false, nil
	}

	return true, nil
}

func TestStreamPresence(ctx context.Context, ddb *dynamodb.Client, stream string) (bool, error) {
	streamKeyMatch := StreamType {
		Stream: stream,
	}
	streamKeyMatchMap, err := attributevalue.MarshalMap(streamKeyMatch)
	if err != nil {
		// true so that we error out even if the err is not considered
		return true, err
	}

	projectionExpression := "#Stream"
	expressionAttributeNames := map[string]string { "#Stream": *StreamNames.Stream.AttributeName }
	input := dynamodb.GetItemInput {
		TableName: StreamNames.TableName,
		Key: streamKeyMatchMap,
		ProjectionExpression: &projectionExpression,
		ExpressionAttributeNames: expressionAttributeNames,
	}

	output, err := ddb.GetItem(ctx, &input)
	if err != nil {
		log.Println(fmt.Sprintf("INFO: Error getting stream for %v: [%v]", stream, err))
		return true, err
	}

	if len(output.Item) == 0 {
		return false, nil
	}

	return true, nil
}

func ConsistentGetInstance(ctx context.Context, ddb *dynamodb.Client, instance string) (*InstanceType, error) {
	instanceKeyMatch := InstanceNameType {
		Instance: instance,
	}

	instanceKeyMatchMap, err := attributevalue.MarshalMap(instanceKeyMatch)
	if err != nil {
		return nil, err
	}

	consistentRead := true
	projectionExpression := fmt.Sprintf("%v, %v, %v", *Instances.Instance.AttributeName, *Instances.Streams.AttributeName, *Instances.Version.AttributeName)
	input := dynamodb.GetItemInput {
		TableName: Instances.TableName,
		Key: instanceKeyMatchMap,
		ProjectionExpression: &projectionExpression,
		ConsistentRead: &consistentRead,
	}

	output, err := ddb.GetItem(ctx, &input)
	if err != nil {
		log.Println(fmt.Sprintf("INFO: Error getting instance %v: [%v]", instance, err))
		return nil, err
	}

	if len(output.Item) == 0 {
		return nil, fmt.Errorf("Instance %v is absent", instance)
	}

	var instanceRecord InstanceType
	err = attributevalue.UnmarshalMap(output.Item, &instanceRecord)
	if err != nil {
		return nil, err
	}

	return &instanceRecord, nil
}



func ConsistentGetShop(ctx context.Context, ddb *dynamodb.Client, shopId string) (*ShopType, error) {
	shopIdKeyMatch := struct {
		ShopId string
	}{
		ShopId: shopId,
	}

	shopIdKeyMatchMap, err := attributevalue.MarshalMap(shopIdKeyMatch)
	if err != nil {
		return nil, err
	}

	consistentRead := true
	input := dynamodb.GetItemInput {
		TableName: Shops.TableName,
		Key: shopIdKeyMatchMap,
		ConsistentRead: &consistentRead,
	}

	output, err := ddb.GetItem(ctx, &input)
	if err != nil {
		log.Println(fmt.Sprintf("INFO: Error getting shop %v: [%v]", shopId, err))
		return nil, err
	}

	if len(output.Item) == 0 {
		return nil, fmt.Errorf("ShopId %v is absent", shopId)
	}

	var shop ShopType
	err = attributevalue.UnmarshalMap(output.Item, &shop)
	if err != nil {
		return nil, err
	}

	return &shop, nil
}

func GetIps(ctx context.Context, ddb *dynamodb.Client, instance string) (string, string, error) {
	instanceKeyMatch := InstanceNameType {
		Instance: instance,
	}
	instanceKeyMatchMap, err := attributevalue.MarshalMap(instanceKeyMatch)
	if err != nil {
		// true so that we error out even if the err is not considered
		return "", "", err
	}

	consistentRead := true
	projectionExpression := fmt.Sprintf("%v, %v", *InstanceIp.PublicIp.AttributeName, *InstanceIp.PrivateIp.AttributeName)
	input := dynamodb.GetItemInput {
		TableName: InstanceIp.TableName,
		Key: instanceKeyMatchMap,
		ProjectionExpression: &projectionExpression,
		ConsistentRead: &consistentRead,
	}

	output, err := ddb.GetItem(ctx, &input)
	if err != nil {
		log.Println(fmt.Sprintf("INFO: Error getting instance ip info for %v: [%v]", instance, err))
		return "", "", err
	}

	if len(output.Item) == 0 {
		return "", "", fmt.Errorf("Instance ip information for %v is absent", instance)
	}

	var instanceIpRecord InstanceIpType
	err = attributevalue.UnmarshalMap(output.Item, &instanceIpRecord)
	if err != nil {
		return "", "", err
	}

	return instanceIpRecord.PublicIp, instanceIpRecord.PrivateIp, nil
}
