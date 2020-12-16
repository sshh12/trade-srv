package scraping

import (
	"regexp"
	"strings"
	"unicode"
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

// FixForUTF8 removes bytes to get things to work
func FixForUTF8(input []byte) string {
	clean := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, string(input))
	return clean
}
