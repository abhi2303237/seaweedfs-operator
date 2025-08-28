package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	monitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	seaweedv1 "github.com/seaweedfs/seaweedfs-operator/api/v1"
)

func (r *SeaweedReconciler) createVolumeServerServiceMonitor(m *seaweedv1.Seaweed) *monitorv1.ServiceMonitor {
	labels := labelsForVolumeServer(m.Name)

	dep := &monitorv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-volume",
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: monitorv1.ServiceMonitorSpec{
			Endpoints: []monitorv1.Endpoint{
				{
					Path: "/metrics",
					Port: "volume-metrics",
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	return dep
}

func (r *SeaweedReconciler) createVolumeServerTopologyServiceMonitor(m *seaweedv1.Seaweed, topologyName string, topologySpec *seaweedv1.VolumeTopologySpec) *monitorv1.ServiceMonitor {
	labels := labelsForVolumeServerTopology(m.Name, topologyName)

	dep := &monitorv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-volume-" + topologyName,
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: monitorv1.ServiceMonitorSpec{
			Endpoints: []monitorv1.Endpoint{
				{
					Path: "/metrics",
					Port: "volume-metrics",
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	return dep
}
