package wordcount

import (
	"fmt"
	"testing"

	"code.google.com/p/goprotobuf/proto"
	"shipshape/test"

	notespb "shipshape/proto/note_proto"
)

func TestWordCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"This is a text with seven words", 7},
		{"\nThis\n is\n a\n text\n with\n seven\n words", 7},
		{"\tThis\t is\t a\t text\t with\t seven\t words", 7},
	}
	var w WordCountAnalyzer
	for _, test := range tests {
		got := w.CountWords(test.input)
		if got != test.want {
			t.Errorf("CountWords(%q): got %d, want %d", test.input, got, test.want)
		}
	}
}

func TestAnalyze(t *testing.T) {
	tests := []struct {
		file  string
		words int
	}{
		{"simple.txt", 6},
		{"complex.txt", 10},
		{"empty.txt", 0},
		{"whitespace.txt", 0},
	}

	var w WordCountAnalyzer
	for _, pair := range tests {
		ctx, err := test.CreateContext("shipshape/test_data/wordcount", []string{pair.file})
		if err != nil {
			t.Fatalf("error from CreateContext: %v", err)
		}

		actualNotes, err := test.RunAnalyzer(ctx, w, t)

		if err != nil {
			t.Errorf("received an analysis failure: %v", err)
		}
		expectedNotes := []*notespb.Note{
			&notespb.Note{
				Category:    proto.String("WordCount"),
				Description: proto.String(fmt.Sprintf("Word count: %d", pair.words)),
			},
		}

		pass, message := test.CheckNoteContainsContent(expectedNotes, actualNotes)
		if !pass {
			t.Errorf(message)
		}
	}
}

func TestAnalyzeFailure(t *testing.T) {
	tests := []string{"nonexistentfile.txt"}

	var w WordCountAnalyzer
	for _, input := range tests {
		ctx, err := test.CreateContext("shipshape/test_data/wordcount", []string{input})
		if err != nil {
			t.Errorf("error from CreateContext: %v", err)
		}

		notes, err := w.Analyze(ctx)
		if err == nil {
			t.Errorf("expected an analysis failure for input %s", input)
		}
		if notes != nil {
			t.Errorf("received notes %v for input %s", notes, input)
		}
	}
}