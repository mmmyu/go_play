package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [<flags>] <input_file>.\n", os.Args[0])
}

const (
	out_date   = iota
	out_amount = iota
	out_desc   = iota
)

func convertDate(d string) string {
	t, err := time.Parse("1/2/2006", d)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func convertWFB(record []string) []string {
	out := make([]string, 7)
	amt, err := strconv.ParseFloat(record[1], 32)
	if err != nil {
		panic(err)
	}
	if amt >= 0.0 ||
		strings.Contains(record[4], " INVESTMENT") ||
		strings.Contains(record[4], "ONLINE TRANSFER") {
		return out
	}
	if strings.HasPrefix(record[4], "BILL PAY") &&
		!strings.Contains(record[4], "RECURRING") {
		return out
	}

	out[out_date] = convertDate(record[0])
	out[out_desc] = record[3] + " " + record[4]
	out[5] = strconv.FormatFloat(-amt, 'f', -1, 32)
	out[6] = strconv.FormatFloat(amt, 'f', -1, 32)
	return out
}

func convertAmex(record []string) []string {
	out := make([]string, 3)
	out[out_date] = convertDate(record[0])
	out[out_amount] = record[2]
	out[out_desc] = record[3] + " " + record[4]
	return out
}

func convertCap1(record []string) []string {
	out := make([]string, 3)
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

func convertCiti(record []string) []string {
	out := make([]string, 3)
	out[out_date] = convertDate(record[0])
	amt, err := strconv.ParseFloat(record[1][1:], 32)
	if err != nil {
		panic(err)
	}
	out[out_amount] = strconv.FormatFloat(-amt, 'f', -1, 32)
	out[out_desc] = record[2]
	return out
}

func convertChase(record []string) []string {
	out := make([]string, 3)
	out[out_date] = convertDate(record[1])
	if out[out_date] == "" {
		return out
	}
	out[out_amount] = record[4]
	out[out_desc] = record[3]
	return out
}

const (
	is_wfb   = iota
	is_amex  = iota
	is_cap1  = iota
	is_chase = iota
	is_citi  = iota
)

func guessFileType(fname string) (int, string) {
	switch {
	case strings.Contains(fname, "Checking1"):
		fmt.Fprintf(os.Stderr, "Format=WFB\n")
		return is_wfb, "wfb"
	case strings.Contains(fname, "ofx"):
		fmt.Fprintf(os.Stderr, "Format=Amex\n")
		return is_amex, "amex"
	case strings.Contains(fname, "export"):
		fmt.Fprintf(os.Stderr, "Format=Capital1\n")
		return is_cap1, "cap1"
	case strings.Contains(fname, "Activity"):
		fmt.Fprintf(os.Stderr, "Format=Chase\n")
		return is_chase, "chase"
	case strings.Contains(fname, "xls"):
		fmt.Fprintf(os.Stderr, "Format=Citi\n")
		return is_citi, "citi"
	default:
		panic("Unknown file name type")
	}
}

func ftypeToEnum(ft string) (int, string) {
	switch ft {
	case "wfb": return is_wfb, ft
	case "amex": return is_amex, ft
	case "cap1": return is_cap1, ft
	case "chase": return is_chase, ft
	case "citi": return is_citi, ft
	default: panic("Unknown file type " + ft)
	}
}

func convert(ftype int, fin string, out_file *os.File) {
	fi, err := os.Open(fin)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	reader := csv.NewReader(fi)
	switch ftype {
	case is_citi:
		reader.Comma = '\t'
	}
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
		case is_amex:
			out = convertAmex(record)
		case is_cap1:
			out = convertCap1(record)
		case is_citi:
			out = convertCiti(record)
		case is_chase:
			out = convertChase(record)
		default:
			panic("Unknown file type " + string(ftype))
		}
		if out[out_date] != "" {
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
	fname := ftypename + "_" + time.Now().Format("20060102") + ".csv"
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
