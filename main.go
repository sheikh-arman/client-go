package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/user/.kube/config", "path to the kubeconfig file")
	//flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Okay")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {

	}

	pods, err := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {

	}
	for _, pod := range pods.Items {
		fmt.Printf("%s \n", pod.Name)
	}
	//fmt.Println(pods)
}
