package etcdkv

import "testing"

func TestNewClientV3(t *testing.T) {
	if err := NewClientV3(ClientEndpoints("127.0.0.1:2379")); err != nil {
		t.Error(err)
		t.FailNow()
	}

	client := ClientV3(DefaultClientName)
	if client == nil {
		t.Error("client not found")
		t.FailNow()
	}

	t.Log("ClientV3:", client)
}
