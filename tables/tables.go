package tables

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

var shopsStr = "shops"
var instancesStr = "instances"
var instanceIp = "instanceIp"
var streamNames = "streamNames"
var instancePorts = "instancePorts"
var shopId = "ShopId"
var streamStr = "Stream"
var instanceStr = "Instance"
var portStr = "Port"
var streamsStr = "Streams"
var publicIp = "PublicIp"
var privateIp = "PrivateIp"
var versionStr = "Version"
var ShopsGsiStream = "ShopsGsiStream"
var InstancesGsiStreamsInstance = "InstancesGsiStreamsInstance"
var projectionAll = types.Projection { ProjectionType: types.ProjectionTypeAll }
var readCapacity int64 = 5
var writeCapacity int64 = 5
var provisionedThroughput = types.ProvisionedThroughput { ReadCapacityUnits: &readCapacity, WriteCapacityUnits: &writeCapacity }

type shopsTableType struct {
	TableName *string
	ProvisionedThroughput *types.ProvisionedThroughput
	ShopId types.AttributeDefinition
	Stream types.AttributeDefinition
	Instance types.AttributeDefinition
	Port types.AttributeDefinition
	Version types.AttributeDefinition
	KeySchema []types.KeySchemaElement
	Gsi []types.GlobalSecondaryIndex
}

var Shops = shopsTableType {
	TableName: &shopsStr,
	ProvisionedThroughput: &provisionedThroughput,
	ShopId: types.AttributeDefinition { AttributeName: &shopId, AttributeType: types.ScalarAttributeTypeS },
	Stream: types.AttributeDefinition { AttributeName: &streamStr, AttributeType: types.ScalarAttributeTypeS },
	Instance: types.AttributeDefinition { AttributeName: &instanceStr, AttributeType: types.ScalarAttributeTypeS },
	Port: types.AttributeDefinition { AttributeName: &portStr, AttributeType: types.ScalarAttributeTypeN },
	Version: types.AttributeDefinition { AttributeName: &versionStr, AttributeType: types.ScalarAttributeTypeS },
	KeySchema: []types.KeySchemaElement {
		types.KeySchemaElement { AttributeName: &shopId, KeyType: types.KeyTypeHash },
	},
	Gsi: []types.GlobalSecondaryIndex {
		types.GlobalSecondaryIndex {
		    IndexName: &ShopsGsiStream,
			Projection: &projectionAll,
			ProvisionedThroughput: &provisionedThroughput,
		    KeySchema: []types.KeySchemaElement {
			    types.KeySchemaElement { AttributeName: &streamStr, KeyType: types.KeyTypeHash },
			},
		},
	},
}

type instanceTableType struct {
	TableName *string
	ProvisionedThroughput *types.ProvisionedThroughput
	Instance types.AttributeDefinition
	Streams types.AttributeDefinition
	KeySchema []types.KeySchemaElement
	Gsi []types.GlobalSecondaryIndex

	// Used as for application managed optimistic locking
	Version types.AttributeDefinition
}

var Instances = instanceTableType {
	TableName: &instancesStr,
	ProvisionedThroughput: &provisionedThroughput,
	Instance: types.AttributeDefinition { AttributeName: &instanceStr, AttributeType: types.ScalarAttributeTypeS },
	Streams: types.AttributeDefinition { AttributeName: &streamsStr, AttributeType: types.ScalarAttributeTypeN },
	Version: types.AttributeDefinition { AttributeName: &versionStr, AttributeType: types.ScalarAttributeTypeS },
	KeySchema: []types.KeySchemaElement {
	    types.KeySchemaElement { AttributeName: &instanceStr, KeyType: types.KeyTypeHash },
	},
	Gsi: []types.GlobalSecondaryIndex {
	    types.GlobalSecondaryIndex {
		    IndexName: &InstancesGsiStreamsInstance,
			Projection: &projectionAll,
			ProvisionedThroughput: &provisionedThroughput,
		    KeySchema: []types.KeySchemaElement {
				types.KeySchemaElement { AttributeName: &streamsStr, KeyType: types.KeyTypeHash },
				types.KeySchemaElement { AttributeName: &instanceStr, KeyType: types.KeyTypeRange },
			},
		},
	},
}

type instanceIpTableType struct {
	TableName *string
	ProvisionedThroughput *types.ProvisionedThroughput
	Instance types.AttributeDefinition
	PublicIp types.AttributeDefinition
	PrivateIp types.AttributeDefinition
	KeySchema []types.KeySchemaElement
}

var InstanceIp = instanceIpTableType {
	TableName: &instanceIp,
	ProvisionedThroughput: &provisionedThroughput,
	Instance: types.AttributeDefinition { AttributeName: &instanceStr, AttributeType: types.ScalarAttributeTypeS },
	PublicIp: types.AttributeDefinition { AttributeName: &publicIp, AttributeType: types.ScalarAttributeTypeS },
	PrivateIp: types.AttributeDefinition { AttributeName: &privateIp, AttributeType: types.ScalarAttributeTypeS },
	KeySchema: []types.KeySchemaElement {
	    types.KeySchemaElement { AttributeName: &instanceStr, KeyType: types.KeyTypeHash },
	},
}

type instancePortsType struct {
	TableName *string
	ProvisionedThroughput *types.ProvisionedThroughput
	Instance types.AttributeDefinition
	Port types.AttributeDefinition
	KeySchema []types.KeySchemaElement
}

var InstancePorts = instancePortsType {
	TableName: &instancePorts,
	ProvisionedThroughput: &provisionedThroughput,
	Instance: types.AttributeDefinition { AttributeName: &instanceStr, AttributeType: types.ScalarAttributeTypeS },
	Port: types.AttributeDefinition { AttributeName: &portStr, AttributeType: types.ScalarAttributeTypeN },
	KeySchema: []types.KeySchemaElement {
		types.KeySchemaElement { AttributeName: &portStr, KeyType: types.KeyTypeHash },
		types.KeySchemaElement { AttributeName: &instanceStr, KeyType: types.KeyTypeRange },
	},
}

type streamNamesType struct {
	TableName *string
	ProvisionedThroughput *types.ProvisionedThroughput
	Stream types.AttributeDefinition
	KeySchema []types.KeySchemaElement
}

var StreamNames = streamNamesType {
	TableName: &streamNames,
	ProvisionedThroughput: &provisionedThroughput,
	Stream: types.AttributeDefinition { AttributeName: &streamStr, AttributeType: types.ScalarAttributeTypeS },
	KeySchema: []types.KeySchemaElement {
	    types.KeySchemaElement { AttributeName: &streamStr, KeyType: types.KeyTypeHash },
	},
}

type ShopType struct {
	ShopId string
	Stream string
	Instance string
	Port uint16
	Version string
}

type InstanceType struct {
	Instance string
	Streams uint8
	Version string
}

type InstancePortType struct {
	Instance string
	Port uint16
}

type InstanceNameType struct {
	Instance string
}

type InstanceIpType struct {
	PublicIp string
	PrivateIp string
}

type StreamType struct {
	Stream string
}
