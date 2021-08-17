package parser

import (
	"fmt"

	"github.com/devfile/library/pkg/devfile/parser/json"
	"github.com/devfile/library/pkg/testingutil/filesystem"
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/yaml"
)

// WriteYamlDevfile creates a devfile.yaml file
func (d *DevfileObj) WriteYamlDevfile() error {

	// Encode data into YAML format
	yamlData, err := Marshal(d.Data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal devfile object into yaml")
	}
	// Write to devfile.yaml
	fs := d.Ctx.GetFs()
	if fs == nil {
		fs = filesystem.DefaultFs{}
	}
	err = fs.WriteFile(d.Ctx.GetAbsPath(), yamlData, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to create devfile yaml file")
	}
	fmt.Printf(string(yamlData))
	// Successful
	klog.V(2).Infof("devfile yaml created at: '%s'", OutputDevfileYamlPath)
	return nil
}

// Marshal marshals the object into JSON then converts JSON to YAML and returns the
// YAML.
func Marshal(o interface{}) ([]byte, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("error marshaling into JSON: %v", err)
	}
	fmt.Printf("JSON")
	fmt.Printf(string(j))
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return nil, fmt.Errorf("error converting JSON to YAML: %v", err)
	}

	return y, nil
}
func JSONToYAML(j []byte) ([]byte, error) {
	// Convert the JSON to an object.
	var jsonObj interface{}
	// We are using yaml.Unmarshal here (instead of json.Unmarshal) because the
	// Go JSON library doesn't try to pick the right number type (int, float,
	// etc.) when unmarshalling to interface{}, it just picks float64
	// universally. go-yaml does go through the effort of picking the right
	// number type, so we can preserve number type throughout this process.
	err := yaml.Unmarshal(j, &jsonObj)
	fmt.Printf("YAML to JSON")
	fmt.Print(jsonObj)
	if err != nil {
		return nil, err
	}

	// Marshal this object into YAML.
	return yaml.Marshal(jsonObj)
}

/*func (d *v2.OrderedDevfileV2) marshalJSON() ([]byte, error) {
    var b []byte
    buf:=bytes.NewBuffer(b)
    buf.WriteRune('{')
    l:=len(d.)
    for i,key:=range om.Order {
        km,err:=json.Marshal(key)
        if err!=nil { return nil,err }
        buf.Write(km)
        buf.WriteRune(':')
        vm,err:=json.Marshal(om.Map[key])
        if err!=nil { return nil,err }
        buf.Write(vm)
        if i!=l-1 { buf.WriteRune(',') }
        fmt.Println(buf.String())
    }
    buf.WriteRune('}')
    fmt.Println(buf.String())
    return buf.Bytes(),nil
}*/
