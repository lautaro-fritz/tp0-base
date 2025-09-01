package common

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Batch struct {
	Reader        *csv.Reader
	MaxBatchSize  int
	MaxMessageLen int
	ClientID      string
}

func (b *Batch) NextBatch() ([]Apuesta, string, error) {
	var sb strings.Builder
	sb.WriteString(b.ClientID)

	batch := make([]Apuesta, 0, b.MaxBatchSize)

	for len(batch) < b.MaxBatchSize {
		record, err := b.Reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "EOF") {
				// No more records
				return batch, sb.String(), nil
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
		if sb.Len()+len(apuestaStr) > b.MaxMessageLen {
			break // avoid exceeding byte limit
		}

		sb.WriteString(apuestaStr)
		batch = append(batch, apuesta)
	}

	return batch, sb.String(), nil
}
