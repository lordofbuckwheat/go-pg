package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	//http.HandleFunc("/api/{version}/screen/{screen_code}", identixone)
	http.HandleFunc("/api/v4/screen/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type Content struct {
	Id  uint64
	Url string
}

var content = []Content{
	{1, "http://demo.rerotor.ru/media/movies/Spain_Patagonia_cgtxKEr_dFuUOZf.m4v?channel=s361"},
	{2, "http://public.tvbit.co/dummy_videos/Sample1280.mp4"},
	{3, "http://public.tvbit.co/dummy_videos/small.mp4"},
	{4, "http://public.tvbit.co/dummy_videos/file_example_MP4_640_3MG.mp4"},
	{5, "http://public.tvbit.co/dummy_videos/115mb Barcelona in 4K.mp4"},
	{6, "http://public.tvbit.co/dummy_videos/1mb_big_buck_bunny_720p.mp4"},
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		panic(err)
	} else if len(body) > 0 {
		type RequestBody struct {
			Stats []struct{
				EventDate string `json:"event_date"`
				EventCount int64 `json:"event_count"`
				ContentId uint64 `json:"content_id"`
			} `json:"stats"`
		}
		var requestBody RequestBody
		if err := json.Unmarshal(body, &requestBody); err != nil {
			panic(err)
		}
		fmt.Printf("request stats %v\n", requestBody)
	}
	var parts = strings.Split(r.URL.Path, "/")
	var screenCode = parts[len(parts)-1]
	var videoUrls = func() []Content {
		if len(screenCode) <= 6 || screenCode[:6] != "screen" {
			return content
		}
		i, err := strconv.ParseUint(screenCode[6:], 10, 64)
		if err != nil {
			return content
		}
		switch i {
		case 1:
			return []Content{content[0]}
		case 2:
			return []Content{content[1]}
		case 3:
			return []Content{content[2]}
		case 4:
			return []Content{content[3]}
		case 5:
			return []Content{content[4]}
		case 6:
			return []Content{content[5]}
		case 7:
			return []Content{content[0],content[1]}
		case 8:
			return []Content{content[1],content[2]}
		case 9:
			return []Content{content[2],content[3]}
		case 10:
			return []Content{content[3],content[4]}
		case 11:
			return []Content{content[4],content[5]}
		default:
			return content
		}
	}()
	var responseRaw = []byte(`{
		"binding": {
			"screen": {
				"showcase": 361,
				"master_screen_code": null,
				"type": "TVLinux",
				"pixelSize": {
					"id": 2,
					"name": "FullHD",
					"icon": "/static/rerotor2/icons/horizontal-content.png",
					"width": 1920,
					"height": 1080,
					"unit": "px"
				},
				"physicalSize": null,
				"code": "demo-screen-2k",
				"device_id": "na",
				"device_name": "na",
				"vendor": "na",
				"os_version": "na",
				"mac_address": "AA:AA:AA:AA:AA:AA",
				"last_update": "2019-10-17T10:36:01.600622+03:00"
			},
			"pricetag": null,
			"group_order": 0
		},
		"product_info": {
			"price": 0,
			"motivation": "",
			"motivation_end": null,
			"options": [],
			"content": [],
			"special_content": [],
			"logos": [],
			"special_offers": [],
			"all_content": null,
			"group": null,
			"price_desc": "64990",
			"commands": [],
			"update_schedule": null,
			"showcase_info": null,
			"server_time_utc": "2019-10-17T07:36:05.249982Z",
			"wifi_list": [],
			"takeaway_state": -1,
			"price_offer": null,
			"features": [],
			"code": 10050002,
			"name": "Телевизор Samsung UE55NU7670U",
			"suspended": false,
			"on_order": false,
			"manual": false,
			"suggested_screen_brightness": 100,
			"is_blacklisted": false,
			"device_classes": [
				9
			]
		}
	}`)
	var response interface{}
	_ = json.Unmarshal(responseRaw, &response)
	var content = make([]interface{}, 0)
	fmt.Println("video urls", videoUrls)
	for _, vu := range videoUrls {
		var chunk interface{}
		_ = json.Unmarshal([]byte(fmt.Sprintf(`{
			"id": %d,
			"links": [
				"http://demo.rerotor.ru/media/movies/Spain_Pat", "%s"
			],
			"start": "2019-02-14T00:00:00+03:00",
			"end": "2019-11-07T23:59:59+03:00",
			"name": "Spain_Patagonia",
			"checksum": "6f3cc437751b1d4c68667efd8d89dfd1",
			"type": 0,
			"fillmode": 0,
			"playorder": 7
		}`, vu.Id, vu.Url)), &chunk)
		content = append(content, chunk)
	}
	response.(map[string]interface{})["product_info"].(map[string]interface{})["content"] = content
	responseRaw, _ = json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseRaw)
}
