//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"

	"google.golang.org/grpc"
)

type RegisterEventListenerRequest struct {
	ScName      string `json:"scName"`
	EventFilter string `json:"eventFilter"`
}

type EventWrapper struct {
	ScEvent *ScEvent
	Error   error
}

type ScEvent struct {
	TxId      string `json:"txId"`
	ScName    string `json:"scName"`
	EventName string `json:"eventName"`
	Payload   string `json:"payload"`
}

// RegisterEventListener ...
func (client *Client) RegisterEventListener(scName string, eventFilter string) (*ListenerController, <-chan *EventWrapper, error) {
	// Check regular expression is valid
	_, err := regexp.Compile(eventFilter)
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("CLIENT: %w", err)
	}

	// Prepare the request
	payload := RegisterEventListenerRequest{
		ScName:      scName,
		EventFilter: eventFilter,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("CLIENT: %w", err)
	}

	// Establish stream connection first
	stream, err := client.grpcClient.RegisterEventListener(context.Background())
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("CLIENT: %w", err)
	}

	// Send the event listener parameters
	err = stream.Send(&pb.Request{Payload: payloadBytes})
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("CLIENT: %w", err)
	}

	// Receive a success message from server
	resp, err := stream.Recv()
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("CLIENT: %w", err)
	}
	if resp.Error != nil {
		client.Close()
		return nil, nil, fmt.Errorf("%v", resp.Error)
	}
	successMsg := string(resp.Payload)
	if successMsg != "Successfully registered event listener." {
		client.Close()
		return nil, nil, fmt.Errorf("CLIENT: Did not receive expected success message")
	}

	eventChannel := make(chan *EventWrapper)
	streamAlive := true
	go client.listenEvents(stream, eventChannel, &streamAlive)

	return &ListenerController{eventChannel: eventChannel, streamAlive: &streamAlive, conn: client.conn}, eventChannel, nil
}

func (client *Client) listenEvents(stream pb.RequestHandler_RegisterEventListenerClient, eventChannel chan *EventWrapper, streamAlive *bool) {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			eventChannel <- &EventWrapper{
				ScEvent: nil,
				Error:   fmt.Errorf("CLIENT: Read io.EOF. Stream closed by server."),
			}
			closeEventListener(eventChannel, streamAlive)
			return
		}
		if err != nil {
			eventChannel <- &EventWrapper{
				ScEvent: nil,
				Error:   fmt.Errorf("CLIENT: %w", err),
			}
			closeEventListener(eventChannel, streamAlive)
			return
		}

		// Check whether the pb.Response contains error
		if resp.Error != nil {
			eventChannel <- &EventWrapper{
				ScEvent: nil,
				Error:   fmt.Errorf("%v", resp.Error),
			}
			closeEventListener(eventChannel, streamAlive)
			return
		}

		// Parse out ScEvent
		var scEvent ScEvent
		err = json.Unmarshal(resp.Payload, &scEvent)
		if err != nil {
			eventChannel <- &EventWrapper{
				ScEvent: nil,
				Error:   fmt.Errorf("CLIENT: %v", resp.Error),
			}
			closeEventListener(eventChannel, streamAlive)
			return
		}

		// Send ScEvent to channel
		eventChannel <- &EventWrapper{
			ScEvent: &scEvent,
			Error:   nil,
		}
	}

	closeEventListener(eventChannel, streamAlive)
}

func closeEventListener(eventChannel chan *EventWrapper, streamAlive *bool) {
	close(eventChannel)
	*streamAlive = false
}

type ListenerController struct {
	eventChannel chan *EventWrapper
	conn         *grpc.ClientConn
	streamAlive  *bool
}

func (cc ListenerController) Close() {
	cc.conn.Close()
	for *cc.streamAlive {
		select {
		case _ = <-cc.eventChannel:
		default:
			time.Sleep(time.Millisecond)
		}
	}

	fmt.Println("Event Listener Removed.")
}
