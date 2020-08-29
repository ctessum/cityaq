package cityaq

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	rpc "github.com/ctessum/cityaq/cityaqrpc"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/spatialmodel/inmap/emissions/aep/aeputil"
)

func TestParseMapRequest(t *testing.T) {
	ms := &MapSpecification{
		CityName:   "Accra Metropolitan",
		Emission:   rpc.Emission_PM2_5,
		ImpactType: rpc.ImpactType_Emissions,
		SourceType: "roadways",
	}
	u, err := url.Parse(fmt.Sprintf("https://example.com/maptile?x=10&y=11&z=12&c=%s&it=%d&em=%d&st=%s&sit=%d",
		html.EscapeString(ms.CityName), ms.ImpactType, ms.Emission, ms.SourceType, ms.SimulationType))
	if err != nil {
		t.Fatal(err)
	}
	newMS, x, y, z, err := parseMapRequest(u)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(newMS, ms) {
		t.Errorf("map spec; %+v != %+v", newMS, ms)
	}
	if x != 10 {
		t.Errorf("x: %d != %d", x, 10)
	}
	if y != 11 {
		t.Errorf("y: %d != %d", y, 11)
	}
	if z != 12 {
		t.Errorf("z: %d != %d", z, 12)
	}
}

func TestMapTileServer_ServeHTTP(t *testing.T) {
	c := &CityAQ{
		CityGeomDir: "testdata/cities",
		SpatialConfig: aeputil.SpatialConfig{
			SrgSpecOSM:    "testdata/srgspec_osm.json",
			SCCExactMatch: true,
			GridRef:       []string{"testdata/gridref.txt"},
			OutputSR:      "+proj=longlat",
			InputSR:       "+proj=longlat",
		},
	}
	s := NewMapTileServer(c, 1)
	ms := &MapSpecification{
		CityName:   "Accra Metropolitan",
		Emission:   rpc.Emission_PM2_5,
		ImpactType: rpc.ImpactType_Emissions,
		SourceType: "roadways",
	}
	u := fmt.Sprintf("https://example.com/maptile?x=4090&y=3967&z=13&c=%s&it=%d&em=%d&st=%s&sit=%d",
		html.EscapeString(ms.CityName), ms.ImpactType, ms.Emission, ms.SourceType, ms.SimulationType)

	t.Run("no_compression", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", u, nil)
		if err != nil {
			t.Fatal(err)
		}
		s.ServeHTTP(w, r)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status %d; message: %s", resp.StatusCode, string(body))
		}
		ct := resp.Header.Get("Content-Type")
		ctWant := "application/octet-stream"
		if ct != ctWant {
			t.Fatalf("content type is %s but should be %s", ct, ctWant)
		}

		layers, err := mvt.Unmarshal(body)
		if err != nil {
			t.Fatal(err)
		}

		if len(layers) != 2 {
			t.Fatalf("wrong number of layers %d", len(layers))
		}

		if layers[0].Name != "Accra Metropolitan_1_1_roadways_0" {
			t.Errorf("wrong layer name %s", layers[0].Name)
		}

		var vSum float64
		for _, f := range layers[0].Features {
			vSum += f.Properties["v"].(float64)
		}
		wantVSum := 453437.49810108123
		if !similar(vSum, wantVSum, 1.0e-8) {
			t.Errorf("value sum %g != %g", vSum, wantVSum)
		}

		if layers[1].Name != "Accra Metropolitan" {
			t.Errorf("wrong layer name %s", layers[1].Name)
		}
	})

	t.Run("gzip", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", u, nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Add("Accept-Encoding", "gzip")
		s.ServeHTTP(w, r)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d; message: %s", resp.StatusCode, string(body))
		}
		ct := resp.Header.Get("Content-Type")
		ctWant := "application/x-gzip"
		if ct != ctWant {
			t.Fatalf("content type is %s but should be %s", ct, ctWant)
		}

		layers, err := mvt.UnmarshalGzipped(body)
		if err != nil {
			t.Fatal(err)
		}

		if len(layers) != 2 {
			t.Fatalf("wrong number of layers %d", len(layers))
		}

		if layers[0].Name != "Accra Metropolitan_1_1_roadways_0" {
			t.Errorf("wrong layer name %s", layers[0].Name)
		}

		var vSum float64
		for _, f := range layers[0].Features {
			vSum += f.Properties["v"].(float64)
		}
		wantVSum := 453437.49810108123
		if !similar(vSum, wantVSum, 1.0e-8) {
			t.Errorf("value sum %g != %g", vSum, wantVSum)
		}

		if layers[1].Name != "Accra Metropolitan" {
			t.Errorf("wrong layer name %s", layers[1].Name)
		}
	})
}
