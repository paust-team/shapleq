package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elon0823/paustq/client/consumer"
	"github.com/elon0823/paustq/message"
	"github.com/elon0823/paustq/proto"
	"testing"
	"time"
)

func mockConsumerHandler(serverReceiveChannel chan []byte, serverSendChannel chan []byte, testData RecordMap) {

	for received := range serverReceiveChannel {
		fetchReqMsg := &paustq_proto.FetchRequest{}

		if err := message.UnPackTo(received, fetchReqMsg); err != nil {
			continue
		}

		if testData[fetchReqMsg.TopicName] != nil {
			topicData := testData[fetchReqMsg.TopicName]

			for _, record := range topicData {
				fetchResMsg, err := message.NewFetchResponseMsgData(record, 0)
				if err != nil {
					fmt.Println("Failed to create FetchResponse message")
					continue
				}
				serverSendChannel <- fetchResMsg
			}

			fetchEnd, err := message.NewFetchResponseMsgData(nil, 1)
			if err != nil {
				fmt.Println("Failed to create FetchResponse message")
				continue
			}

			serverSendChannel <- fetchEnd
		}
	}
}

func TestConsumer_Subscribe(t *testing.T) {

	ip := "127.0.0.1"
	port := ":8001"
	timeout := 5

	testRecordMap := map[string][][]byte{
		"topic1": {
			{'g', 'o', 'o', 'g', 'l', 'e'},
			{'p', 'a', 'u', 's', 't', 'q'},
			{'1', '2', '3', '4', '5', '6'}},
	}
	topic := "topic1"
	receivedRecordMap := make(map[string][][]byte)

	host := fmt.Sprintf("%s%s", ip, port)
	ctx := context.Background()

	// Start Server
	server, err := StartTestServer(port, mockConsumerHandler, testRecordMap)
	if err != nil {
		t.Error(err)
		return
	}

	defer server.Stop()

	// Start Client
	client := consumer.NewConsumer(ctx, host, time.Duration(timeout))
	if client.Connect() != nil {
		t.Error("Error on connect")
	}

	for response := range client.Subscribe(topic) {
		if response.Error != nil {
			t.Error(response.Error)
		} else {
			receivedRecordMap[topic] = append(receivedRecordMap[topic], response.Data)
		}
	}

	expectedResults := testRecordMap[topic]
	receivedResults := receivedRecordMap[topic]

	if len(expectedResults) != len(receivedResults) {
		t.Error("Length Mismatch - Expected records: ", len(expectedResults), ", Received records: ", len(receivedResults))
	}
	for i, record := range expectedResults {
		if bytes.Compare(receivedResults[i], record) != 0 {
			t.Error("Record is not same")
		}
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}
}
