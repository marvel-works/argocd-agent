/*
Copyright 2019 The Codefresh Authors.

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

package main

import (
	//"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

/*
for usage in
\\go:generate go run generate generate_template.go <folder name under templates>
reads all files in folder and appends them to template map
*/

var outfileBaseName = "kubeobj.go"
var functionsMap map[string]string = map[string]string{
	"v1.Secret":                "CoreV1().Secrets(namespace)",
	"v1.ConfigMap":             "CoreV1().ConfigMaps(namespace)",
	"v1.Service":               "CoreV1().Services(namespace)",
	"v1.Pod":                   "CoreV1().Pods(namespace)",
	"v1.ServiceAccount":        "CoreV1().ServiceAccounts(namespace)",
	"v1.PersistentVolumeClaim": "CoreV1().PersistentVolumeClaims(namespace)",
	"v1.PersistentVolume":      "CoreV1().PersistentVolumes()",

	"v1beta1.Deployment": "ExtensionsV1beta1().Deployments(namespace)",
	"v1beta1.DaemonSet":  "ExtensionsV1beta1().DaemonSets(namespace)",

	"appsv1.Deployment": "AppsV1().Deployments(namespace)",
	"appsv1.DaemonSet":  "AppsV1().DaemonSets(namespace)",

	"rbacv1beta1.ClusterRole":        "RbacV1beta1().ClusterRoles()",
	"rbacv1beta1.ClusterRoleBinding": "RbacV1beta1().ClusterRoleBindings()",
	"rbacv1beta1.Role":               "RbacV1beta1().Roles(namespace)",
	"rbacv1beta1.RoleBinding":        "RbacV1beta1().RoleBindings(namespace)",

	"rbacv1.ClusterRole":        "RbacV1().ClusterRoles()",
	"rbacv1.ClusterRoleBinding": "RbacV1().ClusterRoleBindings()",
	"rbacv1.Role":               "RbacV1().Roles(namespace)",
	"rbacv1.RoleBinding":        "RbacV1().RoleBindings(namespace)",

	"storagev1.StorageClass": "StorageV1().StorageClasses()",

	"batchv1.Job":          "BatchV1().Jobs(namespace)",
	"batchv1beta1.CronJob": "BatchV1beta1().CronJobs(namespace)",
}

var packageTemplate = template.Must(template.New("").Parse(
	`
// Code generated by go generate; DO NOT EDIT.

package kubeobj

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/api/core/v1"
    v1beta1 "k8s.io/api/extensions/v1beta1"
    appsv1 "k8s.io/api/apps/v1"

    rbacv1beta1 "k8s.io/api/rbac/v1beta1"
    rbacv1 "k8s.io/api/rbac/v1"

	storagev1 "k8s.io/api/storage/v1"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

// CreateObject - creates kubernetes object from *runtime.Object. Returns object name, kind and creation error
func CreateObject(clientset *kubernetes.Clientset, obj runtime.Object, namespace string) (string, string, error){
	
	var name, kind string
	var err error
	switch objT := obj.(type) {
    {{ range $key, $value := .FunctionsMap }}
    case *{{ $key }}:
        name = objT.ObjectMeta.Name
        kind = objT.TypeMeta.Kind
        _, err = clientset.{{ $value }}.Create(objT)
    {{ end }}
    default:
        return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
    }
    return name, kind, err
}

// CheckObject - checks kubernetes object from *runtime.Object. Returns object name, kind and creation error
func CheckObject(clientset *kubernetes.Clientset, obj runtime.Object, namespace string) (string, string, error){
	
	var name, kind string
	var err error
	switch objT := obj.(type) {
    {{ range $key, $value := .FunctionsMap }}
    case *{{ $key }}:
        name = objT.ObjectMeta.Name
        kind = objT.TypeMeta.Kind
        _, err = clientset.{{ $value }}.Get(name, metav1.GetOptions{})
    {{ end }}
    default:
        return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
    }
    return name, kind, err
}

// DeleteObject - checks kubernetes object from *runtime.Object. Returns object name, kind and creation error
func DeleteObject(clientset *kubernetes.Clientset, obj runtime.Object, namespace string) (string, string, error){
	var propagationPolicy metav1.DeletionPropagation = "Background"
	var name, kind string
	var err error
	switch objT := obj.(type) {
    {{ range $key, $value := .FunctionsMap }}
    case *{{ $key }}:
        name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
        err = clientset.{{ $value }}.Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})
    {{ end }}
    default:
        return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
    }
    return name, kind, err
}

// ReplaceObject - replaces kubernetes object from *runtime.Object. Returns object name, kind and creation error
func ReplaceObject(clientset *kubernetes.Clientset, obj runtime.Object, namespace string) (string, string, error){
	var name, kind string
	var err error
	switch objT := obj.(type) {
    {{ range $key, $value := .FunctionsMap }}
    case *{{ $key }}:
        name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
        _, err = clientset.{{ $value }}.Update(objT)
    {{ end }}
    default:
        return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
    }
    return name, kind, err
}

`))

type tempateData struct {
	FunctionsMap map[string]string
}

func main() {

	var currentFilePath string
	if strings.Contains(os.Args[0], "/go-build") {
		_, currentFilePath, _, _ = runtime.Caller(0)
	} else {
		currentFilePath = os.Args[0]
	}

	currentDir := filepath.Dir(currentFilePath)
	var folderName = path.Join(currentDir, "kubeobj")

	outfileName := path.Join(folderName, "kubeobj.go")
	outfile, err := os.Create(outfileName)
	if err != nil {
		fmt.Printf("ERROR: cannot create out file %v", err)
		os.Exit(1)
	}
	defer outfile.Close()

	err = packageTemplate.Execute(outfile, tempateData{
		FunctionsMap: functionsMap,
	})
	if err != nil {
		fmt.Printf("generate_template ERROR: cannot generate template %v \n", err)
		os.Exit(1)
	}
}
