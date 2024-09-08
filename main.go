package main

import (
	"fmt"
	"galgame-ranking-crawler/gameInfos"
	"sort"
	"strconv"
	"time"
)

func main() {
	url := "https://www.ymgal.games/rank?page="
	var gameInfoList []gameInfos.GameInfo
	for i := 0; i <= 35; i++ {
		time.Sleep(100 * time.Millisecond)

		gameInfo, _ := gameInfos.Extract(url + strconv.Itoa(i))
		gameInfoList = append(gameInfoList, gameInfo...)
	}

	sort.Slice(gameInfoList, func(i, j int) bool {
		return gameInfoList[i].Date.After(gameInfoList[j].Date)
	})

	// 打印美化后的切片长度
	fmt.Printf("Length: %v\n", len(gameInfoList))
	// 用制表符且等距再打印
	for _, gameInfo := range gameInfoList {
		fmt.Printf("%s\t%s\t%v\t%v\t%v\n", gameInfo.Date.Format("2006-01-02"), gameInfo.Title, gameInfo.Rank, gameInfo.Point, gameInfo.Chinese)
	}
}
