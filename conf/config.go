package conf

import "os"

var KubernetesDisable = os.Getenv("PD_KUBERNETES_DISABLE") == "true"
