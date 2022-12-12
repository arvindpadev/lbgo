package lb

import (
	"context"
 	"fmt"
 	"testing"
 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

    "loadbalancer/go/test_setup"
)

var testCtx context.Context
var testClient *dynamodb.Client

func TestRegister(t *testing.T) {
	publicIp, privateIp, status, err := Register("shop0", "stream0", 11000)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Register Error: [%v]", err))
	}

	fmt.Println(fmt.Sprintf("SUCCESS: TestRegister IP addresses (%d) --> %v %v", status, publicIp, privateIp))
}

func TestUnregister(t *testing.T) {
	status, err := Unregister("stream0")
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unregister Error: %d [%v]", status, err))
	}
}

func TestRegisterRepeatingShop(t* testing.T) {
	_, _, status, err := Register("shop0", "stream1", 11001)
	if err == nil {
		t.Fatalf("No error was received")
	}

	fmt.Println(fmt.Sprintf("SUCCESS: TestRegisterRepeatingShop Expected error received (%d) --> %v", status, err))
}

func TestRegisterRepeatingStream(t* testing.T) {
	_, _, status, err := Register("shop1", "stream0", 11002)
	if err == nil {
		t.Fatalf("No error was received")
	}

	fmt.Println(fmt.Sprintf("SUCCESS: TestRegisterRepeatingStream Expected error received (%d) --> %v", status, err))
}

func TestRegisterRepeatingPortAfterPortExhaustion(t* testing.T) {
	_, _, status, err := Register("shop1", "stream1", 11000)
	if err != nil {
		t.Fatalf("(%d) Port 11000 should have successfully been allocated to shop1 and stream1 ---> %v", status, err)
	}

	_, _, status, err = Register("shop2", "stream2", 11000)
	if err != nil {
		t.Fatalf("(%d) Port 11000 should have successfully been allocated to shop2 and stream2 ---> %v", status, err)
	}

	_, _, status, err = Register("shop3", "stream3", 11000)
	if err == nil {
		t.Fatalf("Port 11000 should have NOT been allocated to shop3 and stream3 because port 11000 has already been allocated to the 3 different instances.")
	}

	fmt.Println(fmt.Sprintf("SUCCESS: TestRegisterRepeatingPortAfterPortExhaustion Expected error received (%d) --> %v", status, err))
}

func TestRegisterInstanceExhaustion(t* testing.T) {
	var shopSeed uint16 = 10
	var streamSeed uint16 = 10
	var portSeed uint16 = 12000
	var i uint16

	// 3 instances with one stream each already allocated. 6
	// more allocations to go
	for i = 0; i < 6; i++ {
		shopId := fmt.Sprintf("shop%d", shopSeed + i)
		stream := fmt.Sprintf("stream%d", streamSeed + i)
		port := portSeed + i
		_, _, status, err := Register(shopId, stream, port)
		if err != nil {
			t.Fatalf("(%d) Port %d should have successfully been allocated to %s and %s ---> %v", status, port, shopId, stream, err)
		}
	}

	// Now attempt to register with a completely different shopId, stream and port
	// than anything that was previously registered
	_, _, status, err := Register("MYSHOP", "MYSTREAM", 14000)
	if err == nil {
		t.Fatalf("(%d) An error should have been received when registering MYSHOP, MYSTREAM, 14000", status)
	}

	fmt.Println(fmt.Sprintf("SUCCESS: TestRegisterInstanceExhaustion Expected error received (%d) --> %v", status, err))
}

func TestMain(m *testing.M) {
	test_setup.Setup()
	m.Run()
}