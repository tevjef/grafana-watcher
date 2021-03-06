// Copyright 2016 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grafana

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Interface interface {
	Dashboards() DashboardsInterface
	Datasources() DatasourcesInterface
}

type Clientset struct {
	BaseUrl    *url.URL
	HTTPClient *http.Client
}

func New(baseUrl *url.URL) Interface {
	return &Clientset{
		BaseUrl:    baseUrl,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Clientset) Dashboards() DashboardsInterface {
	return NewDashboardsClient(c.BaseUrl, c.HTTPClient)
}

func (c *Clientset) Datasources() DatasourcesInterface {
	return NewDatasourcesClient(c.BaseUrl, c.HTTPClient)
}

func doRequest(c *http.Client, req *http.Request) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		return fmt.Errorf("Unexpected status code returned from Grafana API (got: %d, expected: 200, msg:%s)", resp.StatusCode, resp.Status)
	}
	return nil
}
