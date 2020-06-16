package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	// "k8s.io/apimachinery/pkg/api/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string
var letters = []rune("abcde")

type Job struct {
	client         *kubernetes.Clientset
	name           string
	namespace      string
	containerImage string
}

func (j *Job) Create(args []string) error {
	containers := []corev1.Container{
		corev1.Container{
			Name:  j.name,
			Image: j.containerImage,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:   constructEnvs(),
			Args:  args,
		},
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: j.name,
			Namespace:    j.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers:    containers,
					RestartPolicy: "Never",
				},
			},
		},
	}
	_, err := j.client.BatchV1().Jobs(j.namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	return err
}

func (j *Job) Get() (*batchv1.Job, error) {
	job, err := j.client.BatchV1().Jobs(j.namespace).Get(context.TODO(), j.name, metav1.GetOptions{})
	return job, err
}

func constructEnvs() []corev1.EnvVar {
	return []corev1.EnvVar{
		corev1.EnvVar{
			Name: "MYSQL_PASS",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: os.Getenv("PASS_SECRET"),
					},
					Key: os.Getenv("SECRET_KEY"),
				},
			},
		},
		corev1.EnvVar{
			Name:  "MYSQL_USER",
			Value: os.Getenv("MYSQL_USER"),
		},
		corev1.EnvVar{
			Name:  "MYSQL_HOST",
			Value: os.Getenv("MYSQL_HOST"),
		},
		corev1.EnvVar{
			Name:  "MYSQL_PORT",
			Value: os.Getenv("MYSQL_PORT"),
		},
		corev1.EnvVar{
			Name:  "MYSQL_DB",
			Value: os.Getenv("MYSQL_DB"),
		},
	}
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file")
	flag.Parse()
	// Build config with kubeconfig file if provided, else build using in-cluster config
	var err error
	var config *rest.Config
	if kubeconfig != nil && *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Fatalf("Error creating kubernetes client config. Error: %s", err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating kubernetes client from config. Error: %s", err.Error())
	}
	job := Job{client: clientset, name: "palindrome", namespace: "default", containerImage: os.Getenv("JOB_DOCKER_IMAGE")}
	rand.Seed(time.Now().UnixNano())
	randStr := randSeq(20)
	log.Printf("Generated random string %s", randStr)
	err = job.Create([]string{"-inputString", randStr})
	if err != nil {
		log.Fatalf("Error creating job object. Error: %s", err.Error())
	}
}
