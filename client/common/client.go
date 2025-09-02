package common

import (
	"context"
	"errors"
	"time"
	"strings"
	"os"
	
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
	//conn   net.Conn
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
		socket := NewSocket()
        err := socket.Connect(c.config.ServerAddress)
        if err != nil {
	        log.Errorf("error when opening connection | result: fail | error: %v", err)
	        return
        }

		
		apuesta := Apuesta{
		Nombre: os.Getenv("NOMBRE"),
		Apellido:            os.Getenv("APELLIDO"),
		Documento:            os.Getenv("DOCUMENTO"),
		Nacimiento:            os.Getenv("NACIMIENTO"),
		Numero:            os.Getenv("NUMERO"),
		}
		
		//msgStr := c.config.ID + "/" + apuesta.toString()
		
		if sentMsg, err := socket.Send(c.config.ID, apuesta); err != nil {
	        log.Errorf("action: send_message | result: fail | error: %v", err)
	        return
        } else {
            log.Infof("action: send_message | result: success | msg: %s", sentMsg)
        }
		
		response, err := socket.ReadResponse(ctx)
		socket.Close()
		log.Infof("action: socket_closed | result: success | client_id: %v", c.config.ID)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Infof("action: receive_message | result: cancelled | client_id: %v", c.config.ID)
				return
			}
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		if strings.TrimSpace(response) == "OK" {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", apuesta.Documento, apuesta.Numero)
		} else {
			log.Infof("action: receive_message | result: unexpected_response | client_id: %v | response: %v", c.config.ID, response)
		}
		
		

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
