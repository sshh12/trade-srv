package scraping

import (
	"bytes"
	"regexp"
	"strings"
)

var wordsRepls [][]string = [][]string{
	// HTML codes
	{"&rsquo;", "'"},
	{"&lsquo;", "'"},
	{"&ldquo;", "\""},
	{"&rdquo;", "\""},
	{"&quot;", "\""},
	{"&amp;", "&"},
	{"&ntilde;", "n"},
	{"&copy;", ""},
	{"&#8217;", "'"},
	{"&#160;", " "},
	{"&nbsp;", " "},
	{"&otilde;", "o"},
	{"&ccedil;", "c"},
	{"&lt;", "<"},
	{"&gt;", ">"},
	{"&ndash;", "-"},
	{"&mdash;", "-"},
	{"&uuml;", "u"},
	{"&oacute;", "o"},
	{"&mu;", "μ"},
	{"&eacute;", "e"},
	{"&ouml;", "o"},
	{"&reg;", ""},
	{"&auml;", "a"},
	{"&iacute;", "i"},
	{"&uacute;", "u"},
	{"&raquo;", "\""},
	{"&laquo;", "\""},
	{"&ocirc;", "o"},
	{"&agrave;", "a"},
	{"&Eacute;", "E"},
	{"&ucirc;", "u"},
	{"&Agrave;", "A"},
	{"&egrave;", "e"},
	{"&ugrave;", "u"},
	{"&aacute;", "a"},
	{"&ocirc;", "o"},
	{"&trade;", ""},
	{"&times;", "x"},
	{"&bull;", "* "},
	{"&#x27;", "'"},
	{"&minus;", "-"},
	{"&apos;", "'"},
	// Weird chars
	{"•", "*"},
	{"●", "* "},
	{"\u2019", "'"},
	{"\r", ""},
	{"…", "..."},
	{"—", "-"},
	{"ー", "-"},
	{"‘", "'"},
	{"’", "'"},
	{"“", ""},
	{"”", ""},
	{"»", "\""},
	{"«", "\""},
	{"™", ""},
	{"\u200d", ""},
	{"\u2013", "-"},
	{"Â\xa0", " "},
	{"\xc2", " "},
	{"Â½", ""},
	{"®", ""},
	{"\xa0", " "},
	{"✅", ""},
	{"→", "->"},
	{"💯", ""},
	{"🚨", ""},
}

// RegexReplace does a replace all with the given regex
func RegexReplace(s string, regex string, new string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(s, new)
}

// CleanHTMLText cleans text extracted from HTML
func CleanHTMLText(raw string) string {
	for _, repl := range wordsRepls {
		raw = strings.ReplaceAll(raw, repl[0], repl[1])
	}
	raw = RegexReplace(raw, "<style[_\\s\\w=\":/\\.\\-,\\'!%$&+@\\|{}\\(\\);#~\\?]*>([\\s\\S]+?)<\\/style>", "")
	raw = RegexReplace(raw, "<script[_\\s\\w=\":/\\.\\-,\\'!%$&+@\\|{}\\(\\);#~\\?]*>([\\s\\S]+?)<\\/script>", "")
	raw = RegexReplace(raw, "<\\w+[_\\*\\s\\w=\":/\\.\\-,\\'!%$&+@\\|#~{}\\(\\);\\?]*>", "")
	raw = RegexReplace(raw, "<\\/?[\\w\\-]+>", "")
	raw = RegexReplace(raw, "<!-*[^>]+>", "")
	raw = RegexReplace(raw, "&#[\\w\\d]+;", "")
	raw = RegexReplace(raw, "\\s{3,}", "")
	raw = RegexReplace(raw, "https:\\/\\/t.co\\/[\\w]+", "")
	raw = RegexReplace(raw, "RT @\\w+:", "")
	raw = RegexReplace(raw, "([a-z])\\s{2,}([A-Z])", "$1 $2")
	raw = RegexReplace(raw, "([a-z%])([A-Z])", "$1 $2")
	raw = RegexReplace(raw, "%(\\w)", "% $2")
	raw = RegexReplace(raw, "([\\w\\)]),([+\\w])", "$1, $2")
	raw = RegexReplace(raw, "(\\w):([A-Za-z])", "$1: $2")
	raw = RegexReplace(raw, "([a-z])\\.([A-Z])", "$1. $2")
	raw = RegexReplace(raw, "(\\w)\\(", "$1 (")
	return strings.TrimSpace(raw)
}

// FixForUTF8 removes random bytes to get things to work w/utf-8
func FixForUTF8(input []byte) string {
	out := bytes.ReplaceAll(input, []byte{0xe6, 0x20, 0xbc}, []byte{})
	out = bytes.ReplaceAll(out, []byte{0xc3, 0x61}, []byte{})
	out = bytes.ReplaceAll(out, []byte{0xc3, 0x20}, []byte{})
	out = bytes.ReplaceAll(out, []byte{0xae}, []byte{})
	out = bytes.ReplaceAll(out, []byte{0xb2}, []byte{})
	out = bytes.ReplaceAll(out, []byte{0xb1}, []byte{})
	out = bytes.ReplaceAll(out, []byte{0xb0}, []byte{})
	return string(out)
}
