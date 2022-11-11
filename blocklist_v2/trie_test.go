package blocklist_v2

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/miekg/dns"
)

const TestBlockList = `. [BlockPrefix:AAAA]
bbbb.aaaa. [BlockPrefix:NS]
cccc.bbbb.aaaa. [Block:ANY]
dddd.bbbb.aaaa. [Block:CNAME Block:TXT]
xxxx. [BlockPrefix:ANY]
`

func TestTrieAdd(t *testing.T) {

	root := NewLevel()

	root.Add([]string{}, BlockPrefixQtype{dns.TypeAAAA})
	root.Add([]string{"bbbb", "aaaa"}, BlockPrefixQtype{dns.TypeNS})
	root.Add([]string{"cccc", "bbbb", "aaaa"}, Block{})
	root.Add([]string{"dddd", "bbbb", "aaaa"}, BlockQtype{dns.TypeCNAME})
	root.Add([]string{"dddd", "bbbb", "aaaa"}, BlockQtype{dns.TypeTXT})
	root.Add([]string{"xxxx"}, BlockPrefix{})

	buf := bytes.Buffer{}
	root.PrintTree(&buf, []string{})

	if buf.String() != TestBlockList {
		t.Error(buf.String())
	}

	out := []BlockEntry{}
	root.Dump([]string{}, &out)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	for _, v := range out {
		fmt.Println(v)
	}
}