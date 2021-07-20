/*
Copyright © 2021 muratgu <mgungora@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kivaorg",
	Short: "Command line interface to your kiva.org account",
	Long: `
kivaorg is a CLI tool for your kiva.org account.
`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./app.env)")

	rootCmd.Version = "0.0.1"
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func Get(query string, values url.Values) (map[string]interface{}, error) {
	url := viper.GetString("KIVAORG_API_URL")
	if url == "" {
		log.Fatal("KIVAORG_API_URL undefined")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, query), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Authorization", auth)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		return data, nil
	} else {
		return nil, errors.New(fmt.Sprintf("HTTP Error %d", res.StatusCode))
	}
}

func IfSetElse(value bool, whenSet string, whenNotSet string) string {
	if value {
		return whenSet
	} else {
		return whenNotSet
	}
}

func Println(data *string, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*data)
	}
}

func JsonEncode(data map[string]interface{}, writer io.Writer) {
	enc := json.NewEncoder(writer)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		log.Fatal(err)
	}
}
