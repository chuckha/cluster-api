package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog"

	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type ClientConfig struct {
	Service   string
	Namespace string
	Port      string
}

type Client struct {
	Config *ClientConfig
}

func (c *Client) Reconcile(cluster *clusterv1.Cluster) error {
	b, err := json.Marshal(cluster)
	if err != nil {
		return errors.WithStack(err)
	}
	klog.Infoln("calling reconcile...")
	klog.Infoln(string(b))
	resp, err := http.Post(fmt.Sprintf("http://%s.%s:%s/reconcile", c.Config.Service, c.Config.Namespace, c.Config.Port), "application/json", bytes.NewReader(b))
	if err != nil {
		return errors.WithStack(err)
	}
	klog.Infoln(resp)
	if resp.StatusCode >= 500 {
		return errors.New("internal server error")
	}
	if resp.StatusCode >= 400 {
		return errors.New("bad request")
	}
	return nil
}

func (c *Client) Delete(cluster *clusterv1.Cluster) error {
	return nil
}
