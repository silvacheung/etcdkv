package etcdkv

import "testing"

func TestNewClientV3(t *testing.T) {
	client, err := NewClientV3(ClientEndpoints("127.0.0.1:2379"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("NewClientV3:", client)

	client, ok := ClientV3()
	if !ok {
		if err != nil {
			t.Error("client not found")
			t.FailNow()
		}
	}

	t.Log("ClientV3:", client)
}
