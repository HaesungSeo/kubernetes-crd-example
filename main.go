package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
	clientV1alpha1 "github.com/martin-helmich/kubernetes-crd-example/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defconf := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeconfig := flag.String("k", defconf, "path to Kubernetes config file")
	ns := flag.String("n", "default", "name space")
	flag.Parse()

	// Bootstrap k8s configuration from local       Kubernetes config file
	log.Println("Using kubeconfig file: ", *kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Println(err)
		// Try it again
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Register CRD scheme
	v1alpha1.AddToScheme(scheme.Scheme)

	// Create an rest client not targeting specific API version
	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	projects, err := clientSet.Projects(*ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d projects found:\n", len(projects.Items))
	for _, p := range projects.Items {
		fmt.Printf("%+v\n", p)
	}

	// Set up a watch on the Kubernetes API for changes to crd project
	watcher, err := clientSet.Projects(*ns).Watch(context.Background(), metav1.ListOptions{
		ResourceVersion: projects.ResourceVersion,
	})
	if err != nil {
		log.Fatalln("failed to get watch channel:", err)
	}

	// Iterate the crd event
	for {
		select {
		case event := <-watcher.ResultChan():
			if event.Object == nil {
				// Once the channel broken, need to restart the process
				log.Fatalln("nil event")
			}
			project, ok := event.Object.(*v1alpha1.Project)
			if !ok {
				str := fmt.Sprintf("%s: unexpected event Object type: %T\n", event.Type, event.Object)
				log.Fatalln(str)
				return
			}

			// access crd object just like built-in object!
			log.Printf("%s: %s\n", event.Type, project.Name)
			log.Printf("  spec.replicas %d\n", project.Spec.Replicas)
		}
	}
}
