package model

//DependencyTree object model for building tree
type DependencyTree struct {
	Name         string            `json:"name,omitempty"`
	Version      string            `json:"version,omitempty"`
	Dependencies []*DependencyTree `json:"dependencies,omitempty"`
}

//NpmDependency object model as return from registry
type NpmDependency struct {
	Name         string            `json:"name,omitempty"`
	Version      string            `json:"version,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}
