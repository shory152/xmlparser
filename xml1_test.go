package xmlparser

import (
	"fmt"
	"io"
	"testing"
	"time"
)

var xmlstr string = `
<?xml version=1.0 encoding=UTF-8 ?>
<books>
	<!-- 2 books - - -x- -- -->
	<book p1="v1" p2="v2">
		<name>
		
		 大 道 中 国 </name>
		<price> 89.00 </price>
		<author>张大中</author>
	</book>
	
	<book>
		<name>小猪唏哩呼噜</name>
		<price>22.50</price>
		<author>Alex</author>
	</book>
</books>
`

func TestXml(t *testing.T) {
	scanner := scanXml(xmlstr)
	for {
		if tk, err := scanner(); err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		} else {
			fmt.Println(tk)
		}
	}
	fmt.Println(scanner())

	fmt.Println("----------------------")
	scanner = scanXml2(xmlstr)
	for {
		if tk, err := scanner(); err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		} else {
			fmt.Println(tk)
		}
	}
	fmt.Println(scanner())

	fmt.Println("----------------------")
	scanner = scanXml3(xmlstr)
	for {
		if tk, err := scanner(); err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		} else {
			fmt.Println(tk)
		}
	}
	fmt.Println(scanner())
}

func nrun(label string, N int64, f func()) {
	bt := time.Now()
	for i := 0; i < int(N); i++ {
		f()
	}
	elapse := time.Since(bt)
	fmt.Printf("%s: %d ns/op, %.2f ops/s, elapse %.2fs\n", label,
		elapse.Nanoseconds()/N, float64(N)/elapse.Seconds(), elapse.Seconds())
}

func TestXmlPerf(t *testing.T) {
	N := int64(1000000)

	nrun("scanner1", N, func() {
		s := scanXml(xmlstr)
		for {
			if tk, err := s(); err != nil {
				if err != io.EOF {
					t.Error(err)
				}
				break
			} else {
				_ = tk
			}
		}
	})

	nrun("scanner2", N, func() {
		s := scanXml2(xmlstr)
		for {
			if tk, err := s(); err != nil {
				if err != io.EOF {
					t.Error(err)
				}
				break
			} else {
				_ = tk
			}
		}
	})

	nrun("scanner3", N, func() {
		s := scanXml3(xmlstr)
		for {
			if tk, err := s(); err != nil {
				if err != io.EOF {
					t.Error(err)
				}
				break
			} else {
				_ = tk
			}
		}
	})

}

func BenchmarkXmlScan(b *testing.B) {
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := scanXml(xmlstr)
		for {
			if tk, err := s(); err != nil {
				break
			} else {
				_ = tk
			}
		}
	}
	//b.StopTimer()
}

func TestParseXml(t *testing.T) {
	root, err := ParseXml(xmlstr)
	if err != nil {
		t.Error(err)
		return
	}
	ShowXml(root, nil, -1)
}
