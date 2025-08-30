package common

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

func readStringWithContext(ctx context.Context, conn net.Conn) (string, error) {
	readCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Runs concurrently
	go func() {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		readCh <- msg
	}()

	select {
	case <-ctx.Done():
		// If context is cancelled, stops reading the buffer
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case msg := <-readCh:
		return msg, nil
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {
		case <-ctx.Done():
			// If context is cancelled, stop the loop
			log.Infof("action: loop_cancelled | result: success | client_id: %v", c.config.ID)
			return
		default:
		}
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		// Handles context cancel, error and successful reads.
		msg, err := readStringWithContext(ctx, c.conn)
		c.conn.Close()
		log.Infof("action: socket_closed | result: success | client_id: %v", c.config.ID)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Infof("action: receive_message | result: cancelled | client_id: %v", c.config.ID)
				return
			}
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		select {
		case <-ctx.Done():
			// If context is cancelled, stop the sleep
			log.Infof("action: loop_cancelled_during_sleep | result: success | client_id: %v", c.config.ID)
			return
		case <-time.After(c.config.LoopPeriod):
		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
