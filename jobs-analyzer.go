package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "regexp"
  "sort"
  "strconv"
  "strings"
)

var states = map[string]string{
  "Alabama":        "AL",
  "Alaska":         "AK",
  "Arizona":        "AZ",
  "Arkansas":       "AR",
  "California":     "CA",
  "Colorado":       "CO",
  "Connecticut":    "CT",
  "Delaware":       "DE",
  "Florida":        "FL",
  "Georgia":        "GA",
  "Hawaii":         "HI",
  "Idaho":          "ID",
  "Illinois":       "IL",
  "Indiana":        "IN",
  "Iowa":           "IA",
  "Kansas":         "KS",
  "Kentucky":       "KY",
  "Louisiana":      "LA",
  "Maine":          "ME",
  "Maryland":       "MD",
  "Massachusetts":  "MA",
  "Michigan":       "MI",
  "Minnesota":      "MN",
  "Mississippi":    "MS",
  "Missouri":       "MO",
  "Montana":        "MT",
  "Nebraska":       "NE",
  "Nevada":         "NV",
  "New Hampshire":  "NH",
  "New Jersey":     "NJ",
  "New Mexico":     "NM",
  "New York":       "NY",
  "North Carolina": "NC",
  "North Dakota":   "ND",
  "Ohio":           "OH",
  "Oklahoma":       "OK",
  "Oregon":         "OR",
  "Pennsylvania":   "PA",
  "Rhode Island":   "RI",
  "South Carolina": "SC",
  "South Dakota":   "SD",
  "Tennessee":      "TN",
  "Texas":          "TX",
  "Utah":           "UT",
  "Vermont":        "VT",
  "Virginia":       "VA",
  "Washington":     "WA",
  "West Virginia":  "WV",
  "Wisconsin":      "WI",
  "Wyoming":        "WY",
}

type pair struct {
  key   string
  value int64
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].value < p[j].value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func rankByCount(a map[string]int64) pairList {
  pl := make(pairList, len(a))
  i := 0
  for k, v := range a {
    pl[i] = pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

func replaceSpacesWithPlus(role string) string {
  return fmt.Sprint(strings.Replace(role, " ", "+", -1))
}

func replaceSpacesWithMinus(role string) string {
  return fmt.Sprint(strings.Replace(role, " ", "-", -1))
}

func getData(role, regex, site string) map[string]int64 {
  dataMap := map[string]int64{}
  r, _ := regexp.Compile(regex)
  for k, v := range states {
    url := ""
    switch site {
    case "dice":
      url = fmt.Sprint("https://www.dice.com/jobs?q=", replaceSpacesWithPlus(role), "&l=", v, "&searchid=")
    case "indeed":
      url = fmt.Sprint("http://www.indeed.com/jobs?q=", replaceSpacesWithPlus(role), "&l=", v)
    case "monster":
      url = fmt.Sprint("http://www.monster.com/jobs/search/?q=", replaceSpacesWithPlus(role), "&where=", v)
    case "cb":
      url = fmt.Sprint("http://www.careerbuilder.com/jobs-", replaceSpacesWithMinus(role), "-in-", strings.ToLower(v), "?keywords=", replaceSpacesWithPlus(role), "&location=", v)
    }
    resp, err := http.Get(url)
    if err != nil {
      fmt.Printf("Error! %v / %v \n", resp.Status, err)
      os.Exit(1)
    }
    defer resp.Body.Close()

    page, _ := ioutil.ReadAll(resp.Body)

    value := r.FindStringSubmatch(string(page))
    if len(value) >= 2 {
      r, _ := regexp.Compile(",|\\+")
      valueInt, _ := strconv.ParseInt(r.ReplaceAllString(value[1], ""), 10, 64)
      dataMap[fmt.Sprint(k, "/", v)] = valueInt
    } else {
      dataMap[fmt.Sprint(k, "/", v)] = 0
    }
  }
  return dataMap
}

func topTenAmongAllStates(a ...map[string]int64) pairList {
  totalMap := map[string]int64{}

  for i := range a {
    for k, v := range a[i] {
      totalMap[k] += v
    }
  }

  result := rankByCount(totalMap)
  result = append(result[:10])
  return result
}

func main() {
  role := flag.String("p", "", "Job title for search. Example: ./job-analyzer -p 'system engineer'")
  flag.Parse()

  //dice
  diceRegex := "\\d+ - \\d+</span> of <span>(\\d+)</span> positions"
  dice := getData(*role, diceRegex, "dice")

  //indeed
  indeedRegex := "<div style=\"padding-top:9px;\"><div id=\"searchCount\">Jobs \\d+ to \\d+ of (\\d+|\\d+,\\d+)</div>"
  indeed := getData(*role, indeedRegex, "indeed")

  //monster
  monsterRegex := "\"eVar23\":\"(\\d+|\\d+\\+)\""
  monster := getData(*role, monsterRegex, "monster")

  //careerbuilder
  cbRegex := "\\((\\d+) Jobs\\)"
  cb := getData(*role, cbRegex, "cb")

  //analytic functions
  topTen := topTenAmongAllStates(dice, indeed, monster, cb)

  // debug
  // fmt.Println("Dice")
  // for k, v := range dice {
  //   fmt.Println(k, " — ", v)
  // }
  // fmt.Println("-------------------")
  // fmt.Println("Indeed")
  // for k, v := range indeed {
  //   fmt.Println(k, " — ", v)
  // }
  // fmt.Println("-------------------")
  // fmt.Println("Monster")
  // for k, v := range monster {
  //   fmt.Println(k, " — ", v)
  // }
  // fmt.Println("-------------------")
  // fmt.Println("CareerBuilder")
  // for k, v := range cb {
  //   fmt.Println(k, " — ", v)
  // }

  fmt.Printf("Top 10 US states for '%s' role:\n", *role)
  for i := range topTen {
    fmt.Println(topTen[i].key, " — ", topTen[i].value)
  }
}
