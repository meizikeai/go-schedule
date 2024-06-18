package tool

func GetApiHost(key string) string {
	result := ""

	for k, v := range zookeeperApi {
		if k == key {
			i := GetRandmod(len(v))
			result = v[i]
		}
	}

	return result
}
