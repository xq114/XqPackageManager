package xpm

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

/*ByteCountIEC convert byte to human readable string*/
func ByteCountIEC(b int64) string {
	const unit int64 = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit && exp <= 5; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

/*GetResponseWithHeader Get the response from the url with proper headers*/
func GetResponseWithHeader(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error connection to %s", url)
	}
	content, err := ioutil.ReadAll(res.Body)
	contents := string(content)
	return contents, nil
}

/*DownloadAndSave Download a file from url and save it to prefix/filename*/
func DownloadAndSave(url string, prefix string, filename string) error {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:81.0) Gecko/20100101 Firefox/81.0",
		// "Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		// "Accept-Encoding":           "gzip, deflate, br",
		// "Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		// "Connection":                "keep-alive",
		// "Upgrade-Insecure-Requests": "1",
		// "Cache-Control":             "max-age=0",
	}

	fmt.Fprintf(os.Stdout, "%s%c%s downloading... ", prefix, os.PathSeparator, path.Base(url))
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError downloading, please check your network connection!\n")
		return err
	}
	defer res.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading response!\n")
		return err
	}

	path := filepath.Join(prefix, filename)
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError creating file!\n")
		return err
	}
	defer file.Close()
	size, err := io.Copy(file, res.Body)
	fmt.Fprintf(os.Stdout, "%s downloaded.\n", ByteCountIEC(size))

	return nil
}
