package SearchRegistry

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
	"sync"
)

type Result struct {
	Id       int    `json:"id"`
	Key      string `json:"key"`
	Path     string `json:"path"`
	Accuracy string `json:"accuracy"`
}

var Percentage chan string

var level string

var SearchChan chan Result

var lastPercentage float64

func init() {
	SearchChan = make(chan Result, 1)
	Percentage = make(chan string, 1)
}

func SearchRegistry(input string) {
	//start := time.Now()
	//color.Green("注册表残留软件信息检测开始,请稍等片刻--------------------------")
	//这里是windows的读取方式，如果是mac \n
	//keywords := strings.Split(input, "\r\n")
	keywords := strings.Split(input, "\n")
	for index, value := range keywords {
		fmt.Printf("Index: %d, Value: %s\n", index, value)
	}
	regData := make(map[string]string)
	// 初始进度为0%
	Percentage <- "0"
	lastPercentage = 0.0
	//读取software注册表进行初始化map
	initRegistryMap("SOFTWARE", regData)

	searchKeyInMap(keywords, regData)
	//close(Percentage)  //搜索完成后关闭进度条通道
}

func initRegistryMap(path string, regData map[string]string) {
	hives := []registry.Key{
		//SearchRegistry.CLASSES_ROOT,
		registry.CURRENT_USER,
		registry.LOCAL_MACHINE,
		//SearchRegistry.USERS,
		//SearchRegistry.CURRENT_CONFIG,
	}
	// 设置并发限制为2个goroutine
	concurrencyLimit := 2
	sem := make(chan struct{}, concurrencyLimit)

	// 创建一个 channel 用于从 goroutine 中接收结果
	resultCh := make(chan map[string]string)

	// 创建一个 goroutine 处理每个注册表键
	var wg sync.WaitGroup

	for _, hive := range hives {
		wg.Add(1)
		go func(h registry.Key) {
			defer wg.Done()
			// 等待信号量
			sem <- struct{}{}
			key, err := registry.OpenKey(h, path, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
			if err != nil {
				fmt.Println("Error opening key:", err)
				// 将空 map 发送到结果 channel 中以表示错误
				resultCh <- make(map[string]string)
				<-sem
				return
			}
			defer key.Close()
			result := make(map[string]string)
			searchRegistryToMap(h, key, path, result)
			// 发送结果到 channel
			resultCh <- result
			<-sem

		}(hive)
	}

	// 启动一个goroutine来收集所有goroutine的结果
	go func() {
		wg.Wait()
		close(resultCh)
		Percentage <- "100"
	}()
	// 处理结果
	for range hives {
		result := <-resultCh
		// 合并结果到主 map
		for k, v := range result {
			regData[k] = v
		}
	}

}

func searchRegistryToMap(hive registry.Key, key registry.Key, keyPath string, regData map[string]string) {
	// 读取当前键的值
	values, err := key.ReadValueNames(-1)
	if err != nil {
		fmt.Println("Error reading values:", err)
		return
	}
	for _, valueName := range values {
		val, _, err := key.GetStringValue(valueName)
		if err != nil {
			//fmt.Println("Error reading value:", err)
			continue
		}
		regData[registryKeyToString(hive)+keyPath+"\\"+valueName] = val
	}

	//如果子健数量为0，说明是最后一键 直接返回
	keyinfo, _ := key.Stat()
	if keyinfo.SubKeyCount == 0 {

		// 更新进度条
		currentPercentage := float64(len(regData)) / float64(500000) * 100
		if lastPercentage > currentPercentage {
			currentPercentage = lastPercentage
		}
		Percentage <- fmt.Sprintf("%.1f", currentPercentage)
		lastPercentage = currentPercentage
		return
	}

	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkeys:", err)
		return
	}

	for _, subKey := range subKeys {
		subKeyPath := keyPath + "\\" + subKey
		//fmt.Println(subKeyPath)

		subKeyHandle, err := registry.OpenKey(hive, subKeyPath, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
		if err != nil {
			//fmt.Println("Error opening subkey:", err)
			continue
		}
		defer subKeyHandle.Close()

		searchRegistryToMap(hive, subKeyHandle, subKeyPath, regData)
	}
}

func registryKeyToString(key registry.Key) string {
	switch key {
	case registry.CLASSES_ROOT:
		return "HKEY_CLASSES_ROOT\\"
	case registry.CURRENT_USER:
		return "HKEY_CURRENT_USER\\"
	case registry.LOCAL_MACHINE:
		return "HKEY_LOCAL_MACHINE\\"
	case registry.USERS:
		return "HKEY_USERS\\"
	case registry.CURRENT_CONFIG:
		return "HKEY_CURRENT_CONFIG\\"
	default:
		return "Unknown"
	}
}

func searchKeyInMap(keywords []string, regData map[string]string) chan Result {
	var wg sync.WaitGroup
	wg.Add(len(keywords))
	// 创建一个互斥锁用于保护计数器
	//var mu sync.Mutex
	//var index int
	// 启动goroutines来搜索关键词
	for _, keyword := range keywords {
		wg.Add(1) // 增加WaitGroup的计数器
		go func(k string) {
			defer wg.Done() // 完成时减少WaitGroup的计数器
			for key, value := range regData {
				if keyContains(key, k) || keyContains(value, k) {
					// 使用互斥锁保护计数器的访问，确保每个结果都分配一个唯一的ID
					//mu.Lock()
					//index++
					////id := index
					//mu.Unlock()

					SearchChan <- Result{Key: k, Path: key, Accuracy: level} // 发送结果到通道
				}
			}
		}(keyword)
	}
	// 启动一个goroutine来关闭通道，当所有搜索完成后
	go func() {
		wg.Wait()         // 等待所有搜索goroutine完成
		close(SearchChan) // 关闭通道
	}()
	return SearchChan
}

//匹配注册表路径的话， //关键字// 说明该文件夹都是软件残留信息

func keyContains(key string, keyword string) bool {
	key = strings.ToLower(key)
	keyword = strings.ToLower(keyword)
	if strings.Contains(key, "\\"+keyword+"\\") {
		level = "high"
		return true
	} else if strings.Contains(key, "\\"+keyword+".") {
		level = "medium"
		return true
	} else if strings.Contains(key, "."+keyword+".") {
		level = "low"
		return true
	}
	// 匹配整个路径or值，开头可以是\或.，结尾可以是\或.
	//pattern := `.*[\\.]` + keyword + `[\\.].*`
	//re, err := regexp.Compile(pattern)
	//if err != nil {
	//	panic(err)
	//}
	//return re.MatchString(key)
	return false
}

func getRootKey(rootKeyName string) (registry.Key, error) {
	switch rootKeyName {
	case "HKEY_CLASSES_ROOT":
		return registry.CLASSES_ROOT, nil
	case "HKEY_CURRENT_USER":
		return registry.CURRENT_USER, nil
	case "HKEY_LOCAL_MACHINE":
		return registry.LOCAL_MACHINE, nil
	case "HKEY_USERS":
		return registry.USERS, nil
	case "HKEY_CURRENT_CONFIG":
		return registry.CURRENT_CONFIG, nil
	default:
		return 0, fmt.Errorf("未知的根键名称: %s", rootKeyName)
	}
}

func DeleteRegistry(input string) {

	firstBackslashIndex := strings.Index(input, "\\")
	if firstBackslashIndex != -1 {
		firstPart := input[:firstBackslashIndex]
		fmt.Println("第一个 \\ 前的字符串:", firstPart)
	} else {
		fmt.Println("找不到第一个 \\")
	}

	// 提取最后一个 \ 后的字符串
	lastBackslashIndex := strings.LastIndex(input, "\\")
	if lastBackslashIndex != -1 {
		lastPart := input[lastBackslashIndex+1:]
		fmt.Println("最后一个 \\ 后的字符串:", lastPart)
	} else {
		fmt.Println("找不到最后一个 \\")
	}

	hive, err := getRootKey(input[:firstBackslashIndex])

	if err != nil {
		//log.Fatal(err)
	}
	for i := lastBackslashIndex - 1; i >= 0; i-- {
		if input[i] == '\\' {
			subKey := input[firstBackslashIndex+1:i] + "\\"
			key, err := registry.OpenKey(hive, subKey, registry.ALL_ACCESS)

			if err != nil {
				fmt.Println(subKey)
				fmt.Printf("找到有效zhi：%s\n", input[i:])
				fmt.Println(err)
			}

			if err == nil {
				fmt.Printf("找到有效路径：%s\n", subKey)
				value := input[i+1:]
				fmt.Println(value)
				err = key.DeleteValue(value)
				if err != nil {
					fmt.Println(value)
					fmt.Println(err)
					//log.Fatal(err)
				}
				break
			}
		}
	}

}
