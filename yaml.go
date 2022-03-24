package k8s_helm

import (
	"encoding/base64"
	"errors"
	"fmt"
	jsonYaml "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func yamlToMap(fileByte []byte) (resultMap map[interface{}]interface{}, err error) {
	// base64 编码
	encodeString := base64.StdEncoding.EncodeToString(fileByte)
	// base64 解码
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(decodeBytes, &resultMap); err != nil {
		return nil, err
	}
	return
}

func ReplaceYaml(filePath string, params map[string]string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if string(yamlFile) == "" || string(yamlFile) == "{}" {
		return nil
	}
	if err != nil {
		return err
	}
	valueMap, err := yamlToMap(yamlFile)
	if err != nil {
		return nil
	}
	for key, value := range params {
		valueMap, err = replaceJsonValue(valueMap, value, key)
		if err != nil {
			return err
		}
	}
	// base64 编码
	marshal, err := yaml.Marshal(valueMap)
	if err != nil {
		fmt.Println(err)
		return err
	}

	yamlData, err := jsonYaml.JSONToYAML(marshal)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, yamlData, 0777)
	return nil
}

func oneParam(data map[interface{}]interface{}, value interface{}, param []string) (res map[interface{}]interface{}, err error) {
	_, ok := data[param[0]]
	if !ok {
		return nil, errors.New("step.yaml参数与value.yaml不匹配")
	}
	data[param[0]] = value
	return data, nil
}

func twoParam(data map[interface{}]interface{}, value interface{}, param []string) (res map[interface{}]interface{}, err error) {
	_, ok := data[param[0]].(map[interface{}]interface{})[param[1]]
	if !ok {
		return nil, errors.New("step.yaml参数与value.yaml不匹配")
	}
	data[param[0]].(map[interface{}]interface{})[param[1]] = value
	return data, nil
}

func threeParam(data map[interface{}]interface{}, value interface{}, param []string) (res map[interface{}]interface{}, err error) {
	_, ok := data[param[0]].(map[interface{}]interface{})[param[1]].(map[interface{}]interface{})[param[2]]
	if !ok {
		return nil, errors.New("step.yaml参数与value.yaml不匹配")
	}
	data[param[0]].(map[interface{}]interface{})[param[1]].(map[interface{}]interface{})[param[2]] = value
	return data, nil
}

func fourParam(data map[interface{}]interface{}, value interface{}, param []string) (res map[interface{}]interface{}, err error) {
	_, ok := data[param[0]].(map[interface{}]interface{})[param[1]].(map[interface{}]interface{})[param[2]].(map[interface{}]interface{})[param[3]]
	if !ok {
		return nil, errors.New("step.yaml参数与value.yaml不匹配")
	}
	data[param[0]].(map[interface{}]interface{})[param[1]].(map[interface{}]interface{})[param[2]].(map[interface{}]interface{})[param[3]] = value
	return data, nil
}

func fiveParam(data map[interface{}]interface{}, value interface{}, param []string) (res map[interface{}]interface{}, err error) {
	_, ok := data[param[0]].(map[interface{}]interface{})[param[1]].(map[interface{}]interface{})[param[2]].(map[interface{}]interface{})[param[3]].(map[interface{}]interface{})[param[4]]
	if !ok {
		return nil, errors.New("step.yaml参数与value.yaml不匹配")
	}
	data[param[0]].(map[interface{}]interface{})[param[1]].(map[interface{}]interface{})[param[2]].(map[interface{}]interface{})[param[3]].(map[interface{}]interface{})[param[4]] = value
	return data, nil
}

func replaceJsonValue(resultMap map[interface{}]interface{}, value string, key string) (res map[interface{}]interface{}, err error) {
	params := StringToSlice(key, ".")
	switch len(params) {
	case 1:
		res, err = oneParam(resultMap, value, params)
	case 2:
		res, err = twoParam(resultMap, value, params)
	case 3:
		res, err = threeParam(resultMap, value, params)
	case 4:
		res, err = fourParam(resultMap, value, params)
	case 5:
		res, err = fiveParam(resultMap, value, params)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func StringToSlice(str string, tag string) []string {
	return strings.Split(str, tag)
}
