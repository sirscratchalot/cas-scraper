package webbook

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	html := `<body>
    <p ><a class="fake" href="#test">
            Jump to content</a>
    </p>
    <header>
    </header>

    <main id="main">
        <h1 id="Top">Methane</h1>
        <ul>
            <li><strong><a title=""></a> </strong></li>
            <li><strong><a title=""
                        href="">Molecular weight</a>:</strong>
                16.0425</li>
            <li><strong>IUPAC Standard InChI:</strong></li>
            <li></li>
            <li><strong>CAS Registry Number:</strong> 74-82-8</li>
            <li><strong>Chemical structure:</strong> <img src="/cgi/cbook.cgi?Struct=C74828&amp;Type=Color"
                    class="struct" alt="CH4"></li>
            <li><strong>Options:</strong>
            </li>
            <p>
                dummy data
            </p>
    </main>
    <footer id="footer">

        <div class="row">
        </div>

    </footer>
<body>`
	results, err := matchXpath(html)
	if err != nil {
		t.Errorf("Got error during HTML parse %s\n", err.Error())
	}
	fmt.Printf("Parsed %+v\n", results)
}
