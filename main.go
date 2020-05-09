package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if err := mainWithError(); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func mainWithError() error {
	/**
	 * Links needed:
	 * 	NGINX http2 blog post: https://www.nginx.com/blog/nginx-1-13-9-http2-server-push/
	 * 	GoogleChromeLabs http2 manifest generator: https://github.com/GoogleChromeLabs/http2-push-manifest
	 * 	NGINX http2 doc: http://nginx.org/en/docs/http/ngx_http_v2_module.html#http2_push
	 *
	 * Goal:
	 *  Take the push_manifest.json from the http2 manifest generator and convert it into a valid http2_push config for nginx
	 */

	// The http2 manifest generator allows for dynamic output names
	inputManifest := flag.String("input-manifest", "push_manifest.json", "output file of the http2 manifest generator")
	nginxConfig := flag.String("nginx-config", "/etc/nginx/nginx.conf", "config for nginx (needs to include the placeholder)")
	placeholder := flag.String("placeholder", "http2_push_/", "placeholder which should be replaced")

	flag.Parse()

	// Format is {"/file/to/push.js" => {"weight": 1, "type":"script"}}, but nginx only needs the filename so we use a map[string]interface{} and only use the keys
	fileMap := make(map[string]interface{})

	// Read in the push_manifest.json
	data, err := ioutil.ReadFile(*inputManifest)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &fileMap)
	if err != nil {
		return err
	}

	// Read the nginx config
	data, err = ioutil.ReadFile(*nginxConfig)
	if err != nil {
		return err
	}

	strContent := string(data)

	// Detect how much spaces are infront of the placeholder
	prefix := ""
	lines := strings.Split(strContent, "\n")
	for _, line := range lines {
		if strings.Contains(line, "%"+*placeholder+"%") {
			prefix = strings.Split(line, "%")[0]
		}
	}

	// We need to put out the http2_push for each line
	replaceValue := ""
	firstLine := true
	for k := range fileMap {
		if !firstLine {
			replaceValue += prefix
		}

		replaceValue += fmt.Sprintf("http2_push %s;\n", k)
		firstLine = false
	}

	// Copy backup nginx config
	err = ioutil.WriteFile(*nginxConfig+".bak", data, os.FileMode(0666))
	if err != nil {
		return err
	}

	// Replace the placeholder
	strContent = strings.Replace(string(data), "%"+*placeholder+"%", replaceValue, -1)
	err = ioutil.WriteFile(*nginxConfig, []byte(strContent), os.FileMode(0666))
	return err
}
