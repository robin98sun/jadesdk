package jadesdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
)

func decodeRawRequestBody(r *rest.Request, v interface{}) ([]byte, error) {
	content, err := ioutil.ReadAll(r.Body)
	// r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, nil
	}
	if v != nil {
		err = json.Unmarshal(content, v)
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}

func decodeRequestFromBytes(reqBytes []byte) (*Request, error) {
	req := &Request{}
	err := json.Unmarshal(reqBytes, req)
	return req, err
}

func decodeInputFromBytes(bodyBytes []byte, v interface{}) error {
	req := &Request{
		Payload: v,
	}
	err := json.Unmarshal(bodyBytes, req)
	return err
}

func decodeRequest(moduleName string, input interface{}, r *rest.Request) (*Request, error) {
	// decode request body into bytes
	var msg string
	bodyBytes, err := decodeRawRequestBody(r, nil)
	if err != nil {
		msg = fmt.Sprintf("[%v] error when decoding request http package: %v", moduleName, err.Error())
	} else {
		// extract taskId, subtaskId, returnTo from request body
		req, err := decodeRequestFromBytes(bodyBytes)
		if err != nil {
			msg = fmt.Sprintf("[%v] error when decoding request from bytes: %v", moduleName, err.Error())
		} else {
			err = decodeInputFromBytes(bodyBytes, input)
			return req, err
		}
	}
	if msg != "" {
		return nil, errors.New(msg)
	}
	return nil, fmt.Errorf("[%v] unknown error during decoding request", moduleName)
}
