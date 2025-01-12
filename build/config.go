package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-linkerd/internal/config"

	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var LatestVersion string
var WorkloadPath string
var MeshModelPath string
var AllVersions []string
var CRDnamesURL map[string]string

const Component = "Linkerd"

var meshmodelmetadata = map[string]interface{}{
	"Primary Color":   "#0DA6E5",
	"Secondary Color": "#0DA6E5",
	"Shape":           "circle",
	"Logo URL":        "https://github.com/cncf/artwork/blob/master/projects/linkerd/icon/white/linkerd-icon-white.svg",
	"SVG_Color":       "<svg xmlns=\"http://www.w3.org/2000/svg\" role=\"img\" viewBox=\"-1.21 12.29 504.92 469.42\"><style>svg {enable-background:new 0 0 500 500}</style><style>.st2{fill:#2beda7}</style><linearGradient id=\"SVGID_1_\" x1=\"477.221\" x2=\"477.221\" y1=\"106.515\" y2=\"308.8\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_1_)\" d=\"M460.4 106.5v182.8l33.7 19.5V126z\"/><linearGradient id=\"SVGID_2_\" x1=\"25.459\" x2=\"25.459\" y1=\"106.52\" y2=\"308.812\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_2_)\" d=\"M8.6 308.8l33.7-19.5V106.5L8.6 126z\"/><path d=\"M173.8 307l155.1 89.5v-38.9l-145.2-83.8-9.9 5.7zm164 141.4l-164-94.7v38.9l43.8 25.3-52.7 30.4c-3.5 2-3.5 7.2 0 9.2l25.7 14.9 60.7-35 60.7 35 25.7-14.9c3.6-1.9 3.6-7.1.1-9.1z\" class=\"st2\"/><linearGradient id=\"SVGID_3_\" x1=\"477.221\" x2=\"477.221\" y1=\"196.062\" y2=\"382.938\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_3_)\" d=\"M460.4 215.5v162.1c0 4.1 4.4 6.7 8 4.6l23.1-13.3c1.6-.9 2.7-2.7 2.7-4.6V196.1l-33.8 19.4z\"/><linearGradient id=\"SVGID_4_\" x1=\"403.048\" x2=\"403.048\" y1=\"238.884\" y2=\"425.76\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_4_)\" d=\"M394.2 425l21.7-12.6c2.5-1.4 4-4.1 4-6.9V238.9l-33.7 19.5v162.1c0 4 4.4 6.6 8 4.5z\"/><linearGradient id=\"SVGID_5_\" x1=\"328.877\" x2=\"328.877\" y1=\"281.704\" y2=\"472.469\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_5_)\" d=\"M312 472.5l31.1-17.9c1.6-.9 2.7-2.7 2.7-4.6V281.7L312 301.2v171.3z\"/><linearGradient id=\"SVGID_6_\" x1=\"173.82\" x2=\"173.82\" y1=\"281.704\" y2=\"472.466\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_6_)\" d=\"M159.6 454.5l31.1 17.9V301.2L157 281.7v168.2c0 1.9 1 3.7 2.6 4.6z\"/><linearGradient id=\"SVGID_7_\" x1=\"99.649\" x2=\"99.649\" y1=\"238.883\" y2=\"425.76\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_7_)\" d=\"M86.8 412.5l21.7 12.6c3.5 2 8-.5 8-4.6V258.3l-33.7-19.5v166.7c0 2.9 1.5 5.6 4 7z\"/><linearGradient id=\"SVGID_8_\" x1=\"25.478\" x2=\"25.478\" y1=\"196.059\" y2=\"382.936\" gradientUnits=\"userSpaceOnUse\"><stop offset=\"0\" stop-color=\"#2beda7\"/><stop offset=\"1\" stop-color=\"#018afd\"/></linearGradient><path fill=\"url(#SVGID_8_)\" d=\"M12.6 369.7l21.7 12.6c3.5 2 8-.5 8-4.6V215.5L8.6 196.1v166.7c0 2.8 1.5 5.4 4 6.9z\"/><path d=\"M494.1 126l-33.7-19.5-60.7 35-40.4-23.3L412 87.8c3.5-2 3.5-7.2 0-9.2L390.2 66c-2.5-1.4-5.5-1.4-8 0l-56.7 32.7-40.4-23.3L337.8 45c3.5-2 3.5-7.2 0-9.2L316 23.2c-2.5-1.4-5.5-1.4-8 0l-56.7 32.7-56.7-32.7c-2.5-1.4-5.5-1.4-8 0l-21.8 12.6c-3.5 2-3.5 7.2 0 9.2l295.4 170.6 33.7-19.5-60.7-35 60.9-35.1zM112.5 66L90.8 78.6c-3.5 2-3.5 7.2 0 9.2l295.4 170.6 33.7-19.5L120.5 66c-2.5-1.4-5.5-1.4-8 0zM8.6 126l60.7 35-60.7 35.1 33.8 19.4 60.6-35 40.5 23.4-60.7 35 33.7 19.5 60.7-35.1 40.4 23.4-60.7 35 33.8 19.5 60.6-35.1 60.7 35.1 33.7-19.5L42.3 106.5z\" class=\"st2\"/></svg>",
	"SVG_White":       "<svg xmlns=\"http://www.w3.org/2000/svg\" role=\"img\" viewBox=\"-1.21 12.29 504.92 469.42\"><style>svg {enable-background:new 0 0 500 500}</style><style>.st0{fill:#fff}</style><path d=\"M460.4 106.5v182.8l33.7 19.5V126zM8.6 308.8l33.7-19.5V106.5L8.6 126zm165.2-1.8l155.1 89.5v-38.9l-145.2-83.8-9.9 5.7zm164 141.4l-164-94.7v38.9l43.8 25.3-52.7 30.4c-3.5 2-3.5 7.2 0 9.2l25.7 14.9 60.7-35 60.7 35 25.7-14.9c3.6-1.9 3.6-7.1.1-9.1z\" class=\"st0\"/><path d=\"M460.4 215.5v162.1c0 4.1 4.4 6.7 8 4.6l23.1-13.3c1.6-.9 2.7-2.7 2.7-4.6V196.1l-33.8 19.4zM394.2 425l21.7-12.6c2.5-1.4 4-4.1 4-6.9V238.9l-33.7 19.5v162.1c0 4 4.4 6.6 8 4.5zM312 472.5l31.1-17.9c1.6-.9 2.7-2.7 2.7-4.6V281.7L312 301.2v171.3zm-152.4-18l31.1 17.9V301.2L157 281.7v168.2c0 1.9 1 3.7 2.6 4.6zm-72.8-42l21.7 12.6c3.5 2 8-.5 8-4.6V258.3l-33.7-19.5v166.7c0 2.9 1.5 5.6 4 7zm-74.2-42.8l21.7 12.6c3.5 2 8-.5 8-4.6V215.5L8.6 196.1v166.7c0 2.8 1.5 5.4 4 6.9z\" class=\"st0\"/><path d=\"M494.1 126l-33.7-19.5-60.7 35-40.4-23.3L412 87.8c3.5-2 3.5-7.2 0-9.2L390.2 66c-2.5-1.4-5.5-1.4-8 0l-56.7 32.7-40.4-23.3L337.8 45c3.5-2 3.5-7.2 0-9.2L316 23.2c-2.5-1.4-5.5-1.4-8 0l-56.7 32.7-56.7-32.7c-2.5-1.4-5.5-1.4-8 0l-21.8 12.6c-3.5 2-3.5 7.2 0 9.2l295.4 170.6 33.7-19.5-60.7-35 60.9-35.1zM112.5 66L90.8 78.6c-3.5 2-3.5 7.2 0 9.2l295.4 170.6 33.7-19.5L120.5 66c-2.5-1.4-5.5-1.4-8 0zM8.6 126l60.7 35-60.7 35.1 33.8 19.4 60.6-35 40.5 23.4-60.7 35 33.7 19.5 60.7-35.1 40.4 23.4-60.7 35 33.8 19.5 60.6-35.1 60.7 35.1 33.7-19.5L42.3 106.5z\" class=\"st0\"/></svg>",
}

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category:    "Orchestration & Management",
	SubCategory: "Service Mesh",
	Metadata:    meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_LINKERD)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			manifests.RemoveHelmTemplatingFromCRD(&manifest)
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	vs, err := config.GetLatestReleaseNames(30)
	if len(vs) == 0 {
		fmt.Println("dynamic component generation failure: ", err.Error())
		return
	}
	for _, v := range vs {
		AllVersions = append(AllVersions, string(v))
	}
	LatestVersion = AllVersions[0]
	DefaultGenerationMethod = adapter.Manifests
	names, err := config.GetFileNames("linkerd", "linkerd2", "charts/linkerd-crds/templates/**")
	if err != nil {
		fmt.Println("dynamic component generation failure: ", err.Error())
		return
	}
	for n := range names {
		if !strings.HasSuffix(n, ".yaml") {
			delete(names, n)
		}
	}
	CRDnamesURL = names
}
