// Код сгенерирован timoni.
// Обратите внимание: этот файл обязателен и должен содержать
// схему значений и рабочий процесс timoni.

package main

import (
	"strconv"
	"strings"

	templates "timoni.sh/mongodb/templates"
)

// Определяем схему для значений, переданных пользователем.
// Во время выполнения Timoni подставляет переданные значения
// и проверяет их согласно схеме Config.
values: templates.#Config

// Определяем, как Timoni должен собирать, проверять
// и применять ресурсы Kubernetes.
timoni: {
	apiVersion: "v1alpha1"

	// Определяем экземпляр, который выводит ресурсы Kubernetes.
	// Во время выполнения Timoni собирает экземпляр и проверяет
	// итоговые ресурсы согласно их Kubernetes-схеме.
	instance: templates.#Instance & {
		// Значения, переданные пользователем, объединяются со
		// значениями по умолчанию во время выполнения Timoni.
		config: values
		// Эти значения подставляются Timoni во время выполнения.
		config: {
			metadata: {
				name:      string @tag(name)
				namespace: string @tag(namespace)
			}
			moduleVersion: string @tag(mv, var=moduleVersion)
			kubeVersion:   string @tag(kv, var=kubeVersion)
		}
	}

	// Требуем минимальную версию Kubernetes.
	kubeMinorVersion: int & >=20
	kubeMinorVersion: strconv.Atoi(strings.Split(instance.config.kubeVersion, ".")[1])

	// Передаем ресурсы Kubernetes, выведенные экземпляром,
	// в многошаговое применение Timoni.
	apply: app: [ for obj in instance.objects {obj}]

	// Условно запускаем тесты после установки или обновления.
	if instance.config.test.enabled {
		apply: test: [ for obj in instance.tests {obj}]
	}
}
