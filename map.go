package cityaq

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	rpc "github.com/ctessum/cityaq/cityaqrpc"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/simplify"

	"github.com/ctessum/requestcache/v4"
)

// MapTileServer serves map tiles for visualizing simulation results.
type MapTileServer struct {
	c     *CityAQ
	cache *requestcache.Cache
}

// NewMapTileServer creates a new map tile server,
// where cacheSize specifies the number of map layers
// to hold in an in-memory cache.
func NewMapTileServer(c *CityAQ, cacheSize int) *MapTileServer {
	s := &MapTileServer{c: c}
	s.cache = requestcache.NewCache(requestcache.Deduplicate(), requestcache.Memory(cacheSize))
	return s
}

func (s *MapTileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mapSpec, x, y, z, err := parseMapRequest(r.URL)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 404)
		return
	}
	layers, err := s.Layers(r.Context(), mapSpec)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	layers.ProjectToTile(maptile.New(uint32(x), uint32(y), maptile.Zoom(z)))
	layers.Clip(mvt.MapboxGLDefaultExtentBound)
	layers.Simplify(simplify.DouglasPeucker(1.0))
	layers.RemoveEmpty(1.0, 2.0)

	var data []byte
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		data, err = mvt.Marshal(layers)
	} else {
		w.Header().Set("Content-Encoding", "gzip")
		data, err = mvt.MarshalGzipped(layers)
	}
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	if _, err = w.Write(data); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

// MapSpecification holds information to specify the map data to serve.
type MapSpecification struct {
	CityName       string
	ImpactType     rpc.ImpactType
	Emission       rpc.Emission
	SourceType     string
	SimulationType rpc.SimulationType
	s              *MapTileServer
}

// Key returns a unique identifier for the receiver.
func (ms *MapSpecification) Key() string {
	return fmt.Sprintf("%s_%d_%d_%s_%d", ms.CityName, ms.ImpactType, ms.Emission, ms.SourceType, ms.SimulationType)
}

func queryString(u *url.URL, q url.Values, k string) (string, error) {
	v := q.Get(k)
	if v == "" {
		return "", fmt.Errorf("map request %s missing %s", k, u.Path)
	}
	return html.UnescapeString(v), nil
}

func queryInt(u *url.URL, q url.Values, k string) (int, error) {
	s, err := queryString(u, q, k)
	if err != nil {
		return -1, err
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("map request invalid value for %s: %s", k, s)
	}
	return int(i), nil
}

// parseRequest parses a request of the type
// xxx?x={x}&y={y}&z={z}&c={city}&it={ImpactType}&em={Emission}&st={SourceType}&sit={SimulationType}
func parseMapRequest(u *url.URL) (*MapSpecification, int, int, int, error) {
	q := u.Query()
	ms := new(MapSpecification)
	x, err := queryInt(u, q, "x")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	y, err := queryInt(u, q, "y")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	z, err := queryInt(u, q, "z")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	ms.CityName, err = queryString(u, q, "c")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	i, err := queryInt(u, q, "it")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	ms.ImpactType = rpc.ImpactType(i)

	i, err = queryInt(u, q, "em")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	ms.Emission = rpc.Emission(i)

	ms.SourceType, err = queryString(u, q, "st")
	if err != nil {
		return nil, -1, -1, -1, err
	}

	i, err = queryInt(u, q, "sit")
	if err != nil {
		return nil, -1, -1, -1, err
	}
	ms.SimulationType = rpc.SimulationType(i)

	return ms, x, y, z, nil
}

type layersResponse struct {
	L mvt.Layers
}

func (l *layersResponse) MarshalBinary() (data []byte, err error) {
	return mvt.Marshal(l.L)
}

func (l *layersResponse) UnmarshalBinary(data []byte) error {
	var err error
	l.L, err = mvt.Unmarshal(data)
	return err
}

// Layers returns the vector tile layers associated with ms.
func (s *MapTileServer) Layers(ctx context.Context, ms *MapSpecification) (mvt.Layers, error) {
	ms.s = s
	var layers layersResponse
	err := s.cache.NewRequest(ctx, ms).Result(&layers)
	if err != nil {
		return nil, err
	}
	return cloneLayers(layers.L), nil
}

func mapResolution(sourceType, cityName string) float64 {
	if egugridEmissions(sourceType) {
		return 0.1
	}
	// TODO: Revert so that all cities have the same resolution.
	if cityName == "Guadalajara" || cityName == "Tokyo" || cityName == "Melbourne" {
		return 0.005
	}
	return 0.002
}

func (ms *MapSpecification) Run(ctx context.Context, r requestcache.Result) error {
	var dataLayer *mvt.Layer
	switch ms.ImpactType {
	case rpc.ImpactType_Emissions:
		req := &rpc.GriddedEmissionsRequest{
			CityName:       ms.CityName,
			Emission:       ms.Emission,
			SourceType:     ms.SourceType,
			SimulationType: ms.SimulationType,
		}
		var err error
		dataLayer, err = ms.s.c.emissionsMapData(ctx, req)
		if err != nil {
			return err
		}
	case rpc.ImpactType_Concentrations:
		req := &rpc.GriddedConcentrationsRequest{
			CityName:       ms.CityName,
			Emission:       ms.Emission,
			SourceType:     ms.SourceType,
			SimulationType: ms.SimulationType,
		}
		var err error
		dataLayer, err = ms.s.c.concentrationsMapData(ctx, req)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid impact type %s", ms.ImpactType.String())
	}

	cityGeom, err := ms.s.c.geojsonGeometry(ms.CityName)
	if err != nil {
		return err
	}
	cityLayerData := geojson.NewFeatureCollection()
	feature := geojson.NewFeature(geomToOrb(cityGeom))
	feature.ID = uint64(0)
	cityLayerData = cityLayerData.Append(feature)
	cityLayer := mvt.NewLayer(ms.CityName, cityLayerData)

	o := r.(*layersResponse)
	o.L = mvt.Layers{dataLayer, cityLayer}

	if egugridEmissions(ms.SourceType) {
		egugridGeom, err := ms.s.c.countryOrGridBuffer(ms.CityName)
		if err != nil {
			return err
		}
		egugridLayerData := geojson.NewFeatureCollection()
		feature := geojson.NewFeature(geomToOrb(egugridGeom.Polygon))
		feature.ID = uint64(0)
		egugridLayerData = egugridLayerData.Append(feature)
		o.L = append(o.L, mvt.NewLayer(ms.CityName+"_egugrid", egugridLayerData))
	}
	return nil
}

func cloneLayers(layers mvt.Layers) mvt.Layers {
	o := make(mvt.Layers, len(layers))
	for i, layer := range layers {
		o[i] = cloneLayer(layer)
	}
	return o
}

func cloneLayer(l *mvt.Layer) *mvt.Layer {
	o := &mvt.Layer{
		Name:     l.Name,
		Version:  l.Version,
		Extent:   l.Extent,
		Features: make([]*geojson.Feature, len(l.Features)),
	}
	for i, f := range l.Features {
		of := &geojson.Feature{
			ID:         f.ID,
			Type:       f.Type,
			BBox:       f.BBox,
			Geometry:   orb.Clone(f.Geometry),
			Properties: f.Properties,
		}
		o.Features[i] = of
	}
	return o
}
