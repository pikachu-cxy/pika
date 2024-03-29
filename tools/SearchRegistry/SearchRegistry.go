package SearchRegistry

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
	"sync"
)

type Result struct {
	Key  string `json:"key"`
	Path string `json:"path"`
	//Accuracy string `json:"accuracy"`
}

var SearchChan chan Result

var SearchPercentage chan string

func init() {
	SearchChan = make(chan Result, 1)
	SearchPercentage = make(chan string, 1)
}

func SearchRegistry(input string) {
	//start := time.Now()
	//color.Green("注册表残留软件信息检测开始,请稍等片刻--------------------------")
	//这里是windows的读取方式，如果是mac \n
	keywords := strings.Split(input, "\r\n")
	regData := make(map[string]string)
	//读取software注册表进行初始化map
	initRegistryMap("SOFTWARE", regData)

	searchKeyInMap(keywords, regData)

}

func initRegistryMap(path string, regData map[string]string) {
	hives := []registry.Key{
		//SearchRegistry.CLASSES_ROOT,
		registry.CURRENT_USER,
		registry.LOCAL_MACHINE,
		//SearchRegistry.USERS,
		//SearchRegistry.CURRENT_CONFIG,
	}
	// 创建一个 channel 用于从 goroutine 中接收结果
	resultCh := make(chan map[string]string)

	// 创建一个 goroutine 处理每个注册表键
	for _, hive := range hives {
		go func(h registry.Key) {
			key, err := registry.OpenKey(h, path, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
			if err != nil {
				fmt.Println("Error opening key:", err)
				// 将空 map 发送到结果 channel 中以表示错误
				resultCh <- make(map[string]string)
				return
			}
			defer key.Close()
			result := make(map[string]string)
			searchRegistryToMap(h, key, path, result)
			// 发送结果到 channel
			resultCh <- result
		}(hive)
	}

	// 收集所有 goroutine 的结果
	for range hives {
		result := <-resultCh
		// 合并结果到主 map
		for k, v := range result {
			regData[k] = v
		}
	}
	// 关闭 channel
	close(resultCh)
	//
	//for _, hive := range hives {
	//
	//	key, err := SearchRegistry.OpenKey(hive, path, SearchRegistry.ENUMERATE_SUB_KEYS|SearchRegistry.QUERY_VALUE)
	//
	//	if err != nil {
	//		fmt.Println("Error opening key:", err)
	//		continue
	//	}
	//	defer key.Close()
	//	searchRegistryToMap(hive, key, path, regData)
	//}

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

func searchKeyInMap(keywords []string, regData map[string]string) {
	//var resultss []string
	// Create a WaitGroup to track goroutines
	var wg sync.WaitGroup

	// Create a channel to collect results
	//results := make(chan string)

	// Add the number of goroutines to the WaitGroup
	wg.Add(len(keywords))

	// Iterate over the map's key-value pairs
	for _, keyword := range keywords {
		keyword := keyword
		go func(k string) {
			// Decrement the WaitGroup counter when the goroutine completes
			defer wg.Done()
			// Process the key-value pair
			for key, value := range regData {
				// Check if the keyword exists in the key or value
				if keyContains(key, keyword) || keyContains(value, keyword) {
					// If the keyword is matched, print the key-value pair information
					SearchChan <- Result{Key: keyword, Path: key}
					//results <- fmt.Sprintf("Matched keyword %q: Key: %s, Value: %s\n", keyword, key, value)

				}
			}
		}(keyword)
	}

	wg.Wait()
	close(SearchChan)
	// Start a goroutine to close the results channel after all the search goroutines have finished
	//go func() {
	//	wg.Wait()
	//	close(results)
	//}()

	// 在内容末尾加上换行符
	//username := os.Getenv("USERNAME")
	////hostname, _ := os.Hostname()
	////filename := hostname + "-result-" + username + ".txt"
	//timestamp := time.Now().Format("20060102150405")
	//filename := "result-" + username + "." + timestamp
	//file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	//if err != nil {
	//	_ = fmt.Errorf("error opening file: %v", err)
	//}
	//defer func(file *os.File) {
	//	_ = file.Close()
	//}(file)
	//// 创建一个缓冲写入器
	//writer := bufio.NewWriter(file)
	//// Read results from the channel
	//for result := range results {
	//	//color.Red(result)
	//	// Write the result to a file
	//	resultss = append(resultss, result)
	//	_, err = writer.WriteString(result)
	//	if err != nil {
	//		_ = fmt.Errorf("error writing to buffer: %v", err)
	//	}
	//}
	//// 将缓冲区的内容写入文件
	//err = writer.Flush()
	//if err != nil {
	//	_ = fmt.Errorf("error flushing buffer to file: %v", err)
	//}
	//resultss = append(resultss, fmt.Sprintf("全部扫描完成,结果请看文件 %s", filename))

}

//匹配注册表路径的话， //关键字// 说明该文件夹都是软件残留信息

func keyContains(key string, keyword string) bool {
	key = strings.ToLower(key)
	keyword = strings.ToLower(keyword)
	if strings.Contains(key, "\\"+keyword+"\\") || strings.Contains(key, "\\"+keyword+".") || strings.Contains(key, "."+keyword+".") {
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
