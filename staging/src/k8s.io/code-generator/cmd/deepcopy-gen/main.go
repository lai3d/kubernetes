/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// deepcopy-gen is a tool for auto-generating DeepCopy functions.
//
// Given a list of input directories, it will generate DeepCopy and
// DeepCopyInto methods that efficiently perform a full deep-copy of each type.
// If these methods already exist (are predefined by the developer), they are
// used instead of generating new ones.  Generated code will use standard value
// assignment whenever possible.  If that is not possible it will try to call
// its own generated copy function for the type.  Failing that, it will fall
// back on `conversion.Cloner.DeepCopy(val)` to make the copy.  The resulting
// file will be stored in the same directory as the processed source package.
//
// If interfaces are referenced in types, it is expected that corresponding
// DeepCopyInterfaceName methods exist, e.g. DeepCopyObject for runtime.Object.
// These can be predefined by the developer or generated through tags, see
// below.  They must be added to the interfaces themselves manually, e.g.
//
//	type Object interface {
//	  ...
//	  DeepCopyObject() Object
//	}
//
// Generation is governed by comment tags in the source.  Any package may
// request DeepCopy generation by including a comment in the file-comments of
// one file, of the form:
//
//	// +k8s:deepcopy-gen=package
//
// DeepCopy functions can be generated for individual types, rather than the
// entire package by specifying a comment on the type definition of the form:
//
//	// +k8s:deepcopy-gen=true
//
// When generating for a whole package, individual types may opt out of
// DeepCopy generation by specifying a comment on the type definition of the
// form:
//
//	// +k8s:deepcopy-gen=false
//
// Additional DeepCopyInterfaceName methods can be generated by specifying a
// comment on the type definition of the form:
//
//	// +k8s:deepcopy-gen:interfaces=k8s.io/kubernetes/runtime.Object,k8s.io/kubernetes/runtime.List
//
// This leads to the generation of DeepCopyObject and DeepCopyList with the given
// interfaces as return types. We say that the tagged type implements deepcopy for the
// interfaces.
//
// The deepcopy funcs for interfaces using "+k8s:deepcopy-gen:interfaces" use the pointer
// of the type as receiver. For those special cases where the non-pointer object should
// implement the interface, this can be done with:
//
//	// +k8s:deepcopy-gen:nonpointer-interfaces=true
package main

import (
	"flag"

	"github.com/spf13/pflag"
	generatorargs "k8s.io/code-generator/cmd/deepcopy-gen/args"
	"k8s.io/code-generator/cmd/deepcopy-gen/generators"
	"k8s.io/gengo/v2/args"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	genericArgs, customArgs := generatorargs.NewDefaults()

	genericArgs.AddFlags(pflag.CommandLine)
	customArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := generatorargs.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	// Run it.
	if err := genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.GetTargets,
		args.StdBuildTag,
		pflag.Args(),
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
