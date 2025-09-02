package common

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Protocol struct {
	Reader        *csv.Reader
	MaxBatchSize  int
	MaxMessageLen int
	ClientID      string
}

func (p *Protocol) formatMessage(messageType MessageType) (string, error) {
	var sb strings.Builder
	sb.WriteString(p.ClientID)
	
	switch messageType {
	case BetsMessage:
	    sb.WriteString("#" + 'B')
	    batch := make([]Apuesta, 0, p.MaxBatchSize)

	    for len(batch) < p.MaxBatchSize {
		    record, err := p.Reader.Read()
		    if err != nil {
			    if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "EOF") {
				    // No more records
				    return sb.String(), nil
			    }
			    return nil, "", fmt.Errorf("read_csv_row: %w", err)
		    }

		    apuesta := Apuesta{
			    Nombre:     record[0],
			    Apellido:   record[1],
			    Documento:  record[2],
			    Nacimiento: record[3],
			    Numero:     record[4],
		    }

		    apuestaStr := "#" + apuesta.toString()
		    if sb.Len() + len(apuestaStr) > p.MaxMessageLen {
			    break // avoid exceeding byte limit
		    }

		    sb.WriteString(apuestaStr)
		    batch = append(batch, apuesta)
	    }
	case DoneMessage:
	    sb.WriteString("#" + 'D')
	
	case WinnersMessage:
	    sb.WriteString("#" + 'W')
	    
	default:
	} 


	return sb.String(), nil
}
