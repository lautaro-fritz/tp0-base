package common

import (
	"bufio"
	"context"
	"encoding/csv"
	"os"
	"strings"
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
	BatchMaxAmount      int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	file, err := os.Open("/agency.csv")
	if err != nil {
		log.Fatalf("action: open_csv | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}
	defer file.Close()
	
	reader := csv.NewReader(bufio.NewReader(file))

	const maxBytes = 8192
	batchNumber := 1
	
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {
		case <-ctx.Done():
			// If context is cancelled, stop the loop
			log.Infof("action: loop_cancelled | result: success | client_id: %v", c.config.ID)
			return
		default:
		}
		
		batch := &Batch{
		    Reader:        reader,
		    MaxBatchSize:  c.config.BatchMaxAmount,
		    MaxMessageLen: maxBytes,
		    ClientID:      c.config.ID,
	    }
	    
	    bets, msg, err := batch.NextBatch()
		if err != nil {
			log.Warningf("action: build_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			break
		}

		if len(bets) == 0 {
			// No more apuestas to send
			log.Infof("esta entrando aca")
			break
		}

		socket := NewSocket()
		err = socket.Connect(c.config.ServerAddress)
		if err != nil {
			log.Errorf("action: open_connection | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		err = socket.Send(msg)
		if err != nil {
			log.Errorf("action: send_batch | result: fail | batch_number: %d | client_id: %v | error: %v", batchNumber, c.config.ID, err)
			socket.Close()
			return
		}

		log.Infof("action: send_batch | result: success | batch_number: %d | client_id: %v | cantidad: %d | size_bytes: %d",
			batchNumber, c.config.ID, len(bets), len(msg))

		response, err := socket.ReadResponse(ctx)
		socket.Close()

		if err != nil {
			log.Errorf("action: read_response | result: fail | batch_number: %d | client_id: %v | error: %v", batchNumber, c.config.ID, err)
			return
		}

		trimmedResp := strings.TrimSpace(response)
		if trimmedResp == "OK" {
			log.Infof("action: apuesta_enviada | result: success | cantidad: %d", len(bets))
		} else {
			log.Infof("action: apuesta_enviada | result: fail | cantidad: %d | response: %s", len(bets), trimmedResp)
		}

		batchNumber++

		select {
		case <-ctx.Done():
			log.Infof("action: loop_cancelled_during_sleep | result: success | client_id: %v", c.config.ID)
			return
		case <-time.After(c.config.LoopPeriod):
		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
