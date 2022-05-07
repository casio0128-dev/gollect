package utils

import (
	"reflect"
	"regexp"
	"testing"
)

var (
	targetFileNames = []string{
		"sample_12345.png", "sample_67890.txt", "12345_sample.png", "67890_sample.md",
		"csv_data1.csv", "csv_data2.csv", "csv_data3.csv", "s.txt", "s.md", "s.png", "s.jpg",
		"s.s", // s.sは指定文字数がファイル名を超過した場合の確認用
	}
)

func TestCreatePatternMapByRegexp(t *testing.T) {
	var (
		hasCharactersPatternSample, _       = regexp.Compile(`sample`)      // 「sample」を含むものを抽出
		hasCharactersPatternCSV, _          = regexp.Compile(`csv`)         // 「csv」を含むものを抽出
		hasCharactersPatternFiveIntegers, _ = regexp.Compile(`\d{5}`)       // 数字５桁一致での抽出
		extensionClassificationPattern, _   = regexp.Compile(`\.[^\.]\w*$`) // 拡張子での抽出
	)

	type args struct {
		fileNames     []string
		regexpPattern **regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{name: "「sample」を含むファイルを抽出", args: args{fileNames: targetFileNames, regexpPattern: &hasCharactersPatternSample},
			want: map[string][]string{
				"sample": {"sample_12345.png", "sample_67890.txt", "12345_sample.png", "67890_sample.md"},
			}},
		{name: "「csv」を含むファイルを抽出", args: args{fileNames: targetFileNames, regexpPattern: &hasCharactersPatternCSV},
			want: map[string][]string{
				"csv": {"csv_data1.csv", "csv_data2.csv", "csv_data3.csv"},
			}},
		{name: "5桁の数字を含むファイルを抽出", args: args{fileNames: targetFileNames, regexpPattern: &hasCharactersPatternFiveIntegers},
			want: map[string][]string{
				"12345": {"sample_12345.png", "12345_sample.png"},
				"67890": {"sample_67890.txt", "67890_sample.md"},
			}},
		{name: "拡張子でのファイル抽出", args: args{fileNames: targetFileNames, regexpPattern: &extensionClassificationPattern},
			want: map[string][]string{
				".png": {"sample_12345.png", "12345_sample.png", "s.png"},
				".txt": {"sample_67890.txt", "s.txt"},
				".md":  {"67890_sample.md", "s.md"},
				".csv": {"csv_data1.csv", "csv_data2.csv", "csv_data3.csv"},
				".jpg": {"s.jpg"},
				".s":   {"s.s"},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreatePatternMapByRegexp(tt.args.fileNames, tt.args.regexpPattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePatternMapByRegexp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreatePatternMapByStringLength(t *testing.T) {
	var (
		pTrue  = true
		pFalse = false
		five   = uint(5)
		three  = uint(3)
	)
	type args struct {
		fileNames    []string
		isHead       *bool
		isTail       *bool
		isIncludeExt *bool
		length       *uint
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{name: "先頭から5文字の一致で抽出", args: args{
			fileNames:    targetFileNames,
			isHead:       &pTrue,
			isTail:       &pFalse,
			isIncludeExt: &pFalse,
			length:       &five,
		}, want: map[string][]string{
			"sampl": {"sample_12345.png", "sample_67890.txt"},
			"12345": {"12345_sample.png"},
			"67890": {"67890_sample.md"},
			"csv_d": {"csv_data1.csv", "csv_data2.csv", "csv_data3.csv"},
			"s.txt": {"s.txt"},
			"s.md":  {"s.md"},
			"s.png": {"s.png"},
			"s.jpg": {"s.jpg"},
			"s.s":   {"s.s"},
		}},
		{name: "末尾から3文字の一致で抽出（拡張子含める）", args: args{
			fileNames:    targetFileNames,
			isHead:       &pFalse,
			isTail:       &pTrue,
			isIncludeExt: &pTrue,
			length:       &three,
		}, want: map[string][]string{
			"png": {"sample_12345.png", "12345_sample.png", "s.png"},
			"txt": {"sample_67890.txt", "s.txt"},
			".md": {"67890_sample.md", "s.md"},
			"csv": {"csv_data1.csv", "csv_data2.csv", "csv_data3.csv"},
			"jpg": {"s.jpg"},
			"s.s": {"s.s"},
		}},
		{name: "末尾から3文字の一致で抽出（拡張子含めず）", args: args{
			fileNames:    targetFileNames,
			isHead:       &pFalse,
			isTail:       &pTrue,
			isIncludeExt: &pFalse,
			length:       &three,
		}, want: map[string][]string{
			"345": {"sample_12345.png"},
			"890": {"sample_67890.txt"},
			"ple": {"12345_sample.png", "67890_sample.md"},
			"ta1": {"csv_data1.csv"},
			"ta2": {"csv_data2.csv"},
			"ta3": {"csv_data3.csv"},
			"s":   {"s.txt", "s.md", "s.png", "s.jpg", "s.s"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreatePatternMapByStringLength(tt.args.fileNames, tt.args.isHead, tt.args.isTail, tt.args.isIncludeExt, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePatternMapByStringLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
