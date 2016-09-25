package gogobosh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetStemcells from given BOSH
func (c *Client) GetStemcells() (stemcells []Stemcell, err error) {
	r := c.NewRequest("GET", "/stemcells")
	resp, err := c.DoRequest(r)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error requesting stemcells  %v", err)
		return
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading stemcells request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &stemcells)
	if err != nil {
		log.Printf("Error unmarshaling stemcells %v", err)
		return
	}
	return
}

// UploadStemcell to given BOSH
func (c *Client) UploadStemcell(stemcellURI string, sha1 *string) (task Task, err error) {
	task, err = c.upload("/stemcells", stemcellURI, sha1)
	return
}

// UploadRelease to given BOSH
func (c *Client) UploadRelease(releaseURI string, sha1 *string) (task Task, err error) {
	task, err = c.upload("/releases", releaseURI, sha1)
	return
}

// Upload a file to given BOSH path
func (c *Client) upload(path string, fileURI string, sha1 *string) (task Task, err error) {

	var (
		isRemote bool
		resp     *http.Response
		body     []byte
	)

	if isRemote, err = regexp.MatchString("^http(s)?://", fileURI); !isRemote {
		if _, err = os.Stat(fileURI); os.IsNotExist(err) {
			err = fmt.Errorf("Stemcell file '%s' does not exist.", fileURI)
		}
	}
	if err != nil {
		return
	}

	r := c.NewRequest("POST", path)
	if isRemote {
		r.header["Content-Type"] = "application/json"
		body, err = json.Marshal(struct {
			Location string  `json:"location"`
			Sha1     *string `json:"sha1"`
		}{
			fileURI,
			sha1,
		})
		if err != nil {
			return
		}
	} else {
		r.header["Content-Type"] = "application/x-compressed"
		body, err = ioutil.ReadFile(fileURI)
		if err != nil {
			return
		}
	}
	buffer := bytes.NewBuffer(body)
	r.body = buffer

	resp, err = c.DoRequest(r)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Bosh director responded with error: %s", resp.Status)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
		return
	}
	return
}

// GetReleases from given BOSH
func (c *Client) GetReleases() (releases []Release, err error) {
	r := c.NewRequest("GET", "/releases")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting releases  %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading releases request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &releases)
	if err != nil {
		log.Printf("Error unmarshaling releases %v", err)
		return
	}
	return
}

// GetDeployments from given BOSH
func (c *Client) GetDeployments() (deployments []Deployment, err error) {
	r := c.NewRequest("GET", "/deployments")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deployments  %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployments request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &deployments)
	if err != nil {
		log.Printf("Error unmarshaling deployments %v", err)
		return
	}
	return
}

// GetDeployment from given BOSH
func (c *Client) GetDeployment(name string) (manifest Manifest, err error) {
	r := c.NewRequest("GET", "/deployments/"+name)
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deployment manifest %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployment manifest request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &manifest)
	if err != nil {
		log.Printf("Error unmarshaling deployment manifest %v", err)
		return
	}
	return
}

// DeleteDeployment from given BOSH
func (c *Client) DeleteDeployment(name string) (task Task, err error) {
	r := c.NewRequest("DELETE", "/deployments/"+name+"?force=true")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deleting deployment %v", err)
		return
	}
	defer resp.Body.Close()
	url, _ := resp.Location()
	re, _ := regexp.Compile(`(\d+)$`)
	stringId := re.FindStringSubmatch(url.Path)
	id, err := strconv.Atoi(stringId[0])
	if err != nil {
		return
	}
	task, err = c.GetTask(id)
	return
}

// CreateDeployment from given BOSH
func (c *Client) CreateDeployment(manifest string) (task Task, err error) {
	r := c.NewRequest("POST", "/deployments")
	buffer := bytes.NewBufferString(manifest)
	r.body = buffer
	r.header["Content-Type"] = "text/yaml"

	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting create deployment %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
		return
	}
	return
}

// GetDeploymentVMs from given BOSH
func (c *Client) GetDeploymentVMs(name string) (vms []VM, err error) {
	var task Task
	r := c.NewRequest("GET", "/deployments/"+name+"/vms?format=full")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deployment vms %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployment vms request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling tasks %v", err)
		return
	}
	for {
		taskStatus, err := c.GetTask(task.ID)
		if err != nil {
			log.Printf("Error getting task %v", err)
		}
		if taskStatus.State == "done" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	output := c.GetTaskResult(task.ID)
	for _, value := range output {
		if len(value) > 0 {
			var vm VM
			err = json.Unmarshal([]byte(value), &vm)
			if err != nil {
				log.Printf("Error unmarshaling vms %v %v", value, err)
				return
			}
			vms = append(vms, vm)
		}
	}
	return
}

// GetTasks from given BOSH
func (c *Client) GetTasks() (tasks []Task, err error) {
	r := c.NewRequest("GET", "/tasks")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting tasks  %v", err)
		return
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading tasks request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &tasks)
	if err != nil {
		log.Printf("Error unmarshaling tasks %v", err)
		return
	}
	return
}

// GetTask from given BOSH
func (c *Client) GetTask(id int) (task Task, err error) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID)
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting task %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
		return
	}
	return
}

// GetTaskResult from given BOSH
func (c *Client) GetTaskResult(id int) (output []string) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID+"/output?type=result")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting task %v", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
	}
	output = strings.Split(string(resBody), "\n")
	return
}
