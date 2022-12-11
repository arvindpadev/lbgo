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

	fmt.Println(fmt.Sprintf("SUCCESS: TestRegister IP addresses --> %v %v %d", publicIp, privateIp, status))
}

func TestUnregister(t *testing.T) {
	status, err := Unregister("stream0")
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unregister Error: %d [%v]", status, err))
	}
}

func TestMain(m *testing.M) {
	test_setup.Setup()
	m.Run()
}