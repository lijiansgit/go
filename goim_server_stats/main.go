package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/tidwall/gjson"
)

type Stats struct {
	id    int
	count int
}

type StatsList []Stats

func (s StatsList) Len() int {
	return len(s)
}
func (s StatsList) Less(i, j int) bool {
	return s[i].count < s[j].count
}
func (s StatsList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	var (
		serverStats StatsList
		roomStats   StatsList
	)
	if len(os.Args) < 2 {
		fmt.Println("Usage: goim_server_stats logic")
		return
	}

	serverStatsUrl := "http://" + os.Args[1] + ":7172/1/count?type=server"
	roomStatsUrl := "http://" + os.Args[1] + ":7172/1/count?type=room"

	// comet server count
	serverStats, totalCount, err := get(serverStatsUrl, "Server")
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i <= 100; i++ {
		for _, stats := range serverStats {
			if stats.count != 0 && stats.id == i {
				fmt.Printf("comet-%d: %d\n", stats.id, stats.count)
			}
		}
	}
	fmt.Println("")

	// room count
	roomStats, _, err = get(roomStatsUrl, "RoomId")
	if err != nil {
		log.Fatalln(err)
	}
	sort.Sort(sort.Reverse(roomStats))
	for i := 0; i <= 10; i++ {
		if roomStats[i].count != 0 {
			fmt.Printf("room-%d: %d\n", roomStats[i].id, roomStats[i].count)
		}
	}
	fmt.Println("")

	fmt.Printf("TotalCount: %d\n", totalCount)
}

func get(url, keyType string) (stats []Stats, totalCount int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return stats, totalCount, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stats, totalCount, err
	}
	res := string(respBody)
	resCode := gjson.Get(res, "ret").Int()
	if resCode != 1 {
		return stats, totalCount, err
	}
	resData := gjson.Get(res, "data").Array()
	var count int
	stats = make([]Stats, 1000000)
	for k, v := range resData {
		ids := v.Get(keyType).String()
		id, err := strconv.Atoi(ids)
		if err != nil {
			return stats, totalCount, err
		}

		count = int(v.Get("Count").Int())
		stats[k].id = id
		stats[k].count = count

		totalCount = totalCount + count
	}

	return stats, totalCount, nil
}
