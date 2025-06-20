package job

import (
	"github.com/patrickmn/go-cache"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"net/http"
	"strings"
)

func CronJobFactory(param CronJobParam) func() {
	return func() {
		resp, err := http.Get(param.Url)
		if err != nil {
			log.Printf("Failed to get %s: %v\n", param.Name, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to get %s: status code %d\n", param.Name, resp.StatusCode)
			return
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read %s: %v\n", param.Name, err)
			return
		}

		if param.ReplaceRMStatic {
			bytes = replaceRMStatic(bytes)
		}

		svc.Cache.Set(param.Name, bytes, cache.DefaultExpiration)

		log.Printf("%s updated\n", strings.ReplaceAll(cases.Title(language.English).String(param.Name), "_", " "))
	}
}

func replaceRMStatic(data []byte) []byte {
	str := string(data)
	str = strings.ReplaceAll(str, "https://rm-static.djicdn.com", "/api/static/rm-static_djicdn_com")
	str = strings.ReplaceAll(str, "https://terra-cn-oss-cdn-public-pro.oss-cn-hangzhou.aliyuncs.com", "/api/static/terra-cn-oss-cdn-public-pro_oss-cn-hangzhou_aliyuncs_com")
	str = strings.ReplaceAll(str, "https://pro-robomasters-hz-n5i3.oss-cn-hangzhou.aliyuncs.com", "/api/static/pro-robomasters-hz-n5i3_oss-cn-hangzhou_aliyuncs_com")
	return []byte(str)
}
