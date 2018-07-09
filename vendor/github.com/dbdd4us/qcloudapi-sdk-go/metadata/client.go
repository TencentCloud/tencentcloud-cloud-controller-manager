package metadata

import (
	"errors"
	"fmt"
	. "github.com/dbdd4us/qcloudapi-sdk-go/util"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
}

const (
	ENDPOINT = "http://metadata.tencentyun.com/meta-data"

	INSTANCE_ID  = "instance-id"
	UUID         = "uuid"
	MAC          = "mac"
	PRIVATE_IPV4 = "local-ipv4"
	REGION       = "placement/region"
	ZONE         = "placement/zone"
	PUBLIC_IPV4  = "public-ipv4"
)

type IMetaDataClient interface {
	Resource(resource string) IMetaDataClient
	Go() (string, error)
	Url() (string, error)
}

type MetaData struct {
	c IMetaDataClient
}

func NewMetaData(client *http.Client) *MetaData {
	if client == nil {
		client = &http.Client{}
	}
	return &MetaData{
		c: &MetaDataClient{client: client},
	}
}

func (m *MetaData) UUID() (string, error) {

	uuid, err := m.c.Resource(UUID).Go()
	if err != nil {
		return "", err
	}
	return uuid, err
}

func (m *MetaData) InstanceID() (string, error) {

	instanceId, err := m.c.Resource(INSTANCE_ID).Go()
	if err != nil {
		return "", err
	}
	return instanceId, err
}

func (m *MetaData) Mac() (string, error) {

	mac, err := m.c.Resource(MAC).Go()
	if err != nil {
		return "", err
	}
	return mac, nil
}

func (m *MetaData) PrivateIPv4() (string, error) {

	ip, err := m.c.Resource(PRIVATE_IPV4).Go()
	if err != nil {
		return "", err
	}
	return ip, nil
}

func (m *MetaData) PublicIPv4() (string, error) {

	ip, err := m.c.Resource(PUBLIC_IPV4).Go()
	if err != nil {
		return "", err
	}
	return ip, nil
}

func (m *MetaData) Region() (string, error) {

	region, err := m.c.Resource(REGION).Go()
	if err != nil {
		return "", err
	}
	return region, nil
}

func (m *MetaData) Zone() (string, error) {

	zone, err := m.c.Resource(ZONE).Go()
	if err != nil {
		return "", err
	}
	return zone, nil
}

//
type MetaDataClient struct {
	resource string
	client   *http.Client
}

func (m *MetaDataClient) Resource(resource string) IMetaDataClient {
	m.resource = resource
	return m
}

func (m *MetaDataClient) Url() (string, error) {
	if m.resource == "" {
		return "", errors.New("the resource you want to visit must not be nil!")
	}
	return fmt.Sprintf("%s/%s", ENDPOINT, m.resource), nil
}

func (m *MetaDataClient) send() (string, error) {
	u, err := m.Url()
	if err != nil {
		return "", err
	}
	requ, err := http.NewRequest(http.MethodGet, u, nil)

	if err != nil {
		return "", err
	}
	resp, err := m.client.Do(requ)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

var retry = AttemptStrategy{
	Min:   5,
	Total: 5 * time.Second,
	Delay: 200 * time.Millisecond,
}

func (vpc *MetaDataClient) Go() (resu string, err error) {
	for r := retry.Start(); r.Next(); {
		resu, err = vpc.send()
		if !shouldRetry(err) {
			break
		}
	}
	return resu, err
}

type TimeoutError interface {
	error
	Timeout() bool // Is the error a timeout?
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(TimeoutError)
	if ok {
		return true
	}

	switch err {
	case io.ErrUnexpectedEOF, io.EOF:
		return true
	}
	switch e := err.(type) {
	case *net.DNSError:
		return true
	case *net.OpError:
		switch e.Op {
		case "read", "write":
			return true
		}
	case *url.Error:
		// url.Error can be returned either by net/url if a URL cannot be
		// parsed, or by net/http if the response is closed before the headers
		// are received or parsed correctly. In that later case, e.Op is set to
		// the HTTP method name with the first letter uppercased. We don't want
		// to retry on POST operations, since those are not idempotent, all the
		// other ones should be safe to retry.
		switch e.Op {
		case "Get", "Put", "Delete", "Head":
			return shouldRetry(e.Err)
		default:
			return false
		}
	}
	return false
}
