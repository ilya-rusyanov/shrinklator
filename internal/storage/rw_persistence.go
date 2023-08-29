package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type RWPersistence struct {
	reader io.Reader
	writer *bufio.Writer
	lastID int
}

type dto struct {
	ID    string `json:"uuid"`
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

func NewRWPersistence(reader io.Reader, writer io.Writer) *RWPersistence {
	w := bufio.NewWriter(writer)
	return &RWPersistence{
		reader: reader,
		writer: w,
		lastID: 1,
	}
}

func (p *RWPersistence) Append(short string, long string) error {
	payload := dto{Short: short, Long: long, ID: strconv.Itoa(p.lastID)}
	encoded, err := json.Marshal(&payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if _, err := p.writer.Write(encoded); err != nil {
		return fmt.Errorf("error writing encoded data: %w", err)
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("error adding newline: %w", err)
	}

	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("error flushing file: %w", err)
	}

	p.lastID++

	return nil
}

func (p *RWPersistence) ReadAll() (values map[string]string, err error) {
	values = make(map[string]string)

	reader := bufio.NewReader(p.reader)
	data, err := reader.ReadBytes('\n')
	for ; !errors.Is(err, io.EOF); data, err = reader.ReadBytes('\n') {
		if err != nil {
			err = fmt.Errorf("error reading line: %w", err)
			return
		}

		dto := dto{}
		err = json.Unmarshal(data, &dto)
		if err != nil {
			err = fmt.Errorf("error unmarshaling JSON: %w", err)
			return
		}

		values[dto.Short] = dto.Long
		p.lastID++
	}

	// clear EOF error
	err = nil

	return
}
