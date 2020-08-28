package gui

import (
	"context"
	"reflect"
	"strings"
	"syscall/js"
	"testing"

	rpc "github.com/ctessum/cityaq/cityaqrpc"
	caqmock "github.com/ctessum/cityaq/cityaqrpc/mock_cityaqrpc"
	"github.com/golang/mock/gomock"
)

func TestDOM(t *testing.T) {
	doc := js.Global().Get("document")
	elem := doc.Call("createElement", "div")
	inputString := "hello world"
	elem.Set("innerText", inputString)
	out := elem.Get("innerText")

	// need Contains because a "\n" gets appended in the output
	if !strings.Contains(out.String(), inputString) {
		t.Errorf("unexpected output string. Expected %q to contain %q", out.String(), inputString)
	}
}

func TestCitySelector(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client := caqmock.NewMockCityAQClient(mockCtrl)

	client.EXPECT().Cities(
		gomock.Any(),
		gomock.AssignableToTypeOf(&rpc.CitiesRequest{}),
	).Return(&rpc.CitiesResponse{Names: []string{"city1", "city2"}}, nil)

	c := &CityAQ{
		CityAQClient: client,
		doc:          js.Global().Get("document"),
	}
	c.mapDiv = c.doc.Call("createElement", "div")
	c.loadMap()
	c.citySelector = c.doc.Call("createElement", "select")

	c.updateCitySelector(context.Background())
	html := c.citySelector.Get("innerHTML").String()
	want := `<option disabled="" hidden="">-- select an option --</option><option value="city1">city1</option><option value="city2">city2</option>`
	if html != want {
		t.Errorf("%v != %v", html, want)
	}

	// Call again to make sure contents get cleared every time.
	client.EXPECT().Cities(
		gomock.Any(), // expect any value for first parameter
		gomock.Any(), // expect any value for second parameter
	).Return(&rpc.CitiesResponse{Names: []string{"city3", "city4"}}, nil)

	c.updateCitySelector(context.Background())
	html = c.citySelector.Get("innerHTML").String()
	want = `<option disabled="" hidden="">-- select an option --</option><option value="city3">city3</option><option value="city4">city4</option>`
	if html != want {
		t.Errorf("%v != %v", html, want)
	}
}

func TestImpactTypeSelector(t *testing.T) {
	c := &CityAQ{
		doc: js.Global().Get("document"),
	}
	c.impactTypeSelector = c.doc.Call("createElement", "select")

	c.updateImpactTypeSelector()
	html := c.impactTypeSelector.Get("innerHTML").String()
	want := `<option disabled="" hidden="">-- select an option --</option><option value="1">Emissions</option>`
	if html != want {
		t.Errorf("%v != %v", html, want)
	}
}

func TestSourceTypeSelector(t *testing.T) {
	c := &CityAQ{
		doc: js.Global().Get("document"),
	}
	c.sourceTypeSelector = c.doc.Call("createElement", "select")

	c.updateSourceTypeSelector()
	html := c.sourceTypeSelector.Get("innerHTML").String()
	want := `<option disabled="" hidden="">-- select an option --</option><option value="electric_gen">electric_gen</option><option value="residential">residential</option><option value="commercial">commercial</option><option value="industrial">industrial</option><option value="builtup">builtup</option><option value="roadways">roadways</option><option value="railways">railways</option><option value="waterways">waterways</option><option value="busways">busways</option><option value="airports">airports</option><option value="agricultural">agricultural</option>`
	if html != want {
		t.Errorf("%v != %v", html, want)
	}
}

func TestEmissionSelector(t *testing.T) {
	c := &CityAQ{
		doc: js.Global().Get("document"),
	}
	c.emissionSelector = c.doc.Call("createElement", "select")

	c.updateEmissionSelector()
	html := c.emissionSelector.Get("innerHTML").String()
	want := `<option disabled="" hidden="">-- select an option --</option><option value="1">PM2_5</option><option value="2">NH3</option><option value="3">NOx</option><option value="4">SOx</option><option value="5">VOC</option>`
	if html != want {
		t.Errorf("%v != %v", html, want)
	}
}

func changeSelector(selector js.Value, index int) {
	selector.Set("selectedIndex", index)
}

func TestSelectors(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client := caqmock.NewMockCityAQClient(mockCtrl)

	client.EXPECT().Cities(
		gomock.Any(), // expect any value for first parameter
		gomock.Any(), // expect any value for second parameter
	).Return(&rpc.CitiesResponse{Names: []string{"city1", "city2"}}, nil)

	c := &CityAQ{
		CityAQClient: client,
		doc:          js.Global().Get("document"),
	}
	c.mapDiv = c.doc.Call("createElement", "div")
	c.loadMap()
	c.citySelector = c.doc.Call("createElement", "select")
	c.impactTypeSelector = c.doc.Call("createElement", "select")
	c.emissionSelector = c.doc.Call("createElement", "select")
	c.sourceTypeSelector = c.doc.Call("createElement", "select")

	if err := c.updateSelectors(context.Background()); err != nil {
		t.Fatal(err)
	}

	changeSelector(c.citySelector, 1)
	changeSelector(c.impactTypeSelector, 1)
	changeSelector(c.emissionSelector, 1)
	changeSelector(c.sourceTypeSelector, 1)
	changeSelector(c.sourceTypeSelector, 1)

	sel, err := c.selectorValues()
	if err != nil {
		t.Fatal(err)
	}
	want := &selections{cityName: "city1", impactType: rpc.ImpactType_Emissions, emission: 1, sourceType: "electric_gen"}

	if !reflect.DeepEqual(want, sel) {
		t.Errorf("%v != %v", sel, want)
	}
}
