package parser

import (
	"reflect"
	"strings"
	"testing"
)

type testDef struct {
	input string
	want  []linkDef
}

func runDef(def testDef, t *testing.T) {
	got, err := ParseHTML(strings.NewReader(def.input))

	if err != nil {
		t.Fatalf("Parsing failed with error:\n%+v", err)
	}

	for i, result := range got {
		if !reflect.DeepEqual(result, def.want[i]) {
			t.Fatalf("want:\n%+v\nReceived:\n%+v", def.want[i], result)
		}
	}
}

func getCases() []testDef {
	noAnchors := `<div>SomeText</div>`
	noAnchorsResult := []linkDef{}
	noAnchorsDef := testDef{noAnchors, noAnchorsResult}

	oneAnchor := `<div><a href="/something.com"><div>Text!</div><div><em>yelling! at! stuff!</em></div></a></div>`
	oneAnchorResult := []linkDef{{Href: "/something.com", Text: "Text! yelling! at! stuff!"}}
	oneAnchorDef := testDef{oneAnchor, oneAnchorResult}

	anchorWithComments := `
	<div><a href="/something.com"><button>click</button><!-- commented text SHOULD NOT be included! --></a></div>
	`
	anchorWithCommentsResult := []linkDef{{Href: "/something.com", Text: "click"}}
	anchorWithCommentsDef := testDef{anchorWithComments, anchorWithCommentsResult}

	multipleAnchors := `
		<html><body><div><a href="/first.com"><button>first button</button></a></div><div><a href="/second.com"><button>second button</button></a></div></body></html>	
	`
	multipleAnchorsResult := []linkDef{{Href: "/first.com", Text: "first button"}, {Href: "/second.com", Text: "second button"}}
	multipleAnchorsDef := testDef{multipleAnchors, multipleAnchorsResult}

	spaceTextAnchor := `
		<html><body><div><a href="/first.com"><button>first button
		
		</button></a></div></body></html>	
	`
	spaceTextAnchorResult := []linkDef{{Href: "/first.com", Text: "first button"}}
	spaceTextAnchorsDef := testDef{spaceTextAnchor, spaceTextAnchorResult}

	return []testDef{noAnchorsDef, oneAnchorDef, anchorWithCommentsDef, multipleAnchorsDef, spaceTextAnchorsDef}
}

func TestParser(t *testing.T) {
	defs := getCases()

	for _, d := range defs {
		runDef(d, t)
	}
}
