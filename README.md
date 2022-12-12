We need to handle uniqueness checks for which I was attempting to use
GSIs in the original design.
But the AWS docs had this info:
https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/transaction-apis.html#transaction-apis-txwriteitems
"Once a transaction completes, the changes made within that transaction are
  propagated to global secondary indexes (GSIs), streams, and backups. ..."

I am therefore adding 2 additional tables in place of those GSIs - a streamNames table, and an instancePorts table.

To recap we have 4 tables now - shops, instances, streamNames and instancePorts
The attributes for these tables are - 
streamNames - Stream (H)
instancePorts - Port (H), Instance (R)
shops - ShopId (H), Stream, Port, Instance, Version
instances - Streams, Instance (H), Version

The GSIs from the previous doc repeated here are still around and they are used
in the deletion logic to make sure that deletions from the streamName and instancePort tables are done correctly, after we have made sure of the stream
still being associated with the instance and port. The GSIs are
instances - Streams (H), Instance (R) to find instances having 0..MAX_INSTANCES-1 streams
shops - Stream (H) to look up the stream to be deleted and get the associated port and instance to delete in instancePorts and change instances

The transacted writes to all 4 of these tables will ensure all of the uniqueness
constraints. Changing the code based on evolving requirements can be tricky though.
Both Register and Unregister will require transact writes on all 4 tables to perform the operations without breaking the data integrity needed

Version attributes are used only in the instances and shops table. This is to
allow us to use it for optimistic locking like condition checking within the transactions. The Version attribute in the instances table is use during registration. The Version attribute in the shops table is used
during deallocation.

There is a 5th table that is used to get a public and private ip address for an
instance. This is more of an assumption to get ip addresses to return, and it isn't
integral to the design. The lookup table to get ip addresses to return
instanceIp - Instance (H), PrivateIp, PublicIp

How to run the tests:
1. Change directory to where DynamoDB local is installed. Run DynamoDB local
```
java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -inMemory -port 22000
```
  (NOTE: This will need to be stopped with CTRL+C and started again, each time the tests are run)

2. In another terminal navigate to the root directory of the repo and run
```
go test
```
  (NOTE: The tests are not well written meaning, they are not independent of
  each other. They have to be run in order)

NOTE: Only the Register function has been implemented.
