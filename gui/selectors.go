package gui

import (
	"context"
	"errors"
	"strconv"
	"syscall/js"

	rpc "github.com/ctessum/cityaq/cityaqrpc"
)

const defaultSelectorText = "-- select an option --"

func updateSelector(doc, selector js.Value, values []interface{}, text []string) {
	selector.Set("innerHTML", "")
	option := doc.Call("createElement", "option")
	option.Set("disabled", true)
	option.Set("selected", true)
	option.Set("hidden", true)
	option.Set("text", defaultSelectorText)
	selector.Call("appendChild", option)
	for i, value := range values {
		option := doc.Call("createElement", "option")
		option.Set("value", value)
		option.Set("text", text[i])
		selector.Call("appendChild", option)
	}
}

// updateCitySelector updates the options of cities.
func (c *CityAQ) updateCitySelector(ctx context.Context) {
	if c.citySelector.IsUndefined() {
		c.citySelector = c.doc.Call("getElementById", "citySelector")
	}
	cities, err := c.Cities(ctx, &rpc.CitiesRequest{})
	if err != nil {
		c.logError(err)
		return
	}
	names := make([]interface{}, len(cities.Names))
	for i, n := range cities.Names {
		names[i] = n
	}
	updateSelector(c.doc, c.citySelector, names, cities.Names)
}

// updateImpactTypeSelector updates the options of impacts.
func (c *CityAQ) updateImpactTypeSelector() {
	if c.impactTypeSelector.IsUndefined() {
		c.impactTypeSelector = c.doc.Call("getElementById", "impactTypeSelector")
	}
	updateSelector(c.doc, c.impactTypeSelector, []interface{}{1, 2}, []string{"Emissions", "Concentrations"})
}

// updateEmissionSelector updates the options of emissions available.
func (c *CityAQ) updateEmissionSelector() {
	if c.emissionSelector.IsUndefined() {
		c.emissionSelector = c.doc.Call("getElementById", "emissionSelector")
	}
	values := make([]interface{}, len(rpc.Emission_value)-1)
	text := make([]string, len(rpc.Emission_value)-1)
	for i := 1; i < len(rpc.Emission_value); i++ {
		n := rpc.Emission_name[int32(i)]
		values[i-1] = i
		text[i-1] = n
	}
	updateSelector(c.doc, c.emissionSelector, values, text)
}

// updateSourceTypeSelector updates the options of source types available.
func (c *CityAQ) updateSourceTypeSelector() error {
	if c.sourceTypeSelector.IsUndefined() {
		c.sourceTypeSelector = c.doc.Call("getElementById", "sourceTypeSelector")
	}
	simulationType, err := c.simulationTypeSelectorValue()
	if err != nil {
		simulationType = 0
	}
	if simulationType == rpc.SimulationType_CityMarginal || simulationType == 0 {
		updateSelector(c.doc, c.sourceTypeSelector,
			[]interface{}{
				"electric_gen_egugrid", "population", "residential",
				"commercial", "industrial", "builtup", "roadways", "roadways_motorway",
				"roadways_trunk", "roadways_primary", "roadways_secondary", "roadways_tertiary",
				"railways", "waterways", "bus_routes", "airports", "agricultural",
			},
			[]string{
				"electric_gen_egugrid", "population", "residential",
				"commercial", "industrial", "builtup", "roadways", "roadways_motorway",
				"roadways_trunk", "roadways_primary", "roadways_secondary", "roadways_tertiary",
				"railways", "waterways", "bus_routes", "airports", "agricultural",
			})
		return nil
	}
	sectors, err := c.EmissionsInventorySectors(context.Background(), &rpc.EmissionsInventorySectorsRequest{})
	if err != nil {
		return err
	}
	secI := make([]interface{}, len(sectors.Sectors))
	for i, sec := range sectors.Sectors {
		secI[i] = sec
	}
	updateSelector(c.doc, c.sourceTypeSelector, secI, sectors.Sectors)
	return nil
}

// updateSimulationTypeSelector updates the options of simulation types available.
func (c *CityAQ) updateSimulationTypeSelector() {
	if c.simulationTypeSelector.IsUndefined() {
		c.simulationTypeSelector = c.doc.Call("getElementById", "simulationTypeSelector")
	}
	updateSelector(c.doc, c.simulationTypeSelector,
		[]interface{}{
			int(rpc.SimulationType_CityMarginal), int(rpc.SimulationType_CityTotal), int(rpc.SimulationType_Total),
		},
		[]string{
			"City Marginal", "City Total", "Total",
		})
}

func (c *CityAQ) updateSelectors(ctx context.Context) error {
	c.updateCitySelector(ctx)
	c.updateImpactTypeSelector()
	c.updateEmissionSelector()
	c.updateSimulationTypeSelector()
	if err := c.updateSourceTypeSelector(); err != nil {
		return err
	}
	return nil
}

func selectorValue(selector js.Value) (string, error) {
	v := selector.Get("value").String()
	if v == defaultSelectorText {
		return v, incompleteSelectionError
	}
	return v, nil
}

func (c *CityAQ) citySelectorValue() (string, error) {
	return selectorValue(c.citySelector)
}

func (c *CityAQ) impactTypeSelectorValue() (rpc.ImpactType, error) {
	v, err := selectorValue(c.impactTypeSelector)
	if err != nil {
		return -1, err
	}
	vInt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1, err
	}
	return rpc.ImpactType(vInt), nil
}

func (c *CityAQ) emissionSelectorValue() (rpc.Emission, error) {
	v, err := selectorValue(c.emissionSelector)
	if err != nil {
		return -1, err
	}
	vInt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1, err
	}
	return rpc.Emission(vInt), nil
}

func (c *CityAQ) sourceTypeSelectorValue() (string, error) {
	return selectorValue(c.sourceTypeSelector)
}

func (c *CityAQ) simulationTypeSelectorValue() (rpc.SimulationType, error) {
	v, err := selectorValue(c.simulationTypeSelector)
	if err != nil {
		return -1, err
	}
	vInt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1, err
	}
	return rpc.SimulationType(vInt), nil
}

type selections struct {
	cityName       string
	impactType     rpc.ImpactType
	sourceType     string
	emission       rpc.Emission
	simulationType rpc.SimulationType
}

var incompleteSelectionError = errors.New("incomplete selection")

func (c *CityAQ) selectorValues() (s *selections, err error) {
	s = new(selections)
	s.cityName, err = c.citySelectorValue()
	if err != nil {
		return
	}

	s.emission, err = c.emissionSelectorValue()
	if err != nil {
		return
	}

	s.impactType, err = c.impactTypeSelectorValue()
	if err != nil {
		return
	}

	s.sourceType, err = c.sourceTypeSelectorValue()
	if err != nil {
		return
	}

	s.simulationType, err = c.simulationTypeSelectorValue()
	if err != nil {
		return
	}

	return
}
