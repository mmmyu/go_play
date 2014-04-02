package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <input_file>.\n", os.Args[0])
}

func parseArgs(args []string) (string, string) {
	if len(args) != 3 {
	}
	return args[1], args[2]
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
		strings.HasPrefix(record[4], "BILL PAY") ||
		strings.Contains(record[4], " INVESTMENT") ||
		strings.Contains(record[4], "ONLINE TRANSFER") {
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
	is_wfb = iota
	is_amex = iota
	is_cap1 = iota
	is_chase = iota
	is_citi = iota
)

func guessFileType(fname string) int {
	switch {
	case strings.Contains(fname, "Checking1"): return is_wfb
	case strings.Contains(fname, "ofx"): return is_amex
	case strings.Contains(fname, "export"): return is_cap1
	case strings.Contains(fname, "Activity"): return is_chase
	case strings.Contains(fname, "xls"): return is_citi
	default: panic("Unknown file name type")
	}
}

func convert(ftype int, fin string) {
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
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()
	for _, r := range output_records {
		writer.Write(r)
	}
}

// Parses banks checking activity csv file
func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}
	fname := os.Args[1]
	convert(guessFileType(fname), fname)
}
