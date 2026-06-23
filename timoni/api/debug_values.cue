@if(debug)

package main

// Значения, используемые debug_tool.cue.
// Пример отладки 'cue cmd -t debug -t name=test -t namespace=test -t mv=1.0.0 -t kv=1.28.0 build'.
values: {
	podAnnotations: "cluster-autoscaler.kubernetes.io/safe-to-evict": "true"
	message: "Hello Debug"
	image: {
		repository: "docker.io/nginx"
		tag:        "1-alpine"
		digest:     ""
	}
	test: {
		enabled: true
		image: {
			repository: "docker.io/curlimages/curl"
			tag:        "latest"
			digest:     ""
		}
	}
	affinity: nodeAffinity: requiredDuringSchedulingIgnoredDuringExecution: nodeSelectorTerms: [{
		matchExpressions: [{
			key:      "kubernetes.io/os"
			operator: "In"
			values: ["linux"]
		}]
	}]
}
