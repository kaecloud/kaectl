package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kaecloud/kaectl/pkg/spec"
	"github.com/pkg/errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type JobClient struct {
	Client
}

type Job struct {
	Id       int64  `json:"id"`
	Name     string `orm:"unique" json:"name"`
	SpecText string `json:"spec_text"`
	Comment  string `json:"comment"`
}

func NewJobClient(baseUrl string, accessTok string) *JobClient {
	c := NewClient()
	c.baseUrl = baseUrl
	c.accessToken = accessTok
	return &JobClient{
		Client: *c,
	}
}

func (c *JobClient) Get(name string) (*Job, error) {
	path := fmt.Sprintf("/api/v1/jobs/%s", name)
	var data Job
	err := c.REST("GET", path, nil, &data)
	return &data, err
}

func (c *JobClient) List() ([]*Job, error) {
	path := "/api/v1/jobs"
	var data []*Job
	err := c.REST("GET", path, nil, &data)
	return data, err
}

func (c *JobClient) Delete(name string) error {
	path := fmt.Sprintf("/api/v1/jobs/%s", name)
	var data Job
	err := c.REST("DELETE", path, nil, &data)
	return err
}

func (c *JobClient) Create(args *spec.CreateJobArgs) (*Job, error) {
	path := "/api/v1/jobs"
	reqBytes, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewReader(reqBytes)
	res := Job{}

	err = c.REST("POST", path, reqBody, &res)
	return &res, err
}

func (c *JobClient) Upload(jobname string, values map[string]io.Reader, respBody interface{}) (err error) {
	path := fmt.Sprintf("/api/v1/jobs/%s/artifacts", jobname)
	uploadURL := c.FullUrl(path)

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", uploadURL, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	return c.rest(req, respBody)
}

func (c *JobClient) UploadArtifact(jobname string, filename string, objKey string) (objUrl string, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	values := map[string]io.Reader{
		"fileUploadName": r,
	}
	if objKey != "" {
		values["objectKey"] = strings.NewReader(objKey)
	}
	var jsonResp struct {
		Data string `json:"data"`
	}
	err = c.Upload(jobname, values, &jsonResp)
	return jsonResp.Data, err
}

func (c *JobClient) Logs(name string, podname string, cluster string, follow bool) (chan interface{}, error) {
	path := fmt.Sprintf("/api/v1/jobs/%s/log?cluster=%s&podname=%s", name, cluster, podname)
	if follow {
		path += "&follow=true"
	}
	wsUrl := c.FullWebsocketUrl(path)
	hdr := http.Header{
		"Authorization": {"Bearer "+c.accessToken},
	}
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, hdr)
	if err != nil {
		log.Fatal("dial:", err)
	}
	res := make(chan interface{})

	go func() {
		defer func() {
			ws.Close()
			close(res)
		}()
		for {
			msgType, msgBytes, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				res <- err
				return
			}
			switch msgType {
			case websocket.PingMessage:
				if err := ws.WriteMessage(websocket.PongMessage, []byte("PP")); err != nil {
					res <- err
					return
				}
			case websocket.TextMessage:
				var data struct {
					Data string `json:"data"`
					Error string `json:"error"`
				}
				err = json.Unmarshal(msgBytes, &data)
				if err != nil {
					res <- err
					return
				}
				if data.Error != "" {
					res <- errors.New(data.Error)
				} else {
					res <- data.Data
				}
			}
		}
	}()
	return res, err
}

