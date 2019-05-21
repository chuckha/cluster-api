package machine

import (
	"bytes"
	"context"
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

type MachineRequest struct {
	Cluster []byte
	Machine []byte
}

func (c *Client) Delete(_ context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	cb, _ := json.Marshal(cluster)
	mb, _ := json.Marshal(machine)
	b, err := json.Marshal(&MachineRequest{
		Cluster: cb,
		Machine: mb,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	klog.Infoln("calling delete...")
	klog.Infoln(string(b))
	resp, err := http.Post(fmt.Sprintf("http://%s.%s:%s/delete", c.Config.Service, c.Config.Namespace, c.Config.Port), "application/json", bytes.NewReader(b))
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

func (c *Client) Exists(_ context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	cb, _ := json.Marshal(cluster)
	mb, _ := json.Marshal(machine)
	b, err := json.Marshal(&MachineRequest{
		Cluster: cb,
		Machine: mb,
	})
	if err != nil {
		return false, errors.WithStack(err)
	}
	klog.Infoln("calling exists...")
	klog.Infoln(string(b))
	resp, err := http.Post(fmt.Sprintf("http://%s.%s:%s/exists", c.Config.Service, c.Config.Namespace, c.Config.Port), "application/json", bytes.NewReader(b))
	if err != nil {
		return false, errors.WithStack(err)
	}
	klog.Infoln(resp)
	if resp.StatusCode >= 500 {
		return false, errors.New("internal server error")
	}
	if resp.StatusCode >= 400 {
		return false, errors.New("bad request")
	}
	return false, nil
}

func (c *Client) Update(_ context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	cb, _ := json.Marshal(cluster)
	mb, _ := json.Marshal(machine)
	b, err := json.Marshal(&MachineRequest{
		Cluster: cb,
		Machine: mb,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	klog.Infoln("calling updte...")
	klog.Infoln(string(b))
	resp, err := http.Post(fmt.Sprintf("http://%s.%s:%s/update", c.Config.Service, c.Config.Namespace, c.Config.Port), "application/json", bytes.NewReader(b))
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

func (c *Client) Create(_ context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	cb, _ := json.Marshal(cluster)
	mb, _ := json.Marshal(machine)
	b, err := json.Marshal(&MachineRequest{
		Cluster: cb,
		Machine: mb,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	klog.Infoln("calling create...")
	klog.Infoln(string(b))
	resp, err := http.Post(fmt.Sprintf("http://%s.%s:%s/create", c.Config.Service, c.Config.Namespace, c.Config.Port), "application/json", bytes.NewReader(b))
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
