package tables

import (
	"context"
	"fmt"
	"log"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func QueryAllInstancesWithNumStreams(ctx context.Context, ddb *dynamodb.Client, num uint8) (*[]InstanceNameType, error) {
	kexpr := expression.Key(*Instances.Streams.AttributeName).Equal(expression.Value(num))
	expr, err := expression.NewBuilder().WithKeyCondition(kexpr).Build()
	if err != nil {
		return nil, fmt.Errorf("Unable to create expression for query [%v]", err)
	}

	projectionExpression := Instances.Instance.AttributeName
	query := dynamodb.QueryInput {
		TableName: Instances.TableName,
		IndexName: &InstancesGsiStreamsInstance,
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression: expr.KeyCondition(),
		ProjectionExpression: projectionExpression,
	}

	output, err := ddb.Query(ctx, &query)
	if err != nil {
		return nil, fmt.Errorf("Could not query Instances table [%v]", err)
	}

	var records []InstanceNameType
	err = attributevalue.UnmarshalListOfMaps(output.Items, &records)
	return &records, err
}

func QueryInstancesUsingPort(ctx context.Context, ddb *dynamodb.Client, port uint16) (*[]InstanceNameType, error) {
	kexpr := expression.Key(*InstancePorts.Port.AttributeName).Equal(expression.Value(port))
	expr, err := expression.NewBuilder().WithKeyCondition(kexpr).Build()
	if err != nil {
		return nil, fmt.Errorf("Unable to create expression for query [%v]", err)
	}

	projectionExpression := InstancePorts.Instance.AttributeName
	input := dynamodb.QueryInput {
		TableName: InstancePorts.TableName,
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression: expr.KeyCondition(),
		ProjectionExpression: projectionExpression,
	}

	output, err := ddb.Query(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("Could not query InstancePorts table [%v]", err)
	}

	var records []InstanceNameType
	err = attributevalue.UnmarshalListOfMaps(output.Items, &records)
	if err != nil {
		return nil, err
	}

	return &records, nil
}

func ConsistentQueryShopByStream(ctx context.Context, ddb *dynamodb.Client, stream string) (*ShopType, error) {
	kexpr := expression.Key(*Shops.Stream.AttributeName).Equal(expression.Value(stream))
	expr, err := expression.NewBuilder().WithKeyCondition(kexpr).Build()
	if err != nil {
		return nil, fmt.Errorf("Unable to create expression for query [%v]", err)
	}

	consistentRead := true
	input := dynamodb.QueryInput {
		TableName: Shops.TableName,
		IndexName: &ShopsGsiStream,
		Select: types.SelectAllAttributes,
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression: expr.KeyCondition(),
		ConsistentRead: &consistentRead,
	}

	output, err := ddb.Query(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("Could not query Shops table with stream %v [%v]", stream, err)
	}

	if len(output.Items) == 0 {
		return nil, fmt.Errorf("No shop record found for stream %s", stream)
	} else if len(output.Items) > 1 {
		// don't panic, since we aren't adding to our data corruption problem here
		log.Println(fmt.Sprintf("ERROR: More than one shop for stream %s was detected"))
	}

	var records []ShopType
	err = attributevalue.UnmarshalListOfMaps(output.Items, &records)
	if err != nil {
		return nil, err
	}

	return &records[0], nil
}
