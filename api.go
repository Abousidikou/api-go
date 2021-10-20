package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

type Location struct {
    Id   int
    Name string
}

type Provider struct {
    Id       int
    ISP      string
    ASNumber string
    ASName   string
}

type id_Getted struct {
    Id int
}

type ProviderSample struct {
    Id       int
    downSam  int
    upSam    int
    Provider Provider
}
type ProviderBW struct {
    AvgBW        int
    MinBW        int
    MaxBW        int
    MedianBW     int
    AvgMinRTT    int
    MinMinRTT    int
    MaxMinRTT    int
    MedianMinRTT int
}

type BW struct {
    BW     int
    MinRTT int
}

type medianByDay struct {
    DayStat_Type         string `json:"Type"`
    DayStat_Year         int    `json:"Year"`
    DayStat_Month        int    `json:"Month"`
    DayStat_Day          int    `json:"Day"`
    DayStat_AvgBW        int    `json:"AvgBW"`
    DayStat_MinBW        int    `json:"MinBW"`
    DayStat_MaxBW        int    `json:"MaxBW"`
    DayStat_MedianBW     int    `json:"MedianBW"`
    DayStat_AvgMinRTT    int    `json:"AvgMinRTT"`
    DayStat_MinMinRTT    int    `json:"MinMinRTT"`
    DayStat_MaxMinRTT    int    `json:"MaxMinRTT"`
    DayStat_MedianMinRTT int    `json:"MedianMinRTT"`
}

type paramTCPInfo struct {
    id    int
    day   int
    month int
    year  int
}
type tcpinfos struct {
    avg    int
    min    int
    max    int
    median int
}

type daysliceData struct {
    x string
    y int
}

type daysliceFromTest struct {
    BBRInfo_id  int
    DaySlice_id int
    Year        int
    Month       int
    Day         int
}

type thirdDaySlice struct {
    DaySlice int `json`
    Bw       int `json:"BW"`
}

func getId(id_need, table, colName, val string) []int {
    db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }
    ////fmt.Println("Successful Connected")
    var sql_statement string
    if colName == "" || val == "" {
        sql_statement = "SELECT " + id_need + " FROM " + table
    } else {
        sql_statement = "SELECT " + id_need + " FROM " + table + " where " + colName + "='" + val + "' "
    }
    //fmt.Println(sql_statement)
    res, err := db.Query(sql_statement)
    defer res.Close()
    if err != nil {
        log.Fatal(err)
    }
    ////fmt.Println("Request executed well")
    var row id_Getted
    ids := []int{}
    for res.Next() {

        err := res.Scan(&row.Id)

        if err != nil {
            log.Fatal(err)
        }
        ids = append(ids, row.Id)
    }
    return ids
}

func getl(l []BW) ([]int, []int) {
    //fmt.Println("In Getl")
    var bl []int
    var ll []int
    for _, st := range l {
        bl = append(bl, st.BW)
        ll = append(ll, st.MinRTT)
    }

    return bl, ll
}

func getAvg(l []int) int {
    //fmt.Println("In getAvgMinMaxMedian")
    //fmt.Println("List given:", l)
    var avg int
    var total int
    if len(l) == 0 {
        avg = 0
        return avg
    }
    for _, num := range l {
        total += num
    }
    avg = int(total / len(l))
    //fmt.Println("Avg:", avg)
    return avg
}
func getAvgMinMaxMedian(l []int) []int {
    //fmt.Println("In getAvgMinMaxMedian")
    //fmt.Println("List given:", l)
    if len(l) == 0 {
        to_return := []int{0, 0, 0, 0}
        return to_return
    }
    var min, max, avg, median, total int
    for i, num := range l {
        //fmt.Println(i, num)
        if i == 0 {
            min = num
            max = num
            total += num
            continue
        }
        if num <= min {
            min = num
        }
        if num >= max {
            max = num
        }
        total += num
        //fmt.Println(min, max)
    }
    //fmt.Println("MIn:", min)
    //fmt.Println("Max:", max)
    if len(l) == 1 {
        median = l[0]
    } else if len(l) == 2 {
        median = int((l[0] + l[1]) / 2)
    } else {
        if len(l)%2 == 0 {
            ind := len(l) / 2
            median = int((l[ind] + l[ind+1]) / 2)
        } else {
            ind := int(len(l) / 2)
            median = l[ind+1]
        }
    }
    avg = int(total / len(l))
    //fmt.Println("Avg:", avg)
    //fmt.Println("Med:", median)
    to_return := []int{avg, min, max, median}
    return to_return
}

func BWProcess(bwl []BW) ProviderBW {
    //fmt.Println("In BWProcess")
    //fmt.Println("BWL: ", bwl)
    if len(bwl) == 0 {
        return ProviderBW{
            AvgBW:        0,
            MinBW:        0,
            MaxBW:        0,
            MedianBW:     0,
            AvgMinRTT:    0,
            MinMinRTT:    0,
            MaxMinRTT:    0,
            MedianMinRTT: 0,
        }
    }
    bl, ll := getl(bwl)
    //fmt.Println(bl, ll)
    blpro := getAvgMinMaxMedian(bl)
    //fmt.Println("Blpro:", blpro)
    llpro := getAvgMinMaxMedian(ll)
    //fmt.Println("Llpro:", llpro)
    var proBw ProviderBW
    proBw.AvgBW = blpro[0]
    proBw.MinBW = blpro[1]
    proBw.MaxBW = blpro[2]
    proBw.MedianBW = blpro[3]
    proBw.AvgMinRTT = llpro[0]
    proBw.MinMinRTT = llpro[1]
    proBw.MaxMinRTT = llpro[2]
    proBw.MedianMinRTT = llpro[3]

    return proBw
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func FindInt(slice []int, val int) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

func FindString(slice []string, val string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

func constructDaySlice(l [][]int) []thirdDaySlice {
    var to_return []thirdDaySlice
    var checked []int
    trat := make(map[int][]int)
    //fmt.Println("l:", l)
    for _, id := range l[0] {
        found := FindInt(checked, id)
        //fmt.Println("found:", found)
        if found == true {
            continue
        }
        var slice []int
        for index, val := range l[0] {
            if val == id {
                slice = append(slice, l[1][index])
            }
        }
        //fmt.Println("slice:", slice)
        checked = append(checked, id)
        trat[id] = slice
    }
    //fmt.Println("checked:", checked)
    //fmt.Println("trat:", trat)
    for ind, val := range trat {
        var t thirdDaySlice
        entier := getAvg(val)
        t.DaySlice = ind
        t.Bw = entier
        to_return = append(to_return, t)
    }
    //fmt.Println(to_return)
    return to_return
}

func unicInt(liste []int) []int {
    var to_return []int
    for _, val := range liste {
        found := FindInt(to_return, val)
        if !found {
            to_return = append(to_return, val)
        }
    }
    return to_return
}

func unicString(liste []string) []string {
    var to_return []string
    for _, val := range liste {
        found := FindString(to_return, val)
        if !found {
            to_return = append(to_return, val)
        }
    }
    return to_return
}

func is_a_After_b(a, b string) bool {
    //c := strings.Split(a, "")
    //d := strings.Split(b, "")
    c := []rune(a)
    d := []rune(b)
    fmt.Println(a, b)
    if c[0] > d[0] {
        return true
    } else if c[0] < d[0] {
        return false
    } else {
        if len(c) < len(d) {
            for ind := range c {
                if c[ind] != d[ind] {
                    if c[ind] > d[ind] {
                        return true
                    } else {
                        return false
                    }
                }
            }
        } else {
            for ind := range d {
                if c[ind] != d[ind] {
                    if c[ind] > d[ind] {
                        return true
                    } else {
                        return false
                    }
                }
            }
        }
    }
    return false
}

func rangeString(l []string) []string {
    liste := l
    for ind := range liste {
        if is_a_After_b(liste[ind], liste[ind+1]) {
            tmp := liste[ind]
            liste[ind] = liste[ind+1]
            liste[ind+1] = tmp
            rangeString(liste)
        }
        if ind+1 == len(liste)-1 {
            break
        }
    }
    return liste
}

func main() {

    router := mux.NewRouter()
    router.HandleFunc("/country", func(w http.ResponseWriter, r *http.Request) {
        var row Location
        country := make(map[int]string)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }
        country_ids := getId("Test_Country_id", "Tests", "", "")
        country_ids = unicInt(country_ids)
        fmt.Println("Country ids:", country_ids)

        ////fmt.Println("Successful Connected")

        for _, id := range country_ids {
            res, err := db.Query("SELECT * FROM Country where Country_id=" + strconv.Itoa(id))
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request executed well")
            for res.Next() {

                err := res.Scan(&row.Id, &row.Name)

                if err != nil {
                    log.Fatal(err)
                }
                country[row.Id] = row.Name
            }
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(country)
        return
    })

    router.HandleFunc("/region/{country}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        urlCountry := vars["country"]
        //fmt.Println("Country: " + urlCountry)
        //Get Country Id
        country_id := getId("Country_id", "Country", "Country_Name", urlCountry)
        //fmt.Println(country_id)
        // Get Rgions Ids
        region_ids := getId("Test_Region_id", "Tests", "Test_Country_id", strconv.Itoa(country_id[0]))
        //fmt.Println(region_ids)
        // regions map
        var row Location
        unordered_region := make(map[int]string)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        i := 0
        for _, id := range region_ids {
            //fmt.Println("region_id: ", id)
            res, err := db.Query("SELECT Region_id,Region_Name FROM Region where Region_id=" + strconv.Itoa(id))
            defer res.Close()
            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request executed well")
            for res.Next() {

                err := res.Scan(&row.Id, &row.Name)

                if err != nil {
                    log.Fatal(err)
                }
                unordered_region[i] = row.Name
                i++
            }
        }
        var region []string
        for _, val := range unordered_region {
            found := FindString(region, val)
            if !found {
                region = append(region, val)
            }
        }
        region = unicString(region)
        //fmt.Println(region)
        region = rangeString(region)
        //fmt.Println(region)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(region)
        return
    })
    router.HandleFunc("/city/{region}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        urlregion := vars["region"]
        //fmt.Println("Region: " + urlregion)
        //Get Region Id
        region_id := getId("Region_id", "Region", "Region_Name", urlregion)
        //fmt.Println("Region_id: " + urlregion)
        // Get City ids
        city_ids := getId("Test_City_id", "Tests", "Test_Region_id", strconv.Itoa(region_id[0]))
        var row Location
        cities := make(map[int]string)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        for _, id := range city_ids {
            ////fmt.Println("region_id: ", id)
            res, err := db.Query("SELECT City_id,City_Name FROM City where City_id=" + strconv.Itoa(id))
            defer res.Close()
            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request executed well")
            i := 0
            for res.Next() {

                err := res.Scan(&row.Id, &row.Name)

                if err != nil {
                    log.Fatal(err)
                }
                cities[i] = row.Name
                i++
            }
        }

        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(cities)
        return
    })

    router.HandleFunc("/Sample/{typeofparam}/{param}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        typeOfParam := vars["typeofparam"]
        param := vars["param"]
        sql_statement := ""
        count := make(map[string]int)
        if typeOfParam == "country" {
            country_id := getId("Country_id", "Country", "Country_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_Country_id=" + strconv.Itoa(country_id[0])
        } else if typeOfParam == "city" {
            city_id := getId("City_id", "City", "City_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_City_id=" + strconv.Itoa(city_id[0])
        } else if typeOfParam == "region" {
            region_id := getId("Region_id", "Region", "Region_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_Region_id=" + strconv.Itoa(region_id[0])
        } else if typeOfParam == "downCountry" {
            country_id := getId("Country_id", "Country", "Country_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_Country_id=" + strconv.Itoa(country_id[0]) + " and Test_Type='Download'"
        } else if typeOfParam == "upCountry" {
            country_id := getId("Country_id", "Country", "Country_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_Country_id=" + strconv.Itoa(country_id[0]) + " and Test_Type='Upload'"
        } else if typeOfParam == "downRegion" {
            region_id := getId("Region_id", "Region", "Region_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_Region_id=" + strconv.Itoa(region_id[0]) + " and Test_Type='Download'"
        } else if typeOfParam == "upRegion" {
            region_id := getId("Region_id", "Region", "Region_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_Region_id=" + strconv.Itoa(region_id[0]) + " and Test_Type='Upload'"
        } else if typeOfParam == "downCity" {
            city_id := getId("City_id", "City", "City_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_City_id=" + strconv.Itoa(city_id[0]) + " and Test_Type='Download'"
        } else if typeOfParam == "upCity" {
            city_id := getId("City_id", "City", "City_Name", param)
            sql_statement = "SELECT count(*) FROM Tests where Test_City_id=" + strconv.Itoa(city_id[0]) + " and Test_Type='Upload'"
        }
        //fmt.Println(sql_statement)
        //Connect to database
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }
        ////fmt.Println("Successful Connected")

        res, err := db.Query(sql_statement)
        defer res.Close()

        var c int

        for res.Next() {
            if err := res.Scan(&c); err != nil {
                log.Fatal(err)
            }
            count["Sample"] = c
        }

        ////fmt.Println(count)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(count)
        return
    })
    router.HandleFunc("/percentageByService/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var down, up []int
        count := make(map[string][]int)
        var vars = mux.Vars(r)
        category := vars["type"]
        category_Name := vars["type_id"]
        category_id := 0
        if category == "Country" {
            category_id = getId("Country_id", "Country", "Country_Name", category_Name)[0]
        } else if category == "Region" {
            category_id = getId("Region_id", "Region", "Region_Name", category_Name)[0]
        } else if category == "City" {
            category_id = getId("City_id", "City", "City_Name", category_Name)[0]
        }
        //fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        startDay, _ := strconv.Atoi(startDate[1])
        startMonth, _ := strconv.Atoi(startDate[0])
        startYear, _ := strconv.Atoi(startDate[2])
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        endDay, _ := strconv.Atoi(endDate[1])
        endMonth, _ := strconv.Atoi(endDate[0])
        endYear, _ := strconv.Atoi(endDate[2])
        //fmt.Println(startDate, endDate)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        done := false
        if !done {
            sql_statement := "SELECT Test_Service_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var c int
            for res.Next() {
                if err := res.Scan(&c); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println(c)
                down = append(down, c)
            }
            done = true
        }
        if done {
            sql_statement := "SELECT Test_Service_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var c int

            for res.Next() {
                if err := res.Scan(&c); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println(c)
                up = append(up, c)
            }
            done = false
        }

        ////fmt.Println(down, up)
        count["Download"] = down
        count["Upload"] = up
        ////fmt.Println(count)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(count)
        return
    })

    router.HandleFunc("/medianByDay/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        category := vars["type"]
        category_Name := vars["type_id"]
        category_id := 0
        if category == "Country" {
            category_id = getId("Country_id", "Country", "Country_Name", category_Name)[0]
        } else if category == "Region" {
            category_id = getId("Region_id", "Region", "Region_Name", category_Name)[0]
        } else if category == "City" {
            category_id = getId("City_id", "City", "City_Name", category_Name)[0]
        }
        fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        startDay, _ := strconv.Atoi(startDate[1])
        startMonth, _ := strconv.Atoi(startDate[0])
        startYear, _ := strconv.Atoi(startDate[2])
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        endDay, _ := strconv.Atoi(endDate[1])
        endMonth, _ := strconv.Atoi(endDate[0])
        endYear, _ := strconv.Atoi(endDate[2])
        //fmt.Println(startDate, endDate)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        to_send := make(map[string]medianByDay)
        done := false
        if !done {
            sql_statement := "SELECT DayStat_Type,DayStat_Year,DayStat_Month,DayStat_Day,DayStat_AvgBW,DayStat_MinBW,DayStat_MaxBW,DayStat_MedianBw,DayStat_AvgMinRTT,DayStat_MinMinRTT,DayStat_MaxMinRTT,DayStat_MedianMinRTT from DayStat where DayStat_Type='Download' and  DayStat_Day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and DayStat_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and DayStat_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m medianByDay
            i := 0
            for res.Next() {
                if err := res.Scan(&m.DayStat_Type, &m.DayStat_Year, &m.DayStat_Month, &m.DayStat_Day, &m.DayStat_AvgBW, &m.DayStat_MinBW, &m.DayStat_MaxBW, &m.DayStat_MedianBW, &m.DayStat_AvgMinRTT, &m.DayStat_MinMinRTT, &m.DayStat_MaxMinRTT, &m.DayStat_MedianMinRTT); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println("Type:", m.DayStat_Type)
                indice := "D_" + strconv.Itoa(i)
                to_send[indice] = m
                i++
            }
            done = true
        }

        if done {
            sql_statement := "SELECT DayStat_Type,DayStat_Year,DayStat_Month,DayStat_Day,DayStat_AvgBW,DayStat_MinBW,DayStat_MaxBW,DayStat_MedianBw,DayStat_AvgMinRTT,DayStat_MinMinRTT,DayStat_MaxMinRTT,DayStat_MedianMinRTT from DayStat where DayStat_Type='Upload' and DayStat_Day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and DayStat_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and DayStat_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m medianByDay
            i := 0
            for res.Next() {
                if err := res.Scan(&m.DayStat_Type, &m.DayStat_Year, &m.DayStat_Month, &m.DayStat_Day, &m.DayStat_AvgBW, &m.DayStat_MinBW, &m.DayStat_MaxBW, &m.DayStat_MedianBW, &m.DayStat_AvgMinRTT, &m.DayStat_MinMinRTT, &m.DayStat_MaxMinRTT, &m.DayStat_MedianMinRTT); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println("Type:", m.DayStat_Type)
                indice := "U_" + strconv.Itoa(i)
                to_send[indice] = m
                i++
            }
            done = true
        }

        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        return
    })

    router.HandleFunc("/providerSample/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        category := vars["type"]
        category_Name := vars["type_id"]
        category_id := 0
        if category == "Country" {
            category_id = getId("Country_id", "Country", "Country_Name", category_Name)[0]
        } else if category == "Region" {
            category_id = getId("Region_id", "Region", "Region_Name", category_Name)[0]
        } else if category == "City" {
            category_id = getId("City_id", "City", "City_Name", category_Name)[0]
        }
        fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        startDay, _ := strconv.Atoi(startDate[1])
        startMonth, _ := strconv.Atoi(startDate[0])
        startYear, _ := strconv.Atoi(startDate[2])
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        endDay, _ := strconv.Atoi(endDate[1])
        endMonth, _ := strconv.Atoi(endDate[0])
        endYear, _ := strconv.Atoi(endDate[2])
        //fmt.Println(startDate, endDate)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        prov := make(map[string]Provider)
        done := false
        if !done {
            sql_statement := "SELECT Provider_id,Provider_ISP,Provider_AS_Number,Provider_AS_Name from Provider"
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m Provider
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Id, &m.ISP, &m.ASNumber, &m.ASName); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println(m)
                s := "Prov_" + strconv.Itoa(i)
                prov[s] = m
                i++
            }
            done = true
        }
        fmt.Println("Provider:", prov)
        var d int
        var u int
        if done {
            for _, provider := range prov {
                ////fmt.Println("Provider select : ", provider.Id)
                sql_statement := "SELECT count(*) from Tests where Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println("Request Successful Executed")
                for res.Next() {
                    if err := res.Scan(&d); err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Down Provider: ", d)
                }
                s := "DownSample_" + strconv.Itoa(d)
                //fmt.Println("Test Provider download:", d)
                prov[s] = provider
            }
            done = false
        }
        fmt.Println("Down prov:", prov)
        if !done {
            for _, provider := range prov {
                ////fmt.Println("Provider select : ", provider.Id)
                sql_statement := "SELECT count(*) from Tests where Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println("Request Successful Executed")
                for res.Next() {
                    if err := res.Scan(&u); err != nil {
                        log.Fatal(err)
                    }
                    //fmt.Println("Up Provider: ", u)
                }
                s := "UpSample_" + strconv.Itoa(u)
                //fmt.Println("Test Provider upload:", u)
                prov[s] = provider
            }

            done = true
        }
        fmt.Println("Up and final prov:", prov)
        fmt.Println(prov)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(prov)
        return
    })

    router.HandleFunc("/providerBW/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        category := vars["type"]
        category_Name := vars["type_id"]
        category_id := 0
        if category == "Country" {
            category_id = getId("Country_id", "Country", "Country_Name", category_Name)[0]
        } else if category == "Region" {
            category_id = getId("Region_id", "Region", "Region_Name", category_Name)[0]
        } else if category == "City" {
            category_id = getId("City_id", "City", "City_Name", category_Name)[0]
        }
        //fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        startDay, _ := strconv.Atoi(startDate[1])
        startMonth, _ := strconv.Atoi(startDate[0])
        startYear, _ := strconv.Atoi(startDate[2])
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        endDay, _ := strconv.Atoi(endDate[1])
        endMonth, _ := strconv.Atoi(endDate[0])
        endYear, _ := strconv.Atoi(endDate[2])
        //fmt.Println(startDate, endDate)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        prov := make(map[string]Provider)
        done := false
        if !done {
            sql_statement := "SELECT Provider_id,Provider_ISP,Provider_AS_Number,Provider_AS_Name from Provider"
            ////fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m Provider
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Id, &m.ISP, &m.ASNumber, &m.ASName); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println(m)
                s := "Prov_" + strconv.Itoa(i)
                prov[s] = m
                i++
            }
            done = true
        }

        proBBR := make(map[string][]int)
        if done {
            for _, provider := range prov {
                sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "'  and  Test_Day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
                ////fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()
                var ids []int
                if err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println("Request Successful Executed")
                var c int
                for res.Next() {
                    if err := res.Scan(&c); err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Down Provider: ", c)
                    ids = append(ids, c)
                }
                proBBR[provider.ISP+"_Down"] = ids
            }

            done = false
        }
        //fmt.Println("ProBBR Down:", proBBR)
        if !done {
            for _, provider := range prov {
                sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "'  and  Test_Day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
                ////fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()
                var ids []int
                if err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println("Request Successful Executed")
                var c int
                for res.Next() {
                    if err := res.Scan(&c); err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Down Provider: ", c)
                    ids = append(ids, c)
                }
                //fmt.Println(provider)
                proBBR[provider.ISP+"_Up"] = ids
            }

            done = true
        }

        //fmt.Println("ProBBR All:", proBBR)
        provBW := make(map[string]ProviderBW)
        if done {
            for pro, idl := range proBBR {
                var bwl []BW
                for _, id := range idl {
                    sql_statement := "SELECT AvgBW,AvgMinRTT from BBRInfo where BBRInfo_id='" + strconv.Itoa(id) + "' "
                    ////fmt.Println(sql_statement)
                    res, err := db.Query(sql_statement)
                    defer res.Close()
                    if err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Request Successful Executed")
                    var bw BW
                    for res.Next() {
                        if err := res.Scan(&bw.BW, &bw.MinRTT); err != nil {
                            //fmt.Println("Scanning Error")
                            log.Fatal(err)
                        }
                        bwl = append(bwl, bw)
                    }
                }
                //fmt.Println("Pro: ", pro)
                //fmt.Println("bwl: ", bwl)
                getted := BWProcess(bwl)
                //fmt.Println("Getted:", getted)
                provBW[pro] = getted
                //fmt.Println("ProvBW:", provBW)
            }

        }

        w.Header().Set("Access-Control-Allow-Origin", "*")
        //json.NewEncoder(w).Encode(prov)
        json.NewEncoder(w).Encode(provBW)
        return
    })

    // Return the highest AvgBW of the Day
    router.HandleFunc("/bandByDaySlice/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        category := vars["type"]
        category_Name := vars["type_id"]
        category_id := 0
        if category == "Country" {
            category_id = getId("Country_id", "Country", "Country_Name", category_Name)[0]
        } else if category == "Region" {
            category_id = getId("Region_id", "Region", "Region_Name", category_Name)[0]
        } else if category == "City" {
            category_id = getId("City_id", "City", "City_Name", category_Name)[0]
        }
        //fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        startDay, _ := strconv.Atoi(startDate[1])
        startMonth, _ := strconv.Atoi(startDate[0])
        startYear, _ := strconv.Atoi(startDate[2])
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        endDay, _ := strconv.Atoi(endDate[1])
        endMonth, _ := strconv.Atoi(endDate[0])
        endYear, _ := strconv.Atoi(endDate[2])
        //fmt.Println(startDate, endDate)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        down := make(map[string][][]int)
        up := make(map[string][][]int)
        done := false
        if !done {
            sql_statement := "SELECT Test_BBRInfo_id,Test_DaySlice_id,Test_Year,Test_Month,Test_day from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m daysliceFromTest
            //var slice [][]int
            q := " "
            i := 0
            for res.Next() {
                if err := res.Scan(&m.BBRInfo_id, &m.DaySlice_id, &m.Year, &m.Month, &m.Day); err != nil {
                    log.Fatal(err)
                }
                d := strconv.Itoa(m.Year) + "-" + strconv.Itoa(m.Month) + "-" + strconv.Itoa(m.Day)
                if i == 0 || q != d {
                    var s1, s2 []int
                    var w [][]int
                    s1 = append(s1, m.DaySlice_id)
                    s2 = append(s2, m.BBRInfo_id)
                    w = append(w, s1)
                    w = append(w, s2)
                    down[d] = w
                    //fmt.Println("downDay1:", down)
                    i++
                    q = d
                    continue
                }
                q = d
                down[d][0] = append(down[d][0], m.DaySlice_id)
                down[d][1] = append(down[d][1], m.BBRInfo_id)
                //fmt.Println("downDay2:", down)
                i++
            }
            done = true
        }
        //fmt.Println("downDay3:", down)
        if done {
            sql_statement := "SELECT Test_BBRInfo_id,Test_DaySlice_id,Test_Year,Test_Month,Test_day from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m daysliceFromTest
            q := " "
            i := 0
            for res.Next() {
                if err := res.Scan(&m.BBRInfo_id, &m.DaySlice_id, &m.Year, &m.Month, &m.Day); err != nil {
                    log.Fatal(err)
                }
                d := strconv.Itoa(m.Year) + "-" + strconv.Itoa(m.Month) + "-" + strconv.Itoa(m.Day)
                if i == 0 || q != d {
                    var s1, s2 []int
                    var w [][]int
                    s1 = append(s1, m.DaySlice_id)
                    s2 = append(s2, m.BBRInfo_id)
                    w = append(w, s1)
                    w = append(w, s2)
                    up[d] = w
                    //fmt.Println("upDay1", up)
                    i++
                    q = d
                    continue
                }
                q = d
                up[d][0] = append(up[d][0], m.DaySlice_id)
                up[d][1] = append(up[d][1], m.BBRInfo_id)
                //fmt.Println("upDay2", up)
                i++
            }
            done = false
        }

        //fmt.Println("upDay3", up)
        var days []string
        if !done {
            for i, y := range down {
                days = append(days, i)
                var bw []int
                for _, id := range y[1] {
                    //fmt.Println("bbrinfo id:", id)
                    sql_statement := "SELECT AvgBW from BBRInfo where BBRInfo_id=" + strconv.Itoa(id)
                    //fmt.Println(sql_statement)
                    res, err := db.Query(sql_statement)
                    defer res.Close()

                    if err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Request Successful Executed")
                    var c int
                    for res.Next() {
                        if err := res.Scan(&c); err != nil {
                            log.Fatal(err)
                        }
                        bw = append(bw, c)
                    }
                }
                var d [][]int
                d = append(d, y[0])
                d = append(d, bw)
                down[i] = d
            }
            done = true
        }
        if done {
            for i, y := range up {
                days = append(days, i)
                var bw []int
                for _, id := range y[1] {
                    //fmt.Println("bbrinfo id:", id)
                    sql_statement := "SELECT AvgBW from BBRInfo where BBRInfo_id=" + strconv.Itoa(id)
                    //fmt.Println(sql_statement)
                    res, err := db.Query(sql_statement)
                    defer res.Close()

                    if err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Request Successful Executed")
                    var c int
                    for res.Next() {
                        if err := res.Scan(&c); err != nil {
                            log.Fatal(err)
                        }
                        bw = append(bw, c)
                    }
                }
                var d [][]int
                d = append(d, y[0])
                d = append(d, bw)
                up[i] = d
            }
            done = true
        }
        //fmt.Println(down, up)
        down_to_send := make(map[string][]thirdDaySlice)
        for d, l := range down {
            var key []thirdDaySlice
            key = constructDaySlice(l)
            down_to_send[d] = key
        }
        //fmt.Println("down_to_send:", down_to_send)
        up_to_send := make(map[string][]thirdDaySlice)
        for d, l := range down {
            var key []thirdDaySlice
            key = constructDaySlice(l)
            up_to_send[d] = key
        }
        //fmt.Println("up_to_send:", up_to_send)
        /*jsonrep, err := json.Marshal(to_send)
          if err != nil {
              log.Fatal(err)
          }
          fmt.Println(string(jsonrep))*/
        to_send := make(map[string]map[string][]thirdDaySlice)
        to_send["Download"] = down_to_send
        to_send["Upload"] = up_to_send
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        //w.Write(jsonrep)
        return
    })

    router.HandleFunc("/tcpinfo/param/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        param := vars["param"]
        category := vars["type"]
        category_Name := vars["type_id"]
        category_id := 0
        if category == "Country" {
            category_id = getId("Country_id", "Country", "Country_Name", category_Name)[0]
        } else if category == "Region" {
            category_id = getId("Region_id", "Region", "Region_Name", category_Name)[0]
        } else if category == "City" {
            category_id = getId("City_id", "City", "City_Name", category_Name)[0]
        }
        //fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        startDay, _ := strconv.Atoi(startDate[1])
        startMonth, _ := strconv.Atoi(startDate[0])
        startYear, _ := strconv.Atoi(startDate[2])
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        endDay, _ := strconv.Atoi(endDate[1])
        endMonth, _ := strconv.Atoi(endDate[0])
        endYear, _ := strconv.Atoi(endDate[2])
        //fmt.Println(startDate, endDate)
        db, err := sql.Open("mysql", "root:Emery@123456789@tcp(127.0.0.1:3306)/monitorDB")
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        done := false
        Test := make(map[string]map[string][]int)
        if !done {
            sql_statement := "SELECT Test_TCPInfo_id,Test_Year,Test_Month,Test_Day from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var c paramTCPInfo
            w := make(map[string][]int)
            q := ""
            i := 0
            for res.Next() {
                if err := res.Scan(&c.id, &c.year, &c.month, &c.day); err != nil {
                    log.Fatal(err)
                }
                //fmt.Println(c)
                d := strconv.Itoa(c.day) + "-" + strconv.Itoa(c.month) + "-" + strconv.Itoa(c.year)
                if i == 0 || q != d {
                    _, found := w[d]
                    if found == false {
                        var s1 []int
                        s1 = append(s1, c.id)
                        w[d] = s1
                        q = d
                        i++
                        continue
                    }
                    w[d] = append(w[d], c.id)
                    q = d
                    i++
                    continue
                }
                w[d] = append(w[d], c.id)
                q = d
                i++
            }
            Test["Download"] = w

            done = true
        }
        //fmt.Println("Download:", Test)
        if done {
            sql_statement := "SELECT Test_TCPInfo_id,Test_Year,Test_Month,Test_Day from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_day between " + strconv.Itoa(startDay) + " and " + strconv.Itoa(endDay) + " and Test_Month between " + strconv.Itoa(startMonth) + " and " + strconv.Itoa(endMonth) + " and Test_Year between " + strconv.Itoa(startYear) + " and " + strconv.Itoa(endYear)
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var c paramTCPInfo
            w := make(map[string][]int)
            //var slice []int
            q := ""
            i := 0
            for res.Next() {
                if err := res.Scan(&c.id, &c.year, &c.month, &c.day); err != nil {
                    log.Fatal(err)
                }
                //fmt.Println(c)
                d := strconv.Itoa(c.day) + "-" + strconv.Itoa(c.month) + "-" + strconv.Itoa(c.year)
                if i == 0 || q != d {
                    _, found := w[d]
                    if found == false {
                        var s1 []int
                        s1 = append(s1, c.id)
                        w[d] = s1
                        q = d
                        i++
                        continue
                    }
                    w[d] = append(w[d], c.id)
                    q = d
                    i++
                    continue
                }
                w[d] = append(w[d], c.id)
                q = d
                i++
            }
            Test["Upload"] = w
            done = false
        }
        fmt.Println("Test:", Test)

        if !done {
            for key, tcpinfo_ids := range Test["Download"] {
                // tcpinfo_ids is list [4,5,2,6,5]
                fmt.Println("ids:", tcpinfo_ids)
                var avgg []int
                var minn []int
                var maxx []int
                var mediann []int
                var avgSlice []int
                for _, id := range tcpinfo_ids {
                    sql_statement := "SELECT Avg" + param + ",Min" + param + ",Max" + param + ",Median" + param + " from TCPInfo where TCPInfo_id=" + strconv.Itoa(id)
                    //fmt.Println(sql_statement)
                    res, err := db.Query(sql_statement)
                    defer res.Close()

                    if err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Request Successful Executed")
                    var c tcpinfos

                    for res.Next() {
                        if err := res.Scan(&c.avg, &c.min, &c.max, &c.median); err != nil {
                            log.Fatal(err)
                        }
                        fmt.Println(c)
                        avgg = append(avgg, c.avg)
                        minn = append(minn, c.min)
                        maxx = append(maxx, c.max)
                        mediann = append(mediann, c.median)
                    }
                }
                //fmt.Println(avgg, minn, maxx, mediann)
                avgSlice = append(avgSlice, getAvg(avgg))
                avgSlice = append(avgSlice, getAvg(minn))
                avgSlice = append(avgSlice, getAvg(maxx))
                avgSlice = append(avgSlice, getAvg(mediann))
                Test["Download"][key] = avgSlice
            }
            done = true
        }

        if done {
            for key, tcpinfo_ids := range Test["Upload"] {
                // tcpinfo_ids is list [4,5,2,6,5]
                //fmt.Println(tcpinfo_ids)
                var avgg []int
                var minn []int
                var maxx []int
                var mediann []int
                var avgSlice []int
                for _, id := range tcpinfo_ids {
                    sql_statement := "SELECT Avg" + param + ", Min" + param + ", Max" + param + ", Median" + param + " from TCPInfo where TCPInfo_id=" + strconv.Itoa(id)
                    //fmt.Println(sql_statement)
                    res, err := db.Query(sql_statement)
                    defer res.Close()

                    if err != nil {
                        log.Fatal(err)
                    }
                    ////fmt.Println("Request Successful Executed")
                    var c tcpinfos

                    for res.Next() {
                        if err := res.Scan(&c.avg, &c.min, &c.max, &c.median); err != nil {
                            log.Fatal(err)
                        }
                        fmt.Println(c)
                        avgg = append(avgg, c.avg)
                        minn = append(minn, c.min)
                        maxx = append(maxx, c.max)
                        mediann = append(mediann, c.median)
                    }
                    //fmt.Println(avgg, minn, maxx, mediann)
                }
                avgSlice = append(avgSlice, getAvg(avgg))
                avgSlice = append(avgSlice, getAvg(minn))
                avgSlice = append(avgSlice, getAvg(maxx))
                avgSlice = append(avgSlice, getAvg(mediann))
                Test["Upload"][key] = avgSlice
            }
            done = false
        }
        fmt.Println("Test['Download']:", Test["Download"])
        fmt.Println("Test['Upload']:", Test["Upload"])

        to_send := make(map[string]interface{})
        for t := range Test {

            var days []string
            var avgs []int
            var mins []int
            var maxs []int
            var medians []int
            all := make(map[string][]int)
            for key, value := range Test[t] {
                days = append(days, key)
                avgs = append(avgs, value[0])
                mins = append(mins, value[1])
                maxs = append(maxs, value[2])
                medians = append(medians, value[3])
            }
            all["avg"] = avgs
            all["min"] = mins
            all["max"] = maxs
            all["median"] = medians
            a := t + "_Day"
            to_send[a] = days
            a = t + "_Data"
            to_send[a] = all
        }

        fmt.Println(to_send)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        return
    })

    // create a custom server
    s := &http.Server{
        Addr:    ":4445",
        Handler: router, // use `http.DefaultServeMux`
    }

    cert := "fullchain.pem"
    key := "privkey.pem"
    // run server on port "9000"
    log.Fatal(s.ListenAndServeTLS(cert, key))

    //http.ListenAndServe(":4445", router)
}
