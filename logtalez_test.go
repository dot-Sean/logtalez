package logtalez

import (
	"testing"

	"github.com/zeromq/goczmq"
)

func TestMakeTopicList(t *testing.T) {
	hosts := "host1,host2"
	programs := "program1,program2"
	topics := MakeTopicList(hosts, programs)

	expected := []string{
		"host1.program1",
		"host1.program2",
		"host2.program1",
		"host2.program2",
	}

	for i, expect := range expected {
		if topics[i] != expect {
			t.Errorf("expected topic '%s', got '%s'", expect[i], topics[i])
		}
	}
}

func TestMakeEndpointList(t *testing.T) {
	conns := "tcp://incproc1,tcp://inproc2"
	endpoints := MakeEndpointList(conns)

	expected := []string{
		"tcp://incproc1",
		"tcp://inproc2",
	}

	for i, expect := range expected {
		if endpoints[i] != expect {
			t.Errorf("expected endpoint '%s', got '%s'", expect[i], endpoints[i])
		}
	}
}

func TestNew(t *testing.T) {
	endpoints := []string{"inproc://test1"}
	topics := []string{"topic1", "topic2"}
	servercert := "./example_certs/example_curve_server_cert"
	clientcert := "./example_certs/example_curve_client_cert"

	auth := goczmq.NewAuth()
	defer auth.Destroy()

	clientCert, err := goczmq.NewCertFromFile(clientcert)
	if err != nil {
		t.Fatal(err)
	}
	defer clientCert.Destroy()

	server := goczmq.NewSock(goczmq.PUB)

	defer server.Destroy()
	server.SetZapDomain("global")

	serverCert, err := goczmq.NewCertFromFile(servercert)
	defer serverCert.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	serverCert.Apply(server)
	server.SetCurveServer(1)

	err = auth.Curve("./example_certs/")
	if err != nil {
		t.Fatal(err)
	}

	server.Bind(endpoints[0])

	lt, err := New(endpoints, topics, servercert, clientcert)
	if err != nil {
		t.Error("NewLogTalez failed: %s", err)
	}

	server.SendFrame([]byte("topic1:hello world"), 0)

	msg := <-lt.TailChan
	if string(msg[0]) != "topic1:hello world" {
		t.Errorf("expected 'topic1:hello world', got '%s'", msg[0])
	}
}