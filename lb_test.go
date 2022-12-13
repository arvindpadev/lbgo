package lb

import (
	"context"
 	"fmt"
	"strings"
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
	_, _, status, err := Register("shopU", "streamU", 7000)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Register of shopU, streamU, 7000 should have succeeded: [%v]", err))
	}

	status, err = Unregister("streamU")
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unregister Error: %d [%v]", status, err))
	}

	fmt.Println(fmt.Sprintf("SUCCESS: TestUnregister (%d)", status))
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

func TestParallelUnregister(t* testing.T) {
	var streamSeed uint16 = 10
	channel := make(chan string, 6)
	var i uint16
	errors := []string{}

	fi := func(st uint16, fin uint16, ch chan<-string) {
		for i = st; i < fin; i++ {
			s := fmt.Sprintf("stream%d", streamSeed + i)
			go func(stream string) {
				status, err := Unregister(stream)
				if err != nil {
					ch <- fmt.Sprintf("Should have been able to unregister %s: [%d] [%v]", stream, status, err)
				} else {
					ch <- ""
				}
			}(s)
		}
	}

	fe := func(st uint16, fin uint16, ch <-chan string) {
		for i = st; i < fin; i++ {
			er := <-ch
			if len(er) > 0 {
				errors = append(errors, er)
			}
		}
    }

	// At more than 3 parallel unregisters, the transaction
	// fails with a conditional check exception presumably
	// because the version check failed.
	fi(0, 3, channel)
	fe(0, 3, channel)
	fi(3, 6, channel)
	fe(3, 6, channel)
	if len(errors) > 0 {
		t.Fatalf(strings.Join(errors, "\n"))
	}

	fmt.Println("SUCCESS: TestParallelUnregister")
}

func TestMain(m *testing.M) {
	test_setup.Setup()
	m.Run()
}