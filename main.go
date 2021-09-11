package main

import (
	"fmt"
)
//創建三個channel 即創建三個Queue
var (
	imagePathChan = make(chan string, 5)
	newImagePathChan = make(chan string, 5)
	exitChan = make(chan bool, 1)
)
//定義結構體(要處理的檔案個數)
type Im struct {
	ImagePath1 string
	ImagePath2 string
	ImagePath3 string
}

//寫入imagePathChan Channel 
func addToQueue(imagePathChan chan string, imagePath1 string, imagePath2 string, imagePath3 string) {
	//將路徑放入Queue中
	imagePathChan <-imagePath1
	imagePathChan <-imagePath2
	imagePathChan <-imagePath3
	close(imagePathChan)
}


//傳入想要照片的路徑(imagePath)，此照片將會被排除待處理的佇列中
func removeFromQueue(imagePathChan chan string, newImagePathChan chan string, ImagePath string) {
	//如果Queue中的值符合傳入的路徑將其從Queue中移除
	for {
		num, ok := <-imagePathChan
		//若imagePathChan為空跳出迴圈
		if !ok {
			break
		}
		//過濾imagePathChan寫入newimagePathChan
		if num != ImagePath {
			newImagePathChan <-num
		}
	}
	close(newImagePathChan)
}

//讀取newImagePathChan Channel與處理圖片
func imageProcess(imagePathChan chan string, exitChan chan bool) {

	for {
		v, ok := <-imagePathChan 
		if !ok {
			break
		}
		fmt.Printf("處理數據完成%v \n", v)
	}

	//處理完數據後 即任務完成後
	exitChan <-true
	close(exitChan)
}



func main() {
	//定義imagePath
	var imp = Im { 
		ImagePath1 : "C:/Users/TD/Pictures/001",
		ImagePath2 : "C:/Users/TD/Pictures/002",
		ImagePath3 : "C:/Users/TD/Pictures/003",
	}
	//goroutine並發處理
	go addToQueue(imagePathChan, imp.ImagePath1, imp.ImagePath2, imp.ImagePath3)
	go removeFromQueue(imagePathChan, newImagePathChan, imp.ImagePath2)
	go imageProcess(newImagePathChan, exitChan)
	//避免主執行緒提前結束，當exitChan讀不到資料，代表副執行緒執行結束，可以結束主執行緒。
	for {
		_, ok := <-exitChan
		if  !ok {
			break
		}
	}
}