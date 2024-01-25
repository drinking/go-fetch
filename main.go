package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

type Result struct {
	Value Data `json:"result"`
}

type Data struct {
	Value Candidate `json:"data"`
}

type Candidate struct {
	Status   int `json:"signStatus"`
	LimitNum int `json:"limitNum"`
	SignNumi int `json:"signNum"`
}

func fetch(url string) (Candidate, error) {
	client := &http.Client{}
	var data = strings.NewReader(`{"shopPath":"pdszqyg","shopProvince":"sh","trainId":"3e3e92356e5740bd8b6f661040aff64d"}`)
	req, err := http.NewRequest("POST", "https://sh.train.wenhuayun.cn/api/train/artTrain/artTrainDetail", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "sh.train.wenhuayun.cn")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9,zh-CN;q=0.8,zh-TW;q=0.7,zh;q=0.6")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("jwttoken", "")
	req.Header.Set("origin", "https://sh.train.wenhuayun.cn")
	req.Header.Set("referer", "https://sh.train.wenhuayun.cn/pdszqyg/cloud-train/train-detail?trainId=3e3e92356e5740bd8b6f661040aff64d")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", `"Android"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sysplatform", "wap")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// bs, _ := json.MarshalIndent(body, "", "    ")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Candidate{}, err
	}

	return result.Value.Value, nil
}

func notify(content string) {
	note := gosxnotifier.NewNotification("Hello")

	// Set the notification title and subtitle
	note.Title = content
	note.Link = "https://sh.train.wenhuayun.cn/pdszqyg/cloud-train/train-detail?trainId=3e3e92356e5740bd8b6f661040aff64d"

	// Set the notification message and sound
	note.Message = "This is a side notification"
	note.Sound = gosxnotifier.Default

	// Set the notification group to override any previous notification
	note.Group = "com.example.notification"

	// Show the notification
	err := note.Push()
	if err != nil {
		fmt.Println("Error displaying notification:", err)
		return
	}
}

// func sendToSlack() {
// 	// Create a new Slack client
// 	api := slack.New("YOUR_SLACK_API_TOKEN")

// 	// Set the channel and message
// 	channelID := "YOUR_CHANNEL_ID"
// 	message := "Hello, Slack!"

// 	// Send the notification
// 	_, _, err := api.PostMessage(
// 		channelID,
// 		slack.MsgOptionText(message, false),
// 	)
// 	if err != nil {
// 		log.Fatalf("Error sending Slack notification: %s", err)
// 	} else {
// 		fmt.Println("Slack notification sent successfully!")
// 	}
// }

func sendToSlack(msg string) {
	SLACK_CHANEEL := os.Getenv("SLACK_CHANEEL")
	client := &http.Client{}
	message := fmt.Sprintf(`{"text":"%s"}`, msg)
	var data = strings.NewReader(message)
	req, err := http.NewRequest("POST", SLACK_CHANEEL, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}

func main() {

	for {
		result, _ := fetch("https://sh.train.wenhuayun.cn/api/train/artTrain/artTrainDetail")
		if result.LimitNum != result.SignNumi || result.Status != 4 {
			// notify("Fetch completed")
			sendToSlack("https://sh.train.wenhuayun.cn/pdszqyg/cloud-train/train-detail?trainId=3e3e92356e5740bd8b6f661040aff64d")
			break
		} else {
			fmt.Println("fetching...", result)
		}
		time.Sleep(60 * time.Second)
	}

}
