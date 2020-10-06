package scraping

import (
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
	{"Â½", ""},
	{"®", ""},
	{"\xa0", " "},
	{"✅", ""},
	{"→", "->"},
	{"💯", ""},
	{"🚨", ""},
	// Tokens
	{" (Updated)", ""},
}

func regexReplace(s string, regex string, new string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(s, new)
}

func CleanHTMLText(raw string) string {
	for _, repl := range wordsRepls {
		raw = strings.ReplaceAll(raw, repl[0], repl[1])
	}
	raw = regexReplace(raw, "<style[\\s\\w=\":/\\.\\-,\\'!%&+@\\|{}\\(\\);#~\\?]*>([\\s\\S]+?)<\\/style>", "")
	raw = regexReplace(raw, "<script[\\s\\w=\":/\\.\\-,\\'!%&+@\\|{}\\(\\);#~\\?]*>([\\s\\S]+?)<\\/script>", "")
	raw = regexReplace(raw, "<\\w+[\\s\\w=\":/\\.\\-,\\'!%&+@\\|#~{}\\(\\);\\?]*>", "")
	raw = regexReplace(raw, "<\\/?[\\w\\-]+>", "")
	raw = regexReplace(raw, "<!-*[^>]+>", "")
	raw = regexReplace(raw, "&#[\\w\\d]+;", "")
	raw = regexReplace(raw, "\\s{3,}", "")
	raw = regexReplace(raw, "https:\\/\\/t.co\\/[\\w]+", "")
	raw = regexReplace(raw, "RT @\\w+:", "")
	raw = regexReplace(raw, "([a-z])\\s{2,}([A-Z])", "\\1 \\2")
	return strings.TrimSpace(raw)
}
