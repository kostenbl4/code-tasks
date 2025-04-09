package httpsender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kostenbl4/code-tasks/code-processor/internal/domain"
)

// HttpSender - вариант отправки результата в http
type HttpSender struct {
	client http.Client
}

func NewHttpSender(client http.Client) *HttpSender {
	return &HttpSender{client: client}
}

func (hs HttpSender) SendResult(ctx context.Context, task domain.Task) error {

	bytesBuffer := new(bytes.Buffer)
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = bytesBuffer.Write(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, "http://task-service:8080/commit", bytesBuffer)
	if err != nil {
		return err
	}
	resp, err := hs.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}
