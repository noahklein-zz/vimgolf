package vim

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

func TestTempFileWithContents(t *testing.T) {
	f := func(s string) bool {
		tf, err := tempFileWithContents("", "temp-file-test", []byte(s))
		if err != nil {
			return false
		}
		defer os.Remove(tf.Name())
		got, err := ioutil.ReadFile(tf.Name())
		return s == string(got)
	}

	quick.Check(f, nil)
}

func BenchmarkTempFileWithContents(b *testing.B) {
	bytes := []byte("testing")
	for i := 0; i < b.N; i++ {
		tf, err := tempFileWithContents("", "benchmark-temp-file", bytes)
		if err != nil {
			b.Errorf("Failed to make temp file: %v", err)
		}
		defer os.Remove(tf.Name())
	}
}

func TestAttemptChallenge(t *testing.T) {
	cases := []struct {
		desc  string
		start string
		input string
		want  string
	}{
		{
			desc:  "single command",
			start: "hello world!",
			input: "f!x0A :)",
			want:  "hello world :)",
		},
		{
			desc:  "multiple commands",
			start: "hello\nworld\nhello!",
			input: "Oworld<Esc>:%s/hello/world/g",
			want:  "world\nworld\nworld\nworld!",
		},
		{
			desc:  "empty command",
			start: "hello world",
			input: " ",
			want:  "hello world",
		},
	}

	for _, c := range cases {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		got, err := AttemptChallenge(ctx, c.start, c.input)
		if err != nil {
			t.Error(err)
		}
		if got != fmt.Sprintf("%s\n", c.want) {
			t.Errorf("Case %s: got %v, want %v", c.desc, got, c.want)
		}
	}
}

func BenchmarkAttemptChallenge(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		AttemptChallenge(ctx, "Test", "dGiTest")
	}
}

func TestLexCommands(t *testing.T) {
	cases := []struct {
		desc  string
		input string
		want  []string
	}{
		{
			desc:  "split on <Esc>",
			input: "one<Esc>two",
			want:  []string{"one", "two"},
		},
		{
			desc:  "multiple <Esc>",
			input: "one<Esc><Esc>two",
			want:  []string{"one", "", "two"},
		},
	}

	for _, c := range cases {
		got := lexCommand(c.input)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Case %s: got %q, want %q", c.desc, got, c.want)
		}
	}
}

func TestScore(t *testing.T) {
	cases := []struct {
		desc  string
		input string
		want  int
	}{
		{desc: "simple", input: "12345", want: 5},
		{desc: "one <Esc>", input: "12<Esc>45", want: 5},
		{desc: "multiple <Esc>", input: "1<Esc><Esc>4<Esc>", want: 5},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			if got := Score(c.input); got != c.want {
				t.Errorf("Score() = %v, want %v", got, c.want)
			}
		})
	}
}
