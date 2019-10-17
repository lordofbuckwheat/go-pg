package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
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

var videos = []string{
	"http://demo.rerotor.ru/media/movies/Spain_Patagonia_cgtxKEr_dFuUOZf.m4v?channel=s361",
	"http://public.tvbit.co/dummy_videos/Sample1280.mp4",
	"http://public.tvbit.co/dummy_videos/small.mp4",
	"http://public.tvbit.co/dummy_videos/file_example_MP4_640_3MG.mp4",
	"http://public.tvbit.co/dummy_videos/115mb Barcelona in 4K.mp4",
	"http://public.tvbit.co/dummy_videos/1mb_big_buck_bunny_720p.mp4",
}

var contentId = 0

func echo(w http.ResponseWriter, r *http.Request) {
	var parts = strings.Split(r.URL.Path, "/")
	var screenCode = parts[len(parts)-1]
	var videoUrls = func() []string {
		switch screenCode {
		case "demo-screen-2k":
			return []string{videos[0], videos[1]}
		case "1":
			return []string{videos[2], videos[3]}
		case "2":
			return []string{videos[4], videos[5]}
		default:
			return []string{videos[0]}
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
	for _, vu := range videoUrls {
		var chunk interface{}
		contentId++
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
		}`, contentId, vu)), &chunk)
		content = append(content, chunk)
	}
	response.(map[string]interface{})["product_info"].(map[string]interface{})["content"] = content
	responseRaw, _ = json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseRaw)
}