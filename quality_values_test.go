package httputil_test

import (
	"testing"

	"github.com/d5/go-httputil"
)

func TestParseQualityValues(t *testing.T) {
	testParseQualityValues(t, "", []httputil.QualityValue{}, nil)

	testParseQualityValues(t, "text/html", []httputil.QualityValue{
		{MIMEType: "text/html", Priority: 1},
	}, nil)

	testParseQualityValues(t, "text/html;q=0.66", []httputil.QualityValue{
		{MIMEType: "text/html", Priority: 0.66},
	}, nil)

	testParseQualityValues(t, "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", []httputil.QualityValue{
		{MIMEType: "text/html", Priority: 1},
		{MIMEType: "application/xhtml+xml", Priority: 1},
		{MIMEType: "application/xml", Priority: 0.9},
		{MIMEType: "*/*", Priority: 0.8},
	}, nil)

	testParseQualityValues(t, "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8", []httputil.QualityValue{
		{MIMEType: "text/html", Priority: 1},
		{MIMEType: "application/xhtml+xml", Priority: 1},
		{MIMEType: "application/xml", Priority: 0.9},
		{MIMEType: "*/*", Priority: 0.8},
	}, nil)

	testParseQualityValuesExpectError(t, "text/html/invalid")
	testParseQualityValuesExpectError(t, "text/html?wrong=0.1")
	testParseQualityValuesExpectError(t, "text/html?q=2.0")
	testParseQualityValuesExpectError(t, "text/html?q=")
}

func TestSortQualityValues(t *testing.T) {
	httputil.SortQualityValues(nil) // nothing happens

	v := make([]httputil.QualityValue, 0) // empty array
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/html", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/html", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type2", Priority: 0},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type2", Priority: 0},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type2", Priority: 0},
		{MIMEType: "text/type1", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type2", Priority: 0},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type2", Priority: 0},
		{MIMEType: "text/type1", Priority: 0.5},
		{MIMEType: "text/type3", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type3", Priority: 1},
		{MIMEType: "text/type1", Priority: 0.5},
		{MIMEType: "text/type2", Priority: 0},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type2", Priority: 0.2},
		{MIMEType: "text/type3", Priority: 0.8},
		{MIMEType: "text/type4", Priority: 0},
		{MIMEType: "text/type5", Priority: 0.3},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type3", Priority: 0.8},
		{MIMEType: "text/type5", Priority: 0.3},
		{MIMEType: "text/type2", Priority: 0.2},
		{MIMEType: "text/type4", Priority: 0},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type2", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type1", Priority: 1},
		{MIMEType: "text/type2", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type2", Priority: 1},
		{MIMEType: "text/type1", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type2", Priority: 1},
		{MIMEType: "text/type1", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/type3", Priority: 1},
		{MIMEType: "text/type2", Priority: 1},
		{MIMEType: "text/type1", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type3", Priority: 1},
		{MIMEType: "text/type2", Priority: 1},
		{MIMEType: "text/type1", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text/*", Priority: 1},
		{MIMEType: "text/type2", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type2", Priority: 1},
		{MIMEType: "text/*", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "text1/*", Priority: 1},
		{MIMEType: "text2/*", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text1/*", Priority: 1},
		{MIMEType: "text2/*", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "*/*", Priority: 1},
		{MIMEType: "text/type2", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/type2", Priority: 1},
		{MIMEType: "*/*", Priority: 1},
	}, v)

	v = []httputil.QualityValue{
		{MIMEType: "*/*", Priority: 1},
		{MIMEType: "text/*", Priority: 1},
	}
	httputil.SortQualityValues(v)
	assertEqualQualityValues(t, []httputil.QualityValue{
		{MIMEType: "text/*", Priority: 1},
		{MIMEType: "*/*", Priority: 1},
	}, v)
}

func testParseQualityValues(t *testing.T, input string, expectedValues []httputil.QualityValue, expectedError error) {
	actualValues, actualError := httputil.ParseQualityValues(input)
	if expectedError != nil {
		if actualError != expectedError {
			t.Errorf("Expected Error: %v, Actual: %v\n", expectedError, actualError)
		}
		return
	}

	assertEqualQualityValues(t, expectedValues, actualValues)
}

func assertEqualQualityValues(t *testing.T, expected, actual []httputil.QualityValue) {
	if len(expected) != len(actual) {
		t.Errorf("Expected Values: %v, Actual: %v\n", expected, actual)
	}

	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected Values: %v, Actual: %v\n", expected, actual)
			return
		}
	}
}

func testParseQualityValuesExpectError(t *testing.T, input string) {
	_, actualError := httputil.ParseQualityValues(input)
	if actualError == nil {
		t.Errorf("Expected Error: %v, Actual: nil\n", actualError)
	}
}
