package lb

import (
	"fmt"
	"log"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"loadbalancer/go/tables"
)

const MAX_INSTANCES uint8 = 3

func Register(shopId string, stream string, port uint16) (string, string, int, error) {
	context, err := tables.Context()
	if err != nil {
		return "", "", 500, err
	}

	ctx := context.Ctx()
	cfg := context.Cfg()
	ddb := dynamodb.NewFromConfig(*cfg)
	shopIdExists, err := tables.TestShopIdPresence(ctx, ddb, shopId)
	if err != nil {
		return "", "", 500, err
	}

	if shopIdExists {
		return "", "",  400, fmt.Errorf(fmt.Sprintf("Shop %v in use", shopId))
	}

	streamExists, err := tables.TestStreamPresence(ctx, ddb, stream)
	if err != nil {
		return "", "",  500, err
	}

	if streamExists {
		return "", "", 400, fmt.Errorf(fmt.Sprintf("Stream %v in use", stream))
	}

	instanceNamesUsingPort, err := tables.QueryInstancesUsingPort(ctx, ddb, port)
	if err != nil {
		return "", "", 500, err
	}

	if len(*instanceNamesUsingPort) > int(MAX_INSTANCES) {
		// bug in code, data integrity compromised
		panic(fmt.Sprintf("More than %d instances [%v] using port %d", MAX_INSTANCES, *instanceNamesUsingPort, port))
	}

	if len(*instanceNamesUsingPort) == int(MAX_INSTANCES) {
		return "", "", 400, fmt.Errorf(fmt.Sprintf("Port %d in use", port))
	}

	instancesSetUsingPort := map[string]interface{}{}
	for _, i := range *instanceNamesUsingPort {
		// we only need a hashset to later test if the instance is
		// using the port or not.
		instancesSetUsingPort[i.Instance] = nil
	}

	var streams uint8 = 0
	for streams = 0; streams < MAX_INSTANCES; streams++ {
		instanceNameRecords, er := tables.QueryAllInstancesWithNumStreams(ctx, ddb, streams)
		err = er
		if err == nil && instanceNameRecords != nil {
			for _, record := range *instanceNameRecords {
				if _, present := instancesSetUsingPort[record.Instance]; !present {
					instanceRecord, er := tables.ConsistentGetInstance(ctx, ddb, record.Instance)
					err = er
					if err == nil && instanceRecord != nil && instanceRecord.Streams < MAX_INSTANCES {
						publicIp, privateIp, er := tables.GetIps(ctx, ddb, record.Instance)
						err = er
						if err == nil {
							er = tables.TransactAddStream(ctx, ddb, shopId, stream, port, instanceRecord)
							err = er
							if err == nil {
								return publicIp, privateIp, 200, nil
							}
						}
					}
				}
			}
	    }

		// Preferred less code nesting over meticulous error logging. The code can be changed to get more
		// precise error logging, if this code ever encounters issues needing deeper troubleshooting.
		if err != nil {
			log.Printf("INFO: Error encountered streams=%d shopId=%v stream=%v port=%d \n instance name records: %v \n instancesSetUsingPort %v --> %v", streams, shopId, stream, port, *instanceNameRecords, instancesSetUsingPort, err)
		}
	}

	if err != nil {
		return "", "", 500, err
	}

	return "", "", 503, fmt.Errorf(fmt.Sprintf("Unable to allocate for %v, %v and %d", shopId, stream, port))
}

func Unregister(stream string) (int, error) {
	context, err := tables.Context()
	if err != nil {
		return 500, err
	}

	ctx := context.Ctx()
	cfg := context.Cfg()
	ddb := dynamodb.NewFromConfig(*cfg)
	shopId, err := tables.QueryShopIdByStream(ctx, ddb, stream)
	if err != nil {
		return 500, err
	}

	if shopId == "" {
		return 400, fmt.Errorf("Stream %s does not exist", stream)
	}

	shop, err := tables.ConsistentGetShop(ctx, ddb, shopId)
	if err != nil {
		return 500, err
	}

	if shop.Stream != stream {
		return 500, fmt.Errorf("That stream may not have been consistently written yet")
	}

	instanceRecord, err := tables.ConsistentGetInstance(ctx, ddb, shop.Instance)
	if err != nil {
		return 500, err
	}

	err = tables.TransactDelete(ctx, ddb, shop, instanceRecord)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
