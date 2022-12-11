package test_setup

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"loadbalancer/go/tables"
)

const (
	awsAccessKeyId = "junk"
	awsSecretKeyId = "bunk"
	awsRegion = "localhost"
	awsEndpoint = "http://localhost:22000"
)

const (
	instance0 = "instance0"
	instance1 = "instance1"
	instance2 = "instance2"
)

type ddbLocalEndpointResolverWithOptions struct {}
func (resolver ddbLocalEndpointResolverWithOptions) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint {
		URL: awsEndpoint,
	}, nil
}

type ddbLocalCredentialProvider struct {}
func (credprovider ddbLocalCredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	credentials := aws.Credentials {
		AccessKeyID: awsAccessKeyId,
		SecretAccessKey: awsSecretKeyId,
	}

	return credentials, nil
}

func Setup() {
	ctx, cfg := ddbLocalConfig()
	ddbLocal := dynamodb.NewFromConfig(*cfg)
	version := uuid.New().String()
	createShopsTable(ctx, ddbLocal)
	createInstancesTable(ctx, ddbLocal, version)
	createInstancePortTable(ctx, ddbLocal)
	createStreams(ctx, ddbLocal)
	createInstanceIpTable(ctx, ddbLocal)
	_, err := ddbLocal.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		panic(fmt.Sprintf("Error listing tables %v", err))
	}

	tables.TestMockAwsCfg = cfg
}

func ddbLocalConfig() (context.Context, *aws.Config) {
	ctx := context.TODO()
	cfg := aws.NewConfig()
	cfg.EndpointResolverWithOptions = ddbLocalEndpointResolverWithOptions {}
	cfg.Credentials = ddbLocalCredentialProvider {}
	cfg.Region = awsRegion
	return ctx, cfg
}

func createShopsTable(ctx context.Context, ddb *dynamodb.Client) {
	input := dynamodb.CreateTableInput {
		TableName: tables.Shops.TableName,
		AttributeDefinitions: []types.AttributeDefinition {
			tables.Shops.ShopId,
			tables.Shops.Stream,
		},
		KeySchema: tables.Shops.KeySchema,
		GlobalSecondaryIndexes: tables.Shops.Gsi,
		ProvisionedThroughput: tables.Shops.ProvisionedThroughput,
	}
	createTable(ctx, ddb, &input)
}

func createInstancePortTable(ctx context.Context, ddb *dynamodb.Client) {
	input := dynamodb.CreateTableInput {
		TableName: tables.InstancePorts.TableName,
		AttributeDefinitions: []types.AttributeDefinition {
			tables.InstancePorts.Instance,
			tables.InstancePorts.Port,
		},
		KeySchema: tables.InstancePorts.KeySchema,
		ProvisionedThroughput: tables.InstancePorts.ProvisionedThroughput,
	}
	createTable(ctx, ddb, &input)
}

func createStreams(ctx context.Context, ddb *dynamodb.Client) {
	input := dynamodb.CreateTableInput {
		TableName: tables.StreamNames.TableName,
		AttributeDefinitions: []types.AttributeDefinition {
			tables.StreamNames.Stream,
		},
		KeySchema: tables.StreamNames.KeySchema,
		ProvisionedThroughput: tables.StreamNames.ProvisionedThroughput,
	}
	createTable(ctx, ddb, &input)
}

type InstanceIpCreateType struct {
	Instance string
	PublicIp string
	PrivateIp string
}

func createInstanceIpTable(ctx context.Context, ddb *dynamodb.Client) {
	createInput := dynamodb.CreateTableInput {
		TableName: tables.InstanceIp.TableName,
		AttributeDefinitions: []types.AttributeDefinition {
			tables.InstanceIp.Instance,
		},
		KeySchema: tables.InstanceIp.KeySchema,
		ProvisionedThroughput: tables.InstanceIp.ProvisionedThroughput,
	}
	createTable(ctx, ddb, &createInput)
	item0 := InstanceIpCreateType {
		Instance: instance0,
		PublicIp: "189.189.189.191",
		PrivateIp: "10.1.1.1",
	}
	putItem(ctx, ddb, tables.InstanceIp.TableName, item0)
	item1 := InstanceIpCreateType {
		Instance: instance1,
		PublicIp: "189.189.189.189",
		PrivateIp: "10.1.1.3",
	}
	putItem(ctx, ddb, tables.InstanceIp.TableName, item1)
	item2 := InstanceIpCreateType {
		Instance: instance2,
		PublicIp: "189.189.189.190",
		PrivateIp: "10.1.1.2",
	}
	putItem(ctx, ddb, tables.InstanceIp.TableName, item2)
}

func createInstancesTable(ctx context.Context, ddb *dynamodb.Client, version string) {
	createInput := dynamodb.CreateTableInput {
		TableName: tables.Instances.TableName,
		AttributeDefinitions: []types.AttributeDefinition {
			tables.Instances.Streams,
			tables.Instances.Instance,
			tables.Instances.Version,
		},
		KeySchema: tables.Instances.KeySchema,
		GlobalSecondaryIndexes: tables.Instances.Gsi,
		ProvisionedThroughput: tables.Instances.ProvisionedThroughput,
	}
	createTable(ctx, ddb, &createInput)
	item0 := tables.InstanceType {
		Instance: instance0,
		Streams: 0,
		Version: version,
	}
	putItem(ctx, ddb, tables.Instances.TableName, item0)
	item1 := tables.InstanceType {
		Instance: instance1,
		Streams: 0,
		Version: version,
	}
	putItem(ctx, ddb, tables.Instances.TableName, item1)
	item2 := tables.InstanceType {
		Instance: instance2,
		Streams: 0,
		Version: version,
	}
	putItem(ctx, ddb, tables.Instances.TableName, item2)
}

func createTable(ctx context.Context, ddb *dynamodb.Client, input *dynamodb.CreateTableInput) {
	_, err := ddb.CreateTable(ctx, input)
	if err != nil {
		panic(fmt.Sprintf("Error creating %v Table [%v]", *input.TableName, err))
	}
}

func putItem(ctx context.Context, ddb *dynamodb.Client, tableName *string, itemPut interface{}) {
	item, err := attributevalue.MarshalMap(itemPut)
	if err != nil {
		panic(fmt.Sprintf("Unable to marshal items to put in %v because of [%v]", *tableName, err))
		return
	}

	putItem := dynamodb.PutItemInput {
		TableName: tableName,
		Item: item,
	}

	_, err = ddb.PutItem(ctx, &putItem)
	if err != nil {
		panic(fmt.Sprintf("Unable to put item %v in %v because of [%v]", putItem, *tableName, err))
	}
}
