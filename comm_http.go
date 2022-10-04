package jadesdk

import (
	"fmt"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	// "github.com/ant0ine/go-json-rest/rest"
	ds "uta.edu/aces/jadesdk/data_structure"
	"io/ioutil"
)

const (
	RetryIntervalNormal int = 30
	RetryIntervalShort      = 3
)

func (j *JadeSDK) retryHTTPCommunication(
	op string, protocol string, method string, path string, targetNode *ds.Node, 
	payload interface{}, logMsg string, seconds int, 
	retryCnt int, retrylimitation int,
) (interface{}, int, []byte, error) {
	if retrylimitation >= 0 && retryCnt > retrylimitation {
		return nil, 0, nil, errors.New(logMsg)
	}
	if logMsg != "" {
		j.log.Println(logMsg, "; going to retry ["+op+"] for the {"+strconv.Itoa(retryCnt)+"}th time in", seconds, "seconds...")
	}
	time.Sleep(time.Second * time.Duration(seconds))
	return j.HTTPCommunicate(op, protocol, method, path, targetNode, payload, retryCnt, retrylimitation)
}

// HTTPCommunicate access target node, if retry > 0, then finite retry mode, if retry == 0, then disable retry mode, if retry < 0, then infinite retry mode
// return value is 
//     the response, bytes of body, 
//     the timestamp of the first time it perform a sending action, before retrying if any
//     the successful transmission duration
//     error
func (j *JadeSDK) HTTPCommunicate(
	operationName string, protocol string, method string, path string, targetNode *ds.Node, 
	payload interface{}, 
	retryCnt int, retryLimitation int,
) (interface{}, int, []byte, error) {
	if targetNode == nil || (strings.ToLower(protocol) != "http" && strings.ToLower(protocol) != "https") {
		msg := "ERROR: invalid target node for " + operationName
		j.log.Println(msg)
		return nil, 0, nil, errors.New(msg)
	}
	retryInterval := RetryIntervalShort

	tailstr := ""
	if retryCnt > 0 {
		tailstr = ", retry count: " + strconv.Itoa(retryCnt)
	}
	j.log.Println("[comm] <"+operationName+"> started toward target node:", targetNode.Key(), tailstr)

	reqbody, err := json.Marshal(payload)
	reqLength := 0
	if err != nil {
		msg := "ERROR during encoding payload: " + err.Error()
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	}
	reqLength = len(reqbody)
	targetURL := protocol + "://" + targetNode.Key() + path
	if strings.Contains(targetNode.Key(), ":") {
		targetURL = targetNode.Key() + path
	}
	j.log.Printf("[comm] <"+operationName+"> target url: %v, method: %v", targetURL, method)
	// Send the register information to upper node
	req, err := http.NewRequest(strings.ToUpper(method), targetURL, bytes.NewBuffer(reqbody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("retry-count", strconv.Itoa(retryCnt))
	client := &http.Client{
		// Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		msg := "Error when sending http request: " + err.Error()
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	}

	if res == nil || res.Body == nil {
		msg := "Error of the communication for the response is nil"
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	}

	// parse the response message of upper node for registering
	resMsg := &Response{}
	// json.NewDecoder(res.Body).Decode(&resMsg)
	content, err := DecodeRequestWithoutClosing(res, resMsg)
	if err != nil {
		msg := fmt.Sprintf("can not decode response body: %v", err)
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	} else if resMsg.Error != "" {
		msg := fmt.Sprintf("target node responded ERROR message: %v", resMsg.Error)
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	} else if resMsg.Status != "OK" {
		msg := "target node responded abnormal status: " + resMsg.Status
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	}
	j.log.Println("[comm] "+operationName+" completed with target node:", targetNode.Key())
	return resMsg.Payload, reqLength, content, nil
}

func DecodeRequestWithoutClosing(r *http.Response, v interface{}) ([]byte, error) {
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("JSON payload is empty")
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return nil, err
	}
	return content, nil
}


