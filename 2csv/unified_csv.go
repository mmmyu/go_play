package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [<flags>] <input_file>.\n", os.Args[0])
}

const (
	out_date        = iota
	out_amount      = iota
	out_desc        = iota
	out_links       = iota
	out_comment     = iota
	out_payment     = iota
	out_cash        = iota
	out_food        = iota
	out_restaurant  = iota
	out_gas         = iota
	out_car         = iota
	out_online      = iota
	out_drugs       = iota
	out_school      = iota
	out_kids        = iota
	out_target      = iota
	out_walmart     = iota
	out_costco      = iota
	out_store       = iota
	out_fashion     = iota
	out_home_improv = iota
	out_electronics = iota
	out_doctors     = iota
	out_utilities   = iota
	out_hoa         = iota
	out_misc        = iota
	out_travel      = iota
	out_tax         = iota
	out_mortgage    = iota
	out_num_fields  = iota
)

var jump_table = [out_num_fields]func(r []string){
	nil, // date
	nil, // amount
	nil, // desc
	nil, // links
	nil, // comment
	classify_payment,
	nil, // cash
	classify_food,
	classify_restaurant,
	classify_gas,
	nil, // car
	classify_online,
	classify_drugs,
	classify_school,
	nil, // kids
	classify_target,
	classify_walmart,
	classify_costco,
	nil, // store
	classify_fashion,
	classify_home_improv,
	classify_electronics,
	nil, // doctors
	classify_utilities,
	classify_hoa,
	nil, // misc
	classify_travel_sports,
	classify_tax_insurance,
	classify_mortgage,
}

func convertDate(d string) string {
	t, err := time.Parse("1/2/2006", d)
	if err != nil {
		t, err = time.Parse("2006-1-2", d)
		if err != nil {
			return ""
		}
	}
	return t.Format("2006-01-02")
}

func has_string(a, b string) bool {
	return strings.Contains(strings.ToLower(a),
		strings.ToLower(b))
}

func classify_fashion(r []string) {
	if has_string(r[out_desc], "FOREVER 21") ||
		has_string(r[out_desc], "6PM.COM") ||
		has_string(r[out_desc], "VICTORIA'S") ||
		has_string(r[out_desc], "EAGLE OUTFTR") ||
		has_string(r[out_desc], "Macys") {
		r[out_fashion] = r[out_amount]
	}
}

func classify_tax_insurance(r []string) {
	if has_string(r[out_desc], "CSAA Auto") {
		r[out_tax] = r[out_amount]
	}
}

func classify_mortgage(r []string) {
	if has_string(r[out_desc], "PROVIDENT FUND") ||
		has_string(r[out_desc], "CALIBER") ||
		has_string(r[out_desc], "Vineyard Village") ||
		has_string(r[out_desc], "CHASE HOME FINAN") {
		r[out_mortgage] = r[out_amount]
	}
}

func classify_school(r []string) {
	if has_string(r[out_desc], "TAEKWON") ||
		has_string(r[out_desc], "TIFFANY'S DANCE") ||
		has_string(r[out_desc], "TIFFANYSDAN") ||
		has_string(r[out_desc], "HAPPY FISH") {
		r[out_school] = r[out_amount]
	}
}

func classify_drugs(r []string) {
	if has_string(r[out_desc], "CVS PHARMACY") ||
		has_string(r[out_desc], "CVS/PHARMACY") {
		r[out_drugs] = r[out_amount]
	}
}

func classify_hoa(r []string) {
	if has_string(r[out_desc], "PARC METRO") {
		r[out_hoa] = r[out_cash]
		r[out_cash] = ""
	}
}

func classify_walmart(r []string) {
	if has_string(r[out_desc], "WAL-MART") {
		r[out_walmart] = r[out_cash]
		r[out_cash] = ""
	}
}

func classify_home_improv(r []string) {
	if has_string(r[out_desc], "Home Depot") ||
		has_string(r[out_desc], "LOWES") ||
		has_string(r[out_desc], "RICHERT LUMBER") ||
		has_string(r[out_desc], "ORCHARD SUPPLY") {
		r[out_home_improv] = r[out_amount]
	}
}

func classify_costco(r []string) {
	if has_string(r[out_desc], "Costco WHSE") {
		r[out_costco] = r[out_amount]
	}
}

func classify_electronics(r []string) {
	if has_string(r[out_desc], "www.newegg.com") {
		r[out_electronics] = r[out_amount]
	}
}

func classify_target(r []string) {
	if has_string(r[out_desc], "target") {
		r[out_target] = r[out_amount]
	}
}

func classify_utilities(r []string) {
	if has_string(r[out_desc], "AT&T*BILL") ||
		has_string(r[out_desc], "OOMA") ||
		has_string(r[out_desc], "PLEASANTON WATER") ||
		has_string(r[out_desc], "COMCAST") ||
		has_string(r[out_desc], "TMOBILE") {
		r[out_utilities] = r[out_amount]
	}
	if has_string(r[out_desc], "PGANDE") {
		r[out_utilities] = r[out_cash]
		r[out_cash] = ""
	}

}

func classify_online(r []string) {
	if has_string(r[out_desc], "netflix.com") ||
		has_string(r[out_desc], "GOOGLE *") ||
		has_string(r[out_desc], "Audible") {
		r[out_online] = r[out_amount]
	}
}

func classify_gas(r []string) {
	if has_string(r[out_desc], "Chevron") ||
		has_string(r[out_desc], "Costco gas") ||
		has_string(r[out_desc], "Union 76") ||
		has_string(r[out_desc], "76 fuel") ||
		has_string(r[out_desc], "76 - ") ||
		has_string(r[out_desc], "SAFEWAY FUEL") ||
		has_string(r[out_desc], "Shell Oil") ||
		has_string(r[out_desc], "valero") ||
		has_string(r[out_desc], "FASTRAK") {
		r[out_gas] = r[out_amount]
	}
}

func classify_restaurant(r []string) {
	if has_string(r[out_desc], "Starbucks") ||
		has_string(r[out_desc], "Tully") ||
		has_string(r[out_desc], "Peet") ||
		has_string(r[out_desc], "SUBWAY") ||
		has_string(r[out_desc], "JUST KOI") ||
		has_string(r[out_desc], "YAKITORI") ||
		has_string(r[out_desc], "CAFE") ||
		has_string(r[out_desc], "GOLDEN SAND HARBOR") ||
		has_string(r[out_desc], "T4") {
		r[out_restaurant] = r[out_amount]
	}
}

func classify_payment(r []string) {
	if has_string(r[out_desc], "PAYMENT") ||
		has_string(r[out_desc], "AUTOPAY") {
		v, err := strconv.ParseFloat(r[out_amount], 32)
		if err == nil && v >= 0.0 {
			r[out_payment] = r[out_amount]
		}
	}
}

func classify_food(r []string) {
	if has_string(r[out_desc], "99 RANCH") ||
		has_string(r[out_desc], "FOOD EXPRESS") ||
		has_string(r[out_desc], "KOI PALACE") ||
		has_string(r[out_desc], "MARINA") ||
		has_string(r[out_desc], "RALEY'S") ||
		has_string(r[out_desc], "TRADER JOE") ||
		has_string(r[out_desc], "WHOLE FOODS") ||
		has_string(r[out_desc], "KEE WAH BAKERY") {
		r[out_food] = r[out_amount]
	}
}

func classify_travel_sports(r []string) {
	is_rei, _ := regexp.MatchString("^REI ", r[out_desc])
	if has_string(r[out_desc], "SPORTS AUTH") ||
		has_string(r[out_desc], "HOTELS.COM") ||
		is_rei ||
		has_string(r[out_desc], "SPORT CHALET") {
		r[out_travel] = r[out_amount]
	}
}

func classify(record []string) {
	for i := 0; i < out_num_fields; i++ {
		if record[i] == "" && jump_table[i] != nil {
			jump_table[i](record)
		}
	}
}

func convertWFB(record []string) []string {
	out := make([]string, out_num_fields)
	amt, err := strconv.ParseFloat(record[1], 32)
	if err != nil {
		panic(err)
	}
	if amt >= 0.0 ||
		strings.Contains(record[4], " INVESTMENT") ||
		strings.Contains(record[4], "ONLINE TRANSFER") {
		return out
	}
	// Skip credit card autopay entries
	if strings.HasPrefix(record[4], "CHASE AUTOPAY") ||
		strings.HasPrefix(record[4], "AMERICAN EXPRESS ACH PMT") ||
		strings.HasPrefix(record[4], "CAPITAL ONE CRCARDPMT") ||
		strings.HasPrefix(record[4], "CITI AUTO PAYMENT") {
		return out
	}
	if !strings.HasPrefix(record[4], "BILL PAY") ||
		strings.Contains(record[4], "RECURRING") ||
		strings.Contains(record[4], "LIFEINS") ||
		strings.Contains(record[4], "CSAA") ||
		strings.Contains(record[4], "AAA") {
		out[out_date] = convertDate(record[0])
		out[out_desc] = record[3] + " " + record[4]
		out[5] = strconv.FormatFloat(-amt, 'f', -1, 32)
		out[6] = strconv.FormatFloat(amt, 'f', -1, 32)
	}
	return out
}

func convertWFBCard(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[0])
	out[out_desc] = record[4]
	out[out_amount] = record[1]
	return out
}

func convertAmex(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[0])
	out[out_amount] = record[2]
	out[out_desc] = record[3] + " " + record[4]
	return out
}

func convertCap1(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[0])
	if out[out_date] == "" {
		return out
	}
	debit, err := strconv.ParseFloat(record[3], 32)
	if err != nil {
		debit = 0
	}
	credit, err := strconv.ParseFloat(record[4], 32)
	if err != nil {
		credit = 0
	}
	amt := credit - debit
	out[out_amount] = strconv.FormatFloat(amt, 'f', -1, 32)
	out[out_desc] = record[2]
	return out
}

func convertCap1Transactions(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[1])
	if out[out_date] == "" {
		return out
	}
	debit, err := strconv.ParseFloat(record[6], 32)
	if err != nil {
		debit = 0
	}
	credit, err := strconv.ParseFloat(record[7], 32)
	if err != nil {
		credit = 0
	}
	amt := credit - debit
	out[out_amount] = strconv.FormatFloat(amt, 'f', -1, 32)
	out[out_desc] = record[4]
	return out
}

func convertCitiStmt(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[0])
	if out[out_date] == "" {
		return out
	}
	amt, _ := strconv.ParseFloat(record[1][1:len(record[1])], 32)
	out[out_amount] = strconv.FormatFloat(-amt, 'f', -1, 32)
	out[out_desc] = strings.Trim(record[2], "\n ")
	return out
}

func convertCiti(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[1])
	if out[out_date] == "" {
		return out
	}
	debit, _ := strconv.ParseFloat(record[3], 32)
	credit, _ := strconv.ParseFloat(record[4], 32)
	amt := credit - debit
	out[out_amount] = strconv.FormatFloat(amt, 'f', -1, 32)
	out[out_desc] = strings.Trim(record[2], "\n ")
	return out
}

func convertChase(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[1])
	if out[out_date] == "" {
		return out
	}
	out[out_amount] = record[4]
	out[out_desc] = record[3]
	return out
}

func convertRed(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[0])
	if out[out_date] == "" {
		return out
	}
	out[out_amount] = record[3]
	var regex_sep = regexp.MustCompile(" +")
	splits := regex_sep.Split(record[2], -1)
	city := splits[len(splits)-2]
	state := splits[len(splits)-1]
	out[out_desc] = "Target " + city + " " + state
	return out
}

func convertMacys(record []string) []string {
	out := make([]string, out_num_fields)
	out[out_date] = convertDate(record[0])
	if out[out_date] == "" {
		return out
	}
	amt, _ := strconv.ParseFloat(record[1][1:len(record[1])], 32)
	out[out_amount] = strconv.FormatFloat(-amt, 'f', -1, 32)
	out[out_desc] = "Macys " + record[2]
	return out
}

const (
	is_wfb       = iota // Wells Fargo Checking Account
	is_amex      = iota // American Express
	is_cap1      = iota // Capital One
	is_cap1trans = iota // Captial One Transactions format
	is_chase     = iota // Chase
	is_citi      = iota // Citi
	is_citiStmt  = iota // Citi from statement download page
	is_wfbc      = iota // Wells Fargo Credit Card
	is_red       = iota // Target Red Card
	is_macys     = iota // Macy's card
)

func guessFileType(fname string) (int, string) {
	citi_stmt, _ := regexp.MatchString("[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9].csv", fname)
	macys_stmt, _ := regexp.MatchString("[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9].txt", fname)
	switch {
	case strings.Contains(fname, "Checking1"):
		fmt.Fprintf(os.Stderr, "Format=WFB\n")
		return is_wfb, "wfb"
	case strings.Contains(fname, "CreditCard"):
		fmt.Fprintf(os.Stderr, "Format=WFBC\n")
		return is_wfbc, "wfbc"
	case strings.Contains(fname, "ofx"):
		fmt.Fprintf(os.Stderr, "Format=Amex\n")
		return is_amex, "amex"
	case strings.Contains(fname, "export"):
		fmt.Fprintf(os.Stderr, "Format=Capital1\n")
		return is_cap1, "cap1"
	case strings.Contains(fname, "Transactions"):
		fmt.Fprintf(os.Stderr, "Format=Capital1Trans\n")
		return is_cap1trans, "cap1t"
	case strings.Contains(fname, "Activity"):
		fmt.Fprintf(os.Stderr, "Format=Chase\n")
		return is_chase, "chase"
	case strings.Contains(fname, "CURRENT_VIEW.CSV"):
		fmt.Fprintf(os.Stderr, "Format=Citi\n")
		return is_citi, "citi"
	case citi_stmt:
		fmt.Fprintf(os.Stderr, "Format=CitiStmt\n")
		return is_citiStmt, "citistmt"
	case strings.Contains(fname, "download"):
		fmt.Fprintf(os.Stderr, "Format=TRed\n")
		return is_red, "tred"
	case macys_stmt:
		fmt.Fprintf(os.Stderr, "Format=Macys\n")
		return is_macys, "macys"
	default:
		panic("Unknown file name type")
	}
}

func ftypeToEnum(ft string) (int, string) {
	switch ft {
	case "wfb":
		return is_wfb, ft
	case "wfbc":
		return is_wfbc, ft
	case "amex":
		return is_amex, ft
	case "cap1":
		return is_cap1, ft
	case "chase":
		return is_chase, ft
	case "citi":
		return is_citi, ft
	case "red":
		return is_red, ft
	case "cap1t":
		return is_cap1trans, ft
	case "macys":
		return is_macys, ft
	default:
		panic("Unknown file type " + ft)
	}
}

func convert(ftype int, fin string, out_file *os.File) {
	fi, err := os.Open(fin)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	reader := csv.NewReader(fi)
	output_records := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		var out []string
		switch ftype {
		case is_wfb:
			out = convertWFB(record)
		case is_wfbc:
			out = convertWFBCard(record)
		case is_amex:
			out = convertAmex(record)
		case is_cap1:
			out = convertCap1(record)
		case is_cap1trans:
			out = convertCap1Transactions(record)
		case is_citi:
			out = convertCiti(record)
		case is_citiStmt:
			out = convertCitiStmt(record)
		case is_chase:
			out = convertChase(record)
		case is_red:
			out = convertRed(record)
		case is_macys:
			out = convertMacys(record)
		default:
			panic("Unknown file type " + string(ftype))
		}
		if out[out_date] != "" {
			classify(out)
			output_records = append(output_records, out)
		}
	}
	//	fo, err := os.Create(fout)
	writer := csv.NewWriter(out_file)
	defer writer.Flush()
	for _, r := range output_records {
		writer.Write(r)
	}
}

func getOutputFile(ftype int, ftypename string) *os.File {
	fname := ftypename + "_" + time.Now().Format("20060102150405") + ".csv"
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	return f
}

// Parses banks checking activity csv file
func main() {
	flag_output := flag.Bool("o", false, "Writes to a csv file.")
	flag_type := flag.String("t", "auto", "Specifies file type.")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("narg=", flag.NArg())
		usage()
		return
	}
	fname := flag.Arg(0)
	var ftype int
	var ftypename string
	if *flag_type == "auto" {
		ftype, ftypename = guessFileType(fname)
	} else {
		ftype, ftypename = ftypeToEnum(*flag_type)
	}
	out_file := os.Stdout
	if *flag_output {
		out_file = getOutputFile(ftype, ftypename)
		defer out_file.Close()
	}
	convert(ftype, fname, out_file)
}
