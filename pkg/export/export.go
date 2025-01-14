package export

import (
	"regexp"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// pipelineExportRemoveExtraFields some of the fields are not getting skipped, when set as nil,
// like creationTimeStamp since it's a time.Time that is not a pointer yaml would display a creationTimeStamp: null
// let's do some regexp magic here, which should hopefully safe.
func pipelineExportRemoveExtraFields(input string) string {
	m1 := regexp.MustCompile(`\n  creationTimestamp: null\n`)
	return m1.ReplaceAllString(input, "\n")
}

// PipelineToYaml convert a pipeline to yaml removing extraneous fields generated by k8/tekton controllers so we can
// easily reimport theme
func PipelineToYaml(p *v1beta1.Pipeline) (string, error) {
	p.ObjectMeta.ManagedFields = nil
	p.ObjectMeta.ResourceVersion = ""
	p.ObjectMeta.UID = ""
	p.ObjectMeta.Generation = 0
	p.ObjectMeta.Namespace = ""
	p.ObjectMeta.CreationTimestamp = metav1.Time{}

	data, err := yaml.Marshal(p)
	return pipelineExportRemoveExtraFields(string(data)), err
}
