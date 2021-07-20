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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"fmt"
)

type stringmap = map[string]interface{}
type arraymap = []interface{}

var listLoansCmd = &cobra.Command{
	Use:   "list-loans",
	Short: "List loans",
	Run: func(cmd *cobra.Command, args []string) {
		data := url.Values{}
		criteria := viper.GetString("KIVAORG_CRITERIA")
		if criteria == "" {
			log.Fatal("KIVAORG_CRITERIA undefined")
		}
		sortBy := viper.GetString("KIVAORG_SORTBY") 
		if sortBy == "" {
			log.Fatal("KIVAORG_SORTBY undefined")
		}
		query := fmt.Sprintf("/loans?limit=3&facets=false&type=lite&sortBy=%s&q=j:{%s}", sortBy, criteria)
		if resp, err := Get(query, data); err != nil {
			log.Fatal(err)
		} else {
			printLoan(resp["entities"].(arraymap)[0].(stringmap))
			printLoan(resp["entities"].(arraymap)[1].(stringmap))			
			JsonEncode(resp, os.Stdout)
		}
	},
}
func printLoan(loan stringmap) {
	properties := loan["properties"].(stringmap)
	geocode := properties["geocode"].(stringmap)
	city := geocode["city"]
	country := geocode["country"].(stringmap)["name"]
	loanAmount := properties["loanAmount"].(stringmap)["amount"]
	id := fmt.Sprintf("%.0f", properties["id"])
	name := properties["name"]
	use := properties["use"]
	fmt.Printf("%s %s from %s, %s needs %s %s\n", id, name, city, country, loanAmount, use)
}
func init() {
	rootCmd.AddCommand(listLoansCmd)
}
