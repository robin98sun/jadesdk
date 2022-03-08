package jadesdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	RetryIntervalNormal int = 30
	RetryIntervalShort      = 3
)

func (j *JadeSDK) retryHTTPCommunication(
	op string, protocol string, method string, path string, targetNode *Node, 
	payload interface{}, logMsg string, seconds int, 
	retryCnt int, retrylimitation int,
) (interface{}, int, error) {
	if retrylimitation >= 0 && retryCnt > retrylimitation {
		return nil, 0, errors.New(logMsg)
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
	operationName string, protocol string, method string, path string, targetNode *Node, 
	payload interface{}, 
	retryCnt int, retryLimitation int,
) (interface{}, int, error) {
	if targetNode == nil || (strings.ToLower(protocol) != "http" && strings.ToLower(protocol) != "https") {
		msg := "ERROR: invalid target node for " + operationName
		j.log.Println(msg)
		return nil, 0, errors.New(msg)
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
	// Send the register information to upper node
	req, err := http.NewRequest(strings.ToUpper(method), targetURL, bytes.NewBuffer(reqbody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("retry-count", strconv.Itoa(retryCnt))
	client := &http.Client{
		Timeout: 30 * time.Second,
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
	resMsg := Response{}
	json.NewDecoder(res.Body).Decode(&resMsg)
	if resMsg.Error != "" {
		msg := "target node responded ERROR message: " + resMsg.Error
		return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
	} else {
		if resMsg.Status == "OK" {
			j.log.Println("[comm] "+operationName+" completed with target node:", targetNode.Key())
			return resMsg.Payload, reqLength, nil
		} else {
			msg := "target node responded abnormal status: " + resMsg.Status
			return j.retryHTTPCommunication(operationName, protocol, method, path, targetNode, payload, msg, retryInterval, retryCnt+1, retryLimitation)
		}
	}
}
