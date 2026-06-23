package templates

import (
	corev1 "k8s.io/api/core/v1"
	timoniv1 "timoni.sh/core/v1alpha1"
)

// Config определяет схему и значения по умолчанию для значений Instance.
#Config: {
	// Информация о версии выполнения
	moduleVersion!: string
	kubeVersion!:   string

	// Метаданные (общие для всех ресурсов)
	metadata: timoniv1.#Metadata & {#Version: moduleVersion}

	// Селектор меток (общий для всех ресурсов)
	selector: timoniv1.#Selector & {#Name: metadata.name}

	// Развертывание
	replicas: *1 | int & >0

	// Pod
	podAnnotations?: {[ string]: string}
	podSecurityContext?: corev1.#PodSecurityContext
	imagePullSecrets?: [...corev1.LocalObjectReference]
	tolerations?: [ ...corev1.#Toleration]
	affinity?: corev1.#Affinity
	topologySpreadConstraints?: [...corev1.#TopologySpreadConstraint]

	// Контейнер
	image!:           timoniv1.#Image
	imagePullPolicy:  *"IfNotPresent" | string
	resources?:       corev1.#ResourceRequirements
	securityContext?: corev1.#SecurityContext

	// Сервис
	service: port: *27017 | int & >0 & <=65535

	// Тестовая Job
	test: {
		enabled: *false | bool
	}
}

// Instance принимает значения конфигурации и выводит объекты Kubernetes.
#Instance: {
	config: #Config

	objects: {
		sa:     #ServiceAccount & {_config: config}
		svc:    #Service & {_config:        config}
		deploy: #Deployment & {_config:     config}
	}
}
