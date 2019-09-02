package influxdb

import (
	"fmt"
	"go-log-analysis/mysql"
	"io/ioutil"
	"net/http"
	"strings"
)

const InfluxdbUrl  = "http://47.94.169.212:8086/write?db=log"


//curl -i -XPOST 'http://47.94.169.212:8086/write?db=log' --data-binary 'cpu_load_short,host=server01,region=us-west value=0.64 1434055562000000000'
func InsertData(data mysql.LogInfo)  {

	client := &http.Client{}

	dataInsert:="nginx_log,ip="+data.Ip +",method="+data.Method +" path=22"

	req, err := http.NewRequest("POST", InfluxdbUrl, strings.NewReader(dataInsert))
	if err != nil {
		fmt.Println("error post:",err)
	}

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	if err!=nil{
		fmt.Println("Do error:",err)
		return
	}

	//defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error read:",err)
	}
	fmt.Println("4444:",resp)
	fmt.Println("5555:",body)
}