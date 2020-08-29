package cityaq

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	rpc "github.com/ctessum/cityaq/cityaqrpc"
	"github.com/spatialmodel/inmap/emissions/aep/aeputil"
	"gonum.org/v1/gonum/floats"
)

func TestCityAQ_GriddedConcentrations_cityMarginal(t *testing.T) {
	dir := fmt.Sprintf("temp_test_%d", time.Now().Unix())
	c := &CityAQ{
		CityGeomDir: "testdata/cities",
		SpatialConfig: aeputil.SpatialConfig{
			SrgSpecOSM:    "testdata/srgspec_osm.json",
			SCCExactMatch: true,
			GridRef:       []string{"testdata/gridref.txt"},
			OutputSR:      "+proj=longlat",
			InputSR:       "+proj=longlat",
		},
		CacheLoc:                    "file://" + dir,
		InMAPCityMarginalConfigFile: "testdata/inmap_config.toml",
	}
	os.Mkdir(dir, os.ModePerm)
	defer os.RemoveAll(dir)

	r := &rpc.GriddedConcentrationsRequest{
		CityName:       "Accra Metropolitan",
		Emission:       rpc.Emission_PM2_5,
		SourceType:     "roadways",
		SimulationType: rpc.SimulationType_CityMarginal,
	}

	conc, err := c.GriddedConcentrations(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	concSum := floats.Sum(conc.Concentrations)
	wantConcSum := 0.039921666164999994
	if !similar(concSum, wantConcSum, 1.0e-10) {
		t.Errorf("concentration sum: %g != %g", concSum, wantConcSum)
	}
}

func TestCityAQ_GriddedConcentrations_cityTotal(t *testing.T) {
	dir := fmt.Sprintf("temp_test_%d", time.Now().Unix())
	c := &CityAQ{
		CityGeomDir: "testdata/cities",
		SpatialConfig: aeputil.SpatialConfig{
			SrgSpecOSM:    "testdata/srgspec_osm.json",
			SCCExactMatch: true,
			GridRef:       []string{"testdata/gridref.txt"},
			OutputSR:      "+proj=longlat",
			InputSR:       "+proj=longlat",
		},
		CacheLoc:                 "file://" + dir,
		InMAPCityTotalConfigFile: "testdata/inmap_config_coards.toml",
	}
	os.Mkdir(dir, os.ModePerm)
	defer os.RemoveAll(dir)

	r := &rpc.GriddedConcentrationsRequest{
		CityName:       "Accra Metropolitan",
		Emission:       rpc.Emission_PM2_5,
		SourceType:     "all",
		SimulationType: rpc.SimulationType_CityTotal,
	}

	conc, err := c.GriddedConcentrations(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	concSum := floats.Sum(conc.Concentrations)
	wantConcSum := 0.0
	if !similar(concSum, wantConcSum, 1.0e-10) {
		t.Errorf("concentration sum: %g != %g", concSum, wantConcSum)
	}
}

func TestCityAQ_GriddedConcentrations_total(t *testing.T) {
	dir := fmt.Sprintf("temp_test_%d", time.Now().Unix())
	c := &CityAQ{
		CityGeomDir: "testdata/cities",
		SpatialConfig: aeputil.SpatialConfig{
			SrgSpecOSM:    "testdata/srgspec_osm.json",
			SCCExactMatch: true,
			GridRef:       []string{"testdata/gridref.txt"},
			OutputSR:      "+proj=longlat",
			InputSR:       "+proj=longlat",
		},
		CacheLoc:             "file://" + dir,
		InMAPTotalConfigFile: "testdata/inmap_config_coards.toml",
	}
	os.Mkdir(dir, os.ModePerm)
	defer os.RemoveAll(dir)

	r := &rpc.GriddedConcentrationsRequest{
		CityName:       "Accra Metropolitan",
		Emission:       rpc.Emission_PM2_5,
		SourceType:     "all",
		SimulationType: rpc.SimulationType_Total,
	}

	conc, err := c.GriddedConcentrations(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	concSum := floats.Sum(conc.Concentrations)
	wantConcSum := 0.0
	if !similar(concSum, wantConcSum, 1.0e-10) {
		t.Errorf("concentration sum: %g != %g", concSum, wantConcSum)
	}
}
