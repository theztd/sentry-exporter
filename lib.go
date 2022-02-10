package main

import "os"

// getnev like in Python
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

/*

func pprint(data {}struct) {
	// // Pretty print
	// prettyResult, prettyErr := json.MarshalIndent(data, "", "  ")
	// if prettyErr != nil {
	// 	panic(prettyErr)
	// }

	// fmt.Println(res.Status)
	// fmt.Printf("%s\n", string(prettyResult))

	// fmt.Println(projects[1].Slug)
}
*/
