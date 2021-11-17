package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

const monthInYear = 12
const dayInMonth = 31
const credential = "root:<password>@tcp(127.0.0.1:3306)/monitorDB"

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
    DayStat_Date         string `json:"Date"`
    DayStat_AvgBW        int    `json:"AvgBW"`
    DayStat_MinBW        int    `json:"MinBW"`
    DayStat_MaxBW        int    `json:"MaxBW"`
    DayStat_MedianBW     int    `json:"MedianBW"`
    DayStat_AvgMinRTT    int    `json:"AvgMinRTT"`
    DayStat_MinMinRTT    int    `json:"MinMinRTT"`
    DayStat_MaxMinRTT    int    `json:"MaxMinRTT"`
    DayStat_MedianMinRTT int    `json:"MedianMinRTT"`
}

type avgMedianByDay struct {
    DayStat_Type         []byte `json:"Type"`
    DayStat_Date         []byte `json:"Date"`
    DayStat_AvgBW        []byte `json:"AvgBW"`
    DayStat_MinBW        []byte `json:"MinBW"`
    DayStat_MaxBW        []byte `json:"MaxBW"`
    DayStat_MedianBW     []byte `json:"MedianBW"`
    DayStat_AvgMinRTT    []byte `json:"AvgMinRTT"`
    DayStat_MinMinRTT    []byte `json:"MinMinRTT"`
    DayStat_MaxMinRTT    []byte `json:"MaxMinRTT"`
    DayStat_MedianMinRTT []byte `json:"MedianMinRTT"`
}
type tcpinfos struct {
    Avg    []byte `json:"Avg"`
    Min    []byte `json:"Min"`
    Max    []byte `json:"Max"`
    Median []byte `json:"Median"`
}
type paramTCPInfo struct {
    Date string
    id   int
}

type daysliceData struct {
    x string
    y int
}

type daysliceFromTest struct {
    Date        string
    BBRInfo_id  int
    DaySlice_id int
}

type thirdDaySlice struct {
    DaySlice int `json:"DaySlice"`
    Bw       int `json:"BW"`
}

//Change

func OneYear() time.Duration {
    t1, _ := time.Parse("2006-01-02", "2021-12-01")
    t2, _ := time.Parse("2006-01-02", "2022-12-01")
    return t2.Sub(t1)
}
func OneMonth() time.Duration {
    t1, _ := time.Parse("2006-01-02", "2021-12-01")
    t2, _ := time.Parse("2006-01-02", "2022-01-01")
    return t2.Sub(t1)
}
func OneDay() time.Duration {
    t1, _ := time.Parse("2006-01-02", "2021-12-01")
    t2, _ := time.Parse("2006-01-02", "2021-12-02")
    return t2.Sub(t1)
}

func TimeDiff(start, end string) (int, int, int) {
    startTime, err := time.Parse("2006-01-02", start)
    endTime, err := time.Parse("2006-01-02", end)
    if err != nil {
        log.Fatal("Time not correct:", err)
    }
    duration := endTime.Sub(startTime)
    year := int(duration.Hours() / OneYear().Hours())
    month := int(duration.Hours() / OneMonth().Hours())
    day := int(duration.Hours() / OneDay().Hours())
    return year, month, day
}

// LastDayOfMonth returns 28-31 - the last day in the month of the time object
// passed in to the function
func LastDayOfMonth(ti string) string {
    t, _ := time.Parse("2006-01-02", ti)
    firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
    lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)

    return lastDay.Format("2006-01-02")
}

func monthInterval(ti string, param int) (string, string) {
    t, _ := time.Parse("2006-01-02", ti)
    y, m, _ := t.Date()
    loc := t.Location()
    firstDay := time.Date(y, m, 1, 0, 0, 0, 0, loc)
    lastDay := time.Date(y, m+time.Month(param), 1, 0, 0, 0, -1, loc)
    return firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")
}

func getMonthListe(st, en string, entierType int) ([]string, []string) {
    startTime, err := time.Parse("2006-01-02", st)
    endTime, err := time.Parse("2006-01-02", en)
    if err != nil {
        log.Fatal("Time not correct:", err)
    }
    var datelistedebut []string
    var datelistefin []string
    last := ""
    for true {
        m := strconv.Itoa(int(startTime.Month()))
        d := strconv.Itoa(startTime.Day())
        if int(startTime.Month()) < 10 {
            m = "0" + m
        }
        if startTime.Day() < 10 {
            d = "0" + d
        }
        w := strconv.Itoa(startTime.Year()) + "-" + m + "-" + d
        if last == "" {
            datelistedebut = append(datelistedebut, w)
        } else {
            if startTime.After(endTime) || startTime.Equal(endTime) {
                m = strconv.Itoa(int(endTime.Month()))
                d = strconv.Itoa(endTime.Day())
                if int(endTime.Month()) < 10 {
                    m = "0" + m
                }
                if endTime.Day() < 10 {
                    d = "0" + d
                }
                w = strconv.Itoa(endTime.Year()) + "-" + m + "-" + d
                datelistefin = append(datelistefin, w)
                if entierType != 0 {
                    break
                }
            } else {
                datelistefin = append(datelistefin, w)
                if startTime.After(endTime) || startTime.Equal(endTime) {
                    break
                }
            }
            if entierType != 0 {
                startTime = startTime.AddDate(0, 0, 1)
                m = strconv.Itoa(int(startTime.Month()))
                d = strconv.Itoa(startTime.Day())
                if int(startTime.Month()) < 10 {
                    m = "0" + m
                }
                if startTime.Day() < 10 {
                    d = "0" + d
                }
                w = strconv.Itoa(startTime.Year()) + "-" + m + "-" + d
            }
            datelistedebut = append(datelistedebut, w)
            if entierType == 0 {
                if startTime.After(endTime) || startTime.Equal(endTime) {
                    break
                }
            }
        }
        last = w
        if entierType == 0 {
            startTime = startTime.AddDate(0, 0, 1)
        } else {
            startTime = startTime.AddDate(0, entierType, 0)
        }

    }
    //fmt.Println(dateliste)
    if entierType == 0 {
        return datelistedebut, datelistedebut
    } else {
        return datelistedebut, datelistefin
    }
}

func rangeDate(dat []string) []string {
    var dateliste []time.Time
    for _, val := range dat {
        q, _ := time.Parse("2006-01-02", val)
        dateliste = append(dateliste, q)
    }
    i := 0
    for true {
        for ind := range dateliste {
            if dateliste[ind].After(dateliste[ind+1]) {
                tmp := dateliste[ind]
                dateliste[ind] = dateliste[ind+1]
                dateliste[ind+1] = tmp
                i = 1
            }
            if ind+1 == len(dateliste)-1 {
                break
            }
        }
        if i == 0 {
            break
        }
        i = 0
    }

    var to_return []string
    for _, val := range dateliste {
        to_return = append(to_return, val.Format("2006-01-01"))
    }

    return to_return
}

func getDateString(st, en string) string {
    startTime, err := time.Parse("2006-01-02", st)
    endTime, err := time.Parse("2006-01-02", en)
    //fmt.Println(st, startTime)
    //fmt.Println(en, endTime)
    if err != nil {
        log.Fatal("Time not correct:", err)
    }
    var y1, y2, to_return string
    firstMonth := startTime.Month().String()[:3]
    //fmt.Println(firstMonth)
    secondMonth := endTime.Month().String()[:3]
    if startTime.Year() == endTime.Year() {
        y1 = strconv.Itoa(startTime.Year())
        if firstMonth == secondMonth {
            if startTime.Day() == endTime.Day() {
                to_return = strconv.Itoa(startTime.Day()) + " " + firstMonth + " " + y1
            } else {
                to_return = strconv.Itoa(startTime.Day()) + "-" + strconv.Itoa(endTime.Day()) + " " + firstMonth + " " + y1
            }
        } else {
            to_return = strconv.Itoa(startTime.Day()) + " " + firstMonth + "-" + strconv.Itoa(endTime.Day()) + " " + secondMonth + " " + y1
        }
    } else {
        y1 = strconv.Itoa(startTime.Year())
        y2 = strconv.Itoa(endTime.Year())
        to_return = strconv.Itoa(startTime.Day()) + " " + firstMonth + " " + y1 + "-" + strconv.Itoa(endTime.Day()) + " " + secondMonth + " " + y2
    }
    return to_return
}

func getFirstDate() string {
    db, err := sql.Open("mysql", credential)
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }
    ////fmt.Println("Successful Connected")
    var sql_statement string
    sql_statement = "SELECT Test_Date from Tests limit 1"
    //fmt.Println(sql_statement)
    res, err := db.Query(sql_statement)
    defer res.Close()
    if err != nil {
        log.Fatal(err)
    }
    ////fmt.Println("Request executed well")
    var row string
    for res.Next() {
        err := res.Scan(&row)
        if err != nil {
            log.Fatal(err)
        }
    }
    return row
}

func is_a_After_bDate(a, b string) bool {
    a_Year, _ := strconv.Atoi(strings.Split(a, "-")[0])
    a_Month, _ := strconv.Atoi(strings.Split(a, "-")[1])
    a_Day, _ := strconv.Atoi(strings.Split(a, "-")[2])
    b_Year, _ := strconv.Atoi(strings.Split(b, "-")[0])
    b_Month, _ := strconv.Atoi(strings.Split(b, "-")[1])
    b_Day, _ := strconv.Atoi(strings.Split(b, "-")[2])

    first := time.Date(a_Year, time.Month(a_Month), a_Day, 0, 0, 0, 0, time.UTC)
    second := time.Date(b_Year, time.Month(b_Month), b_Day, 0, 0, 0, 0, time.UTC)

    return first.After(second)
}
func is_a_equal_bDate(a, b string) bool {
    a_Year, _ := strconv.Atoi(strings.Split(a, "-")[0])
    a_Month, _ := strconv.Atoi(strings.Split(a, "-")[1])
    a_Day, _ := strconv.Atoi(strings.Split(a, "-")[2])
    b_Year, _ := strconv.Atoi(strings.Split(b, "-")[0])
    b_Month, _ := strconv.Atoi(strings.Split(b, "-")[1])
    b_Day, _ := strconv.Atoi(strings.Split(b, "-")[2])

    first := time.Date(a_Year, time.Month(a_Month), a_Day, 0, 0, 0, 0, time.UTC)
    second := time.Date(b_Year, time.Month(b_Month), b_Day, 0, 0, 0, 0, time.UTC)

    return first == second
}
func daySliceToMonth(dateDeb, dateFin string, down, up map[string][]thirdDaySlice) ([]thirdDaySlice, []thirdDaySlice) {
    //dayliste, _ := getMonthListe(dateDeb, dateFin, 0)
    var n thirdDaySlice
    var down_send, up_to_send []thirdDaySlice
    for i := 1; i < 5; i++ {
        var avgDown, avgUp []int
        // Download
        for date, slice := range down {
            if is_a_After_bDate(date, dateDeb) && is_a_After_bDate(dateFin, date) {
                for _, val := range slice {
                    if val.DaySlice == 1 {
                        avgDown = append(avgDown, val.Bw)
                    }
                }
            }
        }
        //avg ready
        n.DaySlice = i
        n.Bw = getAvg(avgDown)
        down_send = append(down_send, n)
        // Upload
        for date, slice := range up {
            if is_a_After_bDate(date, dateDeb) && is_a_After_bDate(dateFin, date) {
                for _, val := range slice {
                    if val.DaySlice == 1 {
                        avgUp = append(avgUp, val.Bw)
                    }
                }
            }
        }
        //avg ready
        n.DaySlice = i
        n.Bw = getAvg(avgUp)
        up_to_send = append(up_to_send, n)
    }
    return down_send, up_to_send
}

////////////////////////////////////////////////////////////////////////////////////Change End

func getId(id_need, table, colName, val string) []int {
    db, err := sql.Open("mysql", credential)
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
    //fmt.Println(a, b)
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
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }
        country_ids := getId("Test_Country_id", "Tests", "", "")
        country_ids = unicInt(country_ids)
        //fmt.Println("Country ids:", country_ids)

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
        db, err := sql.Open("mysql", credential)
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
        db, err := sql.Open("mysql", credential)
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
        db, err := sql.Open("mysql", credential)
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

    //Change
    router.HandleFunc("/percentageByService/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var down, up []int
        count := make(map[string]interface{})
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        /*st := strings.Join(startDate, "-")
          en := strings.Join(endDate, "-")*/
        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]

        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        done := false
        if !done {
            sql_statement := "SELECT Test_Service_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_Date between '" + st + "' and '" + en + "'"
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
            sql_statement := "SELECT Test_Service_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_Date between '" + st + "' and '" + en + "'"
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
        count["len_Down"] = len(down)
        count["Upload"] = up
        count["len_Up"] = len(up)
        ////fmt.Println(count)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(count)
        return
    })

    router.HandleFunc("/percentageByProvider/{provider}/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var down, up []int
        count := make(map[string]interface{})
        var vars = mux.Vars(r)
        provider_name := vars["provider"]
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
        //fmt.Println(category_id, category)
        prov_ID := getId("Provider_id", "Provider", "Provider_AS_Name", provider_name)[0]
        fmt.Println("Prov ID", prov_ID)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        /*st := strings.Join(startDate, "-")
          en := strings.Join(endDate, "-")*/
        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]

        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        done := false
        if !done {
            sql_statement := "SELECT Test_Service_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "' and Test_Date between '" + st + "' and '" + en + "'"
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
            sql_statement := "SELECT Test_Service_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "'  and Test_Date between '" + st + "' and '" + en + "'"
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

        fmt.Println(down, up)
        count["Download"] = down
        count["len_Down"] = len(down)
        count["Upload"] = up
        count["len_Up"] = len(up)
        ////fmt.Println(count)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(count)
        return
    })

    //Change
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
        //fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)
        _, monthDiff, dayDiff := TimeDiff(st, en)
        //fmt.Println(yearDiff, monthDiff, dayDiff)

        // faire la liste des date
        var datelisteDeb, datelisteFin []string
        if dayDiff <= 35 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 0)
        } else if dayDiff > 35 && monthDiff != 0 && monthDiff <= 24 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 1)
        } else if monthDiff > 24 && monthDiff <= 48 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 3)
        } else if monthDiff > 48 && monthDiff <= 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 6)
        } else if monthDiff > 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 12)
        }

        fmt.Println(datelisteDeb, datelisteFin)
        //Base de données
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        to_send := make(map[string]interface{})
        var date []string
        var D_AvgBW []float64
        var D_MinBW []float64
        var D_MaxBW []float64
        var D_MedianBW []float64
        var D_AvgMinRTT []float64
        var D_MinMinRTT []float64
        var D_MaxMinRTT []float64
        var D_MedianMinRTT []float64
        var U_AvgBW []float64
        var U_MinBW []float64
        var U_MaxBW []float64
        var U_MedianBW []float64
        var U_AvgMinRTT []float64
        var U_MinMinRTT []float64
        var U_MaxMinRTT []float64
        var U_MedianMinRTT []float64
        for ind := range datelisteDeb {
            //fmt.Println(datelisteDeb[ind], datelisteFin[ind])
            date = append(date, getDateString(datelisteDeb[ind], datelisteFin[ind]))
            //var d_ids []int
            //var u_ids []int
            done := false
            if !done {
                sql_statement := "SELECT AVG(AvgBW),AVG(MinBw),AVG(MaxBW),AVG(MedianBW),AVG(AvgMinRTT),AVG(MinMinRTT),AVG(MaxMinRTT),AVG(MedianMinRTT) from BBRInfo where BBRInfo_id in (SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                //fmt.Println("Request Successful Executed")
                var m avgMedianByDay
                for res.Next() {
                    if err := res.Scan(&m.DayStat_AvgBW, &m.DayStat_MinBW, &m.DayStat_MaxBW, &m.DayStat_MedianBW, &m.DayStat_AvgMinRTT, &m.DayStat_MinMinRTT, &m.DayStat_MaxMinRTT, &m.DayStat_MedianMinRTT); err != nil {
                        log.Fatal(err)
                    }

                    s, _ := strconv.ParseFloat(string(m.DayStat_AvgBW), 10)
                    D_AvgBW = append(D_AvgBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinBW), 10)
                    D_MinBW = append(D_MinBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxBW), 10)
                    D_MaxBW = append(D_MaxBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianBW), 10)
                    D_MedianBW = append(D_MedianBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_AvgMinRTT), 10)
                    D_AvgMinRTT = append(D_AvgMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinMinRTT), 10)
                    D_MinMinRTT = append(D_MinMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxMinRTT), 10)
                    D_MaxMinRTT = append(D_MaxMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianMinRTT), 10)
                    D_MedianMinRTT = append(D_MedianMinRTT, s)
                }
                done = true
            }
            if done {
                sql_statement := "SELECT AVG(AvgBW),AVG(MinBw),AVG(MaxBW),AVG(MedianBW),AVG(AvgMinRTT),AVG(MinMinRTT),AVG(MaxMinRTT),AVG(MedianMinRTT) from BBRInfo where BBRInfo_id in (SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                var m avgMedianByDay
                for res.Next() {
                    if err := res.Scan(&m.DayStat_AvgBW, &m.DayStat_MinBW, &m.DayStat_MaxBW, &m.DayStat_MedianBW, &m.DayStat_AvgMinRTT, &m.DayStat_MinMinRTT, &m.DayStat_MaxMinRTT, &m.DayStat_MedianMinRTT); err != nil {
                        log.Fatal(err)
                    }
                    //fmt.Println(string(m.DayStat_AvgBW))
                    s, _ := strconv.ParseFloat(string(m.DayStat_AvgBW), 10)
                    U_AvgBW = append(U_AvgBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinBW), 10)
                    U_MinBW = append(U_MinBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxBW), 10)
                    U_MaxBW = append(U_MaxBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianBW), 10)
                    U_MedianBW = append(U_MedianBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_AvgMinRTT), 10)
                    U_AvgMinRTT = append(U_AvgMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinMinRTT), 10)
                    U_MinMinRTT = append(U_MinMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxMinRTT), 10)
                    U_MaxMinRTT = append(U_MaxMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianMinRTT), 10)
                    U_MedianMinRTT = append(U_MedianMinRTT)
                }
                done = false
            }
        }
        to_send["D_Date"] = date
        to_send["D_AvgBW"] = D_AvgBW
        to_send["D_MinBW"] = D_MinBW
        to_send["D_MaxBW"] = D_MaxBW
        to_send["D_MedianBW"] = D_MedianBW
        to_send["D_AvgMinRTT"] = D_AvgMinRTT
        to_send["D_MinMinRTT"] = D_MinMinRTT
        to_send["D_MaxMinRTT"] = D_MaxMinRTT
        to_send["D_MedianMinRTT"] = D_MedianMinRTT
        to_send["U_AvgBW"] = U_AvgBW
        to_send["U_MinBW"] = U_MinBW
        to_send["U_MaxBW"] = U_MaxBW
        to_send["U_MedianBW"] = U_MedianBW
        to_send["U_AvgMinRTT"] = U_AvgMinRTT
        to_send["U_MinMinRTT"] = U_MinMinRTT
        to_send["U_MaxMinRTT"] = U_MaxMinRTT
        to_send["U_MedianMinRTT"] = U_MedianMinRTT

        //fmt.Println("to_send:", to_send)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        return
    })

    //Change
    router.HandleFunc("/medianByProvider/{provider}/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        provider_name := vars["provider"]
        prov_ID := getId("Provider_id", "Provider", "Provider_AS_Name", provider_name)[0]
        fmt.Println("Prov ID", prov_ID)
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)
        _, monthDiff, dayDiff := TimeDiff(st, en)
        //fmt.Println(yearDiff, monthDiff, dayDiff)

        // faire la liste des date
        var datelisteDeb, datelisteFin []string
        if dayDiff <= 35 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 0)
        } else if dayDiff > 35 && monthDiff != 0 && monthDiff <= 24 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 1)
        } else if monthDiff > 24 && monthDiff <= 48 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 3)
        } else if monthDiff > 48 && monthDiff <= 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 6)
        } else if monthDiff > 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 12)
        }

        //fmt.Println(datelisteDeb, datelisteFin)
        //Base de données
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")

        to_send := make(map[string]interface{})
        var date []string
        var D_AvgBW []float64
        var D_MinBW []float64
        var D_MaxBW []float64
        var D_MedianBW []float64
        var D_AvgMinRTT []float64
        var D_MinMinRTT []float64
        var D_MaxMinRTT []float64
        var D_MedianMinRTT []float64
        var U_AvgBW []float64
        var U_MinBW []float64
        var U_MaxBW []float64
        var U_MedianBW []float64
        var U_AvgMinRTT []float64
        var U_MinMinRTT []float64
        var U_MaxMinRTT []float64
        var U_MedianMinRTT []float64
        for ind := range datelisteDeb {
            //fmt.Println(datelisteDeb[ind], datelisteFin[ind])
            date = append(date, getDateString(datelisteDeb[ind], datelisteFin[ind]))
            //var d_ids []int
            //var u_ids []int
            done := false
            if !done {
                sql_statement := "SELECT AVG(AvgBW),AVG(MinBw),AVG(MaxBW),AVG(MedianBW),AVG(AvgMinRTT),AVG(MinMinRTT),AVG(MaxMinRTT),AVG(MedianMinRTT) from BBRInfo where BBRInfo_id in (SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                //fmt.Println("Request Successful Executed")
                var m avgMedianByDay
                for res.Next() {
                    if err := res.Scan(&m.DayStat_AvgBW, &m.DayStat_MinBW, &m.DayStat_MaxBW, &m.DayStat_MedianBW, &m.DayStat_AvgMinRTT, &m.DayStat_MinMinRTT, &m.DayStat_MaxMinRTT, &m.DayStat_MedianMinRTT); err != nil {
                        log.Fatal(err)
                    }

                    s, _ := strconv.ParseFloat(string(m.DayStat_AvgBW), 10)
                    D_AvgBW = append(D_AvgBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinBW), 10)
                    D_MinBW = append(D_MinBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxBW), 10)
                    D_MaxBW = append(D_MaxBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianBW), 10)
                    D_MedianBW = append(D_MedianBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_AvgMinRTT), 10)
                    D_AvgMinRTT = append(D_AvgMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinMinRTT), 10)
                    D_MinMinRTT = append(D_MinMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxMinRTT), 10)
                    D_MaxMinRTT = append(D_MaxMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianMinRTT), 10)
                    D_MedianMinRTT = append(D_MedianMinRTT, s)
                }
                done = true
            }
            if done {
                sql_statement := "SELECT AVG(AvgBW),AVG(MinBw),AVG(MaxBW),AVG(MedianBW),AVG(AvgMinRTT),AVG(MinMinRTT),AVG(MaxMinRTT),AVG(MedianMinRTT) from BBRInfo where BBRInfo_id in (SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                var m avgMedianByDay
                for res.Next() {
                    if err := res.Scan(&m.DayStat_AvgBW, &m.DayStat_MinBW, &m.DayStat_MaxBW, &m.DayStat_MedianBW, &m.DayStat_AvgMinRTT, &m.DayStat_MinMinRTT, &m.DayStat_MaxMinRTT, &m.DayStat_MedianMinRTT); err != nil {
                        log.Fatal(err)
                    }
                    //fmt.Println(string(m.DayStat_AvgBW))
                    s, _ := strconv.ParseFloat(string(m.DayStat_AvgBW), 10)
                    U_AvgBW = append(U_AvgBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinBW), 10)
                    U_MinBW = append(U_MinBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxBW), 10)
                    U_MaxBW = append(U_MaxBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianBW), 10)
                    U_MedianBW = append(U_MedianBW, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_AvgMinRTT), 10)
                    U_AvgMinRTT = append(U_AvgMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MinMinRTT), 10)
                    U_MinMinRTT = append(U_MinMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MaxMinRTT), 10)
                    U_MaxMinRTT = append(U_MaxMinRTT, s)
                    s, _ = strconv.ParseFloat(string(m.DayStat_MedianMinRTT), 10)
                    U_MedianMinRTT = append(U_MedianMinRTT)
                }
                done = false
            }
        }
        to_send["D_Date"] = date
        to_send["D_AvgBW"] = D_AvgBW
        to_send["D_MinBW"] = D_MinBW
        to_send["D_MaxBW"] = D_MaxBW
        to_send["D_MedianBW"] = D_MedianBW
        to_send["D_AvgMinRTT"] = D_AvgMinRTT
        to_send["D_MinMinRTT"] = D_MinMinRTT
        to_send["D_MaxMinRTT"] = D_MaxMinRTT
        to_send["D_MedianMinRTT"] = D_MedianMinRTT
        to_send["U_AvgBW"] = U_AvgBW
        to_send["U_MinBW"] = U_MinBW
        to_send["U_MaxBW"] = U_MaxBW
        to_send["U_MedianBW"] = U_MedianBW
        to_send["U_AvgMinRTT"] = U_AvgMinRTT
        to_send["U_MinMinRTT"] = U_MinMinRTT
        to_send["U_MaxMinRTT"] = U_MaxMinRTT
        to_send["U_MedianMinRTT"] = U_MedianMinRTT

        //fmt.Println("to_send:", to_send)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)
        _, monthDiff, dayDiff := TimeDiff(st, en)
        //fmt.Println(yearDiff, monthDiff, dayDiff)

        //fmt.Println(st, en)
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        //fmt.Println("Successful Connected")
        down := make(map[string][][]int)
        up := make(map[string][][]int)
        done := false
        if !done {
            sql_statement := "SELECT Test_Date,Test_BBRInfo_id,Test_DaySlice_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_Date between '" + st + "' and '" + en + "'"
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m daysliceFromTest
            var q []string
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Date, &m.BBRInfo_id, &m.DaySlice_id); err != nil {
                    log.Fatal(err)
                }
                d := m.Date
                found := FindString(q, d)
                if i == 0 || !found {
                    var s1, s2 []int
                    var w [][]int
                    s1 = append(s1, m.DaySlice_id)
                    s2 = append(s2, m.BBRInfo_id)
                    w = append(w, s1)
                    w = append(w, s2)
                    down[d] = w
                    //fmt.Println("downDay1:", down)
                    i++
                    q = append(q, d)
                    continue
                }
                q = append(q, d)
                down[d][0] = append(down[d][0], m.DaySlice_id)
                down[d][1] = append(down[d][1], m.BBRInfo_id)
                //fmt.Println("downDay2:", down)
                i++
            }
            done = true
        }
        //fmt.Println("downDay3:", down)
        if done {
            sql_statement := "SELECT Test_Date,Test_BBRInfo_id,Test_DaySlice_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_Date between '" + st + "' and '" + en + "'"
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m daysliceFromTest
            var q []string
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Date, &m.BBRInfo_id, &m.DaySlice_id); err != nil {
                    log.Fatal(err)
                }
                d := m.Date
                found := FindString(q, d)
                if i == 0 || !found {
                    var s1, s2 []int
                    var w [][]int
                    s1 = append(s1, m.DaySlice_id)
                    s2 = append(s2, m.BBRInfo_id)
                    w = append(w, s1)
                    w = append(w, s2)
                    up[d] = w
                    //fmt.Println("upDay1", up)
                    i++
                    q = append(q, d)
                    continue
                }
                q = append(q, d)
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
        //fmt.Println("Down with Avg:", down)
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
        //fmt.Println("Up with Avg", up)
        down_to_send := make(map[string][]thirdDaySlice)
        for d, l := range down {
            var key []thirdDaySlice
            key = constructDaySlice(l)
            down_to_send[d] = key
        }
        //fmt.Println("down_to_send:", down_to_send)
        up_to_send := make(map[string][]thirdDaySlice)
        for d, l := range up {
            var key []thirdDaySlice
            key = constructDaySlice(l)
            up_to_send[d] = key
        }
        //fmt.Println("up_to_send:", up_to_send)
        to_send := make(map[string]map[string][]thirdDaySlice)
        to_send["Download"] = down_to_send
        to_send["Upload"] = up_to_send

        // According to month
        if dayDiff > 30 {
            // faire la liste des date
            var datelisteDeb, datelisteFin []string
            if monthDiff <= 24 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 1)
            } else if monthDiff > 24 && monthDiff <= 48 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 3)
            } else if monthDiff > 48 && monthDiff <= 60 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 6)
            } else if monthDiff > 60 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 12)
            }
            //fmt.Println(datelisteDeb, datelisteFin)
            tp1 := make(map[string][]thirdDaySlice)
            tp2 := make(map[string][]thirdDaySlice)
            for ind := range datelisteDeb {
                date := getDateString(datelisteDeb[ind], datelisteFin[ind])
                d_tmp, u_tmp := daySliceToMonth(datelisteDeb[ind], datelisteFin[ind], down_to_send, up_to_send)
                tp1[date] = d_tmp
                tp2[date] = u_tmp
            }
            to_send["Download"] = tp1
            to_send["Upload"] = tp2
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        //w.Write(jsonrep)
        return
    })

    router.HandleFunc("/bandByDaySliceProvider/{provider}/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        provider_name := vars["provider"]
        prov_ID := getId("Provider_id", "Provider", "Provider_AS_Name", provider_name)[0]
        fmt.Println("Prov ID", prov_ID)
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)
        _, monthDiff, dayDiff := TimeDiff(st, en)
        //fmt.Println(yearDiff, monthDiff, dayDiff)

        //fmt.Println(st, en)
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        //fmt.Println("Successful Connected")
        down := make(map[string][][]int)
        up := make(map[string][][]int)
        done := false
        if !done {
            sql_statement := "SELECT Test_Date,Test_BBRInfo_id,Test_DaySlice_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "'  and Test_Type='Download' and Test_Date between '" + st + "' and '" + en + "'"
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m daysliceFromTest
            var q []string
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Date, &m.BBRInfo_id, &m.DaySlice_id); err != nil {
                    log.Fatal(err)
                }
                d := m.Date
                found := FindString(q, d)
                if i == 0 || !found {
                    var s1, s2 []int
                    var w [][]int
                    s1 = append(s1, m.DaySlice_id)
                    s2 = append(s2, m.BBRInfo_id)
                    w = append(w, s1)
                    w = append(w, s2)
                    down[d] = w
                    //fmt.Println("downDay1:", down)
                    i++
                    q = append(q, d)
                    continue
                }
                q = append(q, d)
                down[d][0] = append(down[d][0], m.DaySlice_id)
                down[d][1] = append(down[d][1], m.BBRInfo_id)
                //fmt.Println("downDay2:", down)
                i++
            }
            done = true
        }
        //fmt.Println("downDay3:", down)
        if done {
            sql_statement := "SELECT Test_Date,Test_BBRInfo_id,Test_DaySlice_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Upload' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "' and Test_Date between '" + st + "' and '" + en + "'"
            //fmt.Println(sql_statement)
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var m daysliceFromTest
            var q []string
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Date, &m.BBRInfo_id, &m.DaySlice_id); err != nil {
                    log.Fatal(err)
                }
                d := m.Date
                found := FindString(q, d)
                if i == 0 || !found {
                    var s1, s2 []int
                    var w [][]int
                    s1 = append(s1, m.DaySlice_id)
                    s2 = append(s2, m.BBRInfo_id)
                    w = append(w, s1)
                    w = append(w, s2)
                    up[d] = w
                    //fmt.Println("upDay1", up)
                    i++
                    q = append(q, d)
                    continue
                }
                q = append(q, d)
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
        //fmt.Println("Down with Avg:", down)
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
        //fmt.Println("Up with Avg", up)
        down_to_send := make(map[string][]thirdDaySlice)
        for d, l := range down {
            var key []thirdDaySlice
            key = constructDaySlice(l)
            down_to_send[d] = key
        }
        //fmt.Println("down_to_send:", down_to_send)
        up_to_send := make(map[string][]thirdDaySlice)
        for d, l := range up {
            var key []thirdDaySlice
            key = constructDaySlice(l)
            up_to_send[d] = key
        }
        //fmt.Println("up_to_send:", up_to_send)
        to_send := make(map[string]map[string][]thirdDaySlice)
        to_send["Download"] = down_to_send
        to_send["Upload"] = up_to_send

        // According to month
        if dayDiff > 30 {
            // faire la liste des date
            var datelisteDeb, datelisteFin []string
            if monthDiff <= 24 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 1)
            } else if monthDiff > 24 && monthDiff <= 48 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 3)
            } else if monthDiff > 48 && monthDiff <= 60 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 6)
            } else if monthDiff > 60 {
                datelisteDeb, datelisteFin = getMonthListe(st, en, 12)
            }
            //fmt.Println(datelisteDeb, datelisteFin)
            tp1 := make(map[string][]thirdDaySlice)
            tp2 := make(map[string][]thirdDaySlice)
            for ind := range datelisteDeb {
                date := getDateString(datelisteDeb[ind], datelisteFin[ind])
                d_tmp, u_tmp := daySliceToMonth(datelisteDeb[ind], datelisteFin[ind], down_to_send, up_to_send)
                tp1[date] = d_tmp
                tp2[date] = u_tmp
            }
            to_send["Download"] = tp1
            to_send["Upload"] = tp2
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        //w.Write(jsonrep)
        return
    })

    router.HandleFunc("/tcpinfo/{param}/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)

        _, monthDiff, dayDiff := TimeDiff(st, en)

        // faire la liste des date
        var datelisteDeb, datelisteFin []string
        if dayDiff <= 35 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 0)
        } else if dayDiff > 35 && monthDiff != 0 && monthDiff <= 24 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 1)
        } else if monthDiff > 24 && monthDiff <= 48 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 3)
        } else if monthDiff > 48 && monthDiff <= 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 6)
        } else if monthDiff > 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 12)
        }

        //fmt.Println(datelisteDeb, datelisteFin)
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        var date []string
        var D_Avg []float64
        var D_Min []float64
        var D_Max []float64
        var D_Median []float64
        var U_Avg []float64
        var U_Min []float64
        var U_Max []float64
        var U_Median []float64

        for ind := range datelisteDeb {
            //fmt.Println(datelisteDeb[ind], datelisteFin[ind])
            date = append(date, getDateString(datelisteDeb[ind], datelisteFin[ind]))
            //var d_ids []int
            //var u_ids []int
            done := false
            if !done {
                sql_statement := "SELECT AVG(Avg" + param + "),AVG(Min" + param + "),AVG(Max" + param + "),AVG(Median" + param + ") from TCPInfo where TCPInfo_id in (SELECT Test_TCPInfo_id from Tests where Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                //fmt.Println("Request Successful Executed")
                var m tcpinfos
                for res.Next() {
                    if err := res.Scan(&m.Avg, &m.Min, &m.Max, &m.Median); err != nil {
                        log.Fatal(err)
                    }

                    s, _ := strconv.ParseFloat(string(m.Avg), 10)
                    D_Avg = append(D_Avg, s)
                    s, _ = strconv.ParseFloat(string(m.Min), 10)
                    D_Min = append(D_Min, s)
                    s, _ = strconv.ParseFloat(string(m.Max), 10)
                    D_Max = append(D_Max, s)
                    s, _ = strconv.ParseFloat(string(m.Median), 10)
                    D_Median = append(D_Median, s)
                }
                done = true
            }
            if done {
                sql_statement := "SELECT AVG(Avg" + param + "),AVG(Min" + param + "),AVG(Max" + param + "),AVG(Median" + param + ") from TCPInfo where TCPInfo_id in (SELECT Test_TCPInfo_id from Tests where Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                //fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                var m tcpinfos
                for res.Next() {
                    if err := res.Scan(&m.Avg, &m.Min, &m.Max, &m.Median); err != nil {
                        log.Fatal(err)
                    }

                    s, _ := strconv.ParseFloat(string(m.Avg), 10)
                    U_Avg = append(U_Avg, s)
                    s, _ = strconv.ParseFloat(string(m.Min), 10)
                    U_Min = append(U_Min, s)
                    s, _ = strconv.ParseFloat(string(m.Max), 10)
                    U_Max = append(U_Max, s)
                    s, _ = strconv.ParseFloat(string(m.Median), 10)
                    U_Median = append(U_Median, s)
                }
                done = false
            }
        }

        to_send := make(map[string]interface{})
        to_send["D_Date"] = date
        to_send["D_Avg"] = D_Avg
        to_send["D_Min"] = D_Min
        to_send["D_Max"] = D_Max
        to_send["D_Median"] = D_Median
        to_send["U_Avg"] = U_Avg
        to_send["U_Min"] = U_Min
        to_send["U_Max"] = U_Max
        to_send["U_Median"] = U_Median
        //fmt.Println(to_send)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(to_send)
        return
    })

    router.HandleFunc("/tcpinfoProvider/{provider}/{param}/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
        var vars = mux.Vars(r)
        provider_name := vars["provider"]
        prov_ID := getId("Provider_id", "Provider", "Provider_AS_Name", provider_name)[0]
        fmt.Println("Prov ID", prov_ID)
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)

        _, monthDiff, dayDiff := TimeDiff(st, en)

        // faire la liste des date
        var datelisteDeb, datelisteFin []string
        if dayDiff <= 35 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 0)
        } else if dayDiff > 35 && monthDiff != 0 && monthDiff <= 24 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 1)
        } else if monthDiff > 24 && monthDiff <= 48 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 3)
        } else if monthDiff > 48 && monthDiff <= 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 6)
        } else if monthDiff > 60 {
            datelisteDeb, datelisteFin = getMonthListe(st, en, 12)
        }

        //fmt.Println(datelisteDeb, datelisteFin)
        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        var date []string
        var D_Avg []float64
        var D_Min []float64
        var D_Max []float64
        var D_Median []float64
        var U_Avg []float64
        var U_Min []float64
        var U_Max []float64
        var U_Median []float64

        for ind := range datelisteDeb {
            //fmt.Println(datelisteDeb[ind], datelisteFin[ind])
            date = append(date, getDateString(datelisteDeb[ind], datelisteFin[ind]))
            //var d_ids []int
            //var u_ids []int
            done := false
            if !done {
                sql_statement := "SELECT AVG(Avg" + param + "),AVG(Min" + param + "),AVG(Max" + param + "),AVG(Median" + param + ") from TCPInfo where TCPInfo_id in (SELECT Test_TCPInfo_id from Tests where Test_Type='Download' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                //fmt.Println("Request Successful Executed")
                var m tcpinfos
                for res.Next() {
                    if err := res.Scan(&m.Avg, &m.Min, &m.Max, &m.Median); err != nil {
                        log.Fatal(err)
                    }

                    s, _ := strconv.ParseFloat(string(m.Avg), 10)
                    D_Avg = append(D_Avg, s)
                    s, _ = strconv.ParseFloat(string(m.Min), 10)
                    D_Min = append(D_Min, s)
                    s, _ = strconv.ParseFloat(string(m.Max), 10)
                    D_Max = append(D_Max, s)
                    s, _ = strconv.ParseFloat(string(m.Median), 10)
                    D_Median = append(D_Median, s)
                }
                done = true
            }
            if done {
                sql_statement := "SELECT AVG(Avg" + param + "),AVG(Min" + param + "),AVG(Max" + param + "),AVG(Median" + param + ") from TCPInfo where TCPInfo_id in (SELECT Test_TCPInfo_id from Tests where Test_Type='Download' and Test_Provider_id='" + strconv.Itoa(prov_ID) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "')"
                //sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + datelisteDeb[ind] + "' and '" + datelisteFin[ind] + "'"
                fmt.Println(sql_statement)
                res, err := db.Query(sql_statement)
                defer res.Close()

                if err != nil {
                    log.Fatal(err)
                }
                var m tcpinfos
                for res.Next() {
                    if err := res.Scan(&m.Avg, &m.Min, &m.Max, &m.Median); err != nil {
                        log.Fatal(err)
                    }

                    s, _ := strconv.ParseFloat(string(m.Avg), 10)
                    U_Avg = append(U_Avg, s)
                    s, _ = strconv.ParseFloat(string(m.Min), 10)
                    U_Min = append(U_Min, s)
                    s, _ = strconv.ParseFloat(string(m.Max), 10)
                    U_Max = append(U_Max, s)
                    s, _ = strconv.ParseFloat(string(m.Median), 10)
                    U_Median = append(U_Median, s)
                }
                done = false
            }
        }

        to_send := make(map[string]interface{})
        to_send["D_Date"] = date
        to_send["D_Avg"] = D_Avg
        to_send["D_Min"] = D_Min
        to_send["D_Max"] = D_Max
        to_send["D_Median"] = D_Median
        to_send["U_Avg"] = U_Avg
        to_send["U_Min"] = U_Min
        to_send["U_Max"] = U_Max
        to_send["U_Median"] = U_Median
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
        //fmt.Println(category, category_id)
        //fmt.Println(category, category_id)
        startDate := strings.Split(strings.Split(vars["dayRange"], "-")[0], ",")
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)
        db, err := sql.Open("mysql", credential)
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
            //fmt.Println("Request Successful Executed")
            var m Provider
            i := 0
            for res.Next() {
                if err := res.Scan(&m.Id, &m.ISP, &m.ASNumber, &m.ASName); err != nil {
                    log.Fatal(err)
                }
                //fmt.Println("m:", m)
                s := "Prov_" + strconv.Itoa(i)
                prov[s] = m
                i++
            }
            done = true
        }
        //fmt.Println("Provider:", prov)
        var d int
        var u int
        if done {
            for _, provider := range prov {
                ////fmt.Println("Provider select : ", provider.Id)
                sql_statement := "SELECT count(*) from Tests where Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_Type='Download' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + st + "' and '" + en + "'"
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
                    if d == 0 {
                        continue
                    }
                    ////fmt.Println("Down Provider: ", d)
                }
                s := "DownSample_" + strconv.Itoa(d)
                //fmt.Println("Test Provider download:", d)
                prov[s] = provider
            }
            done = false
        }
        //fmt.Println("Down prov:", prov)
        if !done {
            for _, provider := range prov {
                ////fmt.Println("Provider select : ", provider.Id)
                sql_statement := "SELECT count(*) from Tests where Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_Type='Upload' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Date between '" + st + "' and '" + en + "'"
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
        //fmt.Println("Up and final prov:", prov)
        //fmt.Println(prov)
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]
        //fmt.Println(st, en)
        db, err := sql.Open("mysql", credential)
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
                sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Download' and Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "'  and  Test_Date between '" + st + "' and '" + en + "'"
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
                    //fmt.Println("Down Provider: ", c)
                    ids = append(ids, c)
                }
                if len(ids) != 0 {
                    proBBR[provider.ASName+"_Down"] = ids
                }
            }

            done = false
        }
        //fmt.Println("ProBBR Down:", proBBR)
        if !done {
            for _, provider := range prov {
                sql_statement := "SELECT Test_BBRInfo_id from Tests where Test_Type='Upload' and Test_Provider_id='" + strconv.Itoa(provider.Id) + "' and Test_" + category + "_id='" + strconv.Itoa(category_id) + "'  and  Test_Date between '" + st + "' and '" + en + "'"
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
                    //fmt.Println("Up Provider: ", c)
                    ids = append(ids, c)
                }
                //fmt.Println(provider)
                if len(ids) == 0 {
                    proBBR[provider.ASName+"_Up"] = ids
                }

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

    router.HandleFunc("/providersListe/{type}/{type_id}/{dayRange}", func(w http.ResponseWriter, r *http.Request) {
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
        endDate := strings.Split(strings.Split(vars["dayRange"], "-")[1], ",")
        //fmt.Println(startDate, endDate)

        /*st := strings.Join(startDate, "-")
          en := strings.Join(endDate, "-")*/
        st := startDate[2] + "-" + startDate[0] + "-" + startDate[1]
        en := endDate[2] + "-" + endDate[0] + "-" + endDate[1]

        db, err := sql.Open("mysql", credential)
        defer db.Close()

        if err != nil {
            log.Fatal(err)
        }

        ////fmt.Println("Successful Connected")
        var prov_liste []string
        done := false
        if !done {
            sql_statement := "select Provider_AS_Name from Provider where Provider_id in (select Test_Provider_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_Date between '" + st + "' and '" + en + "')"
            //fmt.Println(sql_statement) "select Provider_AS_Name,Provider_id from Provider where Provider_id in (select Test_Provider_id from Tests where Test_" + category + "_id='" + strconv.Itoa(category_id) + "' and Test_Type='Download' and Test_Date between '"+st+"' and '"+en+"')"
            res, err := db.Query(sql_statement)
            defer res.Close()

            if err != nil {
                log.Fatal(err)
            }
            ////fmt.Println("Request Successful Executed")
            var s string
            for res.Next() {
                if err := res.Scan(&s); err != nil {
                    log.Fatal(err)
                }
                ////fmt.Println(c)
                prov_liste = append(prov_liste, s)

            }
            done = true
        }
        ////fmt.Println(count)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(prov_liste)
        return
    })
    // Lauching server
    //log.Fatal(http.ListenAndServe(":4445", router))

    // create a custom server
    fmt.Println("Server Starting at :4445")
    s := &http.Server{
        Addr:    ":4445",
        Handler: router, // use `http.DefaultServeMux`
    }

    cert := "fullchain.pem"
    key := "privkey.pem"
    // run server on port "9000"
    log.Fatal(s.ListenAndServeTLS(cert, key))
}
