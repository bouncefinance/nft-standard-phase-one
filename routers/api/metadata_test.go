package api

import (
	"Ankr-gin-ERC721/pkg/util"
	"encoding/json"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	tokens1155 := []int{1, 2, 2, 3, 3, 3, 6}

	for i := 0; ; {
		if tokens1155[i+1] == tokens1155[i] {
			tokens1155 = append(tokens1155[:i], tokens1155[i+1:]...)
			continue
		}
		i++
		if i >= len(tokens1155)-1 {
			break
		}
	}

	fmt.Println(tokens1155)
}

func Test_(t *testing.T) {
	var dataI interface{}
	data ,err:= util.GetUrl("https://shalomhu.github.io/bounce/NFT/tokenURL_04.json", nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(data))
	err = json.Unmarshal(data, &dataI)
	if err != nil {
		t.Error(err)
		return
	}
	d, ok := dataI.(map[string]interface{})
	if !ok {
		t.Error(ok)
		return
	}
	d1, ok := d["properties"].(map[string]interface{})
	if !ok {
		t.Error(ok)
		return
	}
	dName, ok := d1["name"].(map[string]interface{})
	if !ok {
		t.Error(ok)
		return
	}

	dDescription, ok := d1["description"].(map[string]interface{})
	if !ok {
		t.Error(ok)
		return
	}

	dImage, ok := d1["image"].(map[string]interface{})
	if !ok {
		t.Error(ok)
		return
	}

	r := make(map[string]interface{})
	r["name"] = dName["description"]
	r["description"] = dDescription["description"]
	r["image"] = dImage["description"]
}
