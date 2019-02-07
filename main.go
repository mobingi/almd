package main

import (
	goflag "flag"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var Version = "?"

var (
	str = "hello world"

	rootCmd = &cobra.Command{
		Use:   "oceand",
		Short: "k8s agent for ocean",
		Long:  "Mobingi Ocean agent for Kubernetes.",
	}
)

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.AddCommand(
		RunCmd(),
	)

	klog.InitFlags(nil)
	goflag.Set("logtostderr", "true")
	goflag.Parse()
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func testK8s() {
	// Creates the in-cluster config.
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Creates the clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		klog.Infof("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		//
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			klog.Errorf("Pod not found")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			klog.Errorf("Error getting pod %v", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			klog.Errorf("Found pod")
		}

		time.Sleep(10 * time.Second)
	}
}

func RunCmd() *cobra.Command {
	runcmd := &cobra.Command{
		Use:   "run",
		Short: "core run command",
		Long:  "Core run command.",
		Run: func(cmd *cobra.Command, args []string) {
			testK8s()
		},
	}

	runcmd.Flags().SortFlags = false
	runcmd.Flags().StringVar(&str, "str", str, "string to print")
	return runcmd
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		klog.Fatalf("root cmd execute failed, err=%v", err)
	}
}
