{
	"apiVersion": "v1",
	"kind": "Pod",
	"metadata": {
		"name": "nginx"
	},
	"spec": {
		"containers": [
			{
				"name": "nginx",
				"image": "nginx:1.7.9",
				"ports": [
					{
						"containerPort": 80
					}
				]
			}
		]
	}
}
