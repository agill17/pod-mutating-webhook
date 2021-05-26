package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

type MutateContainerImage struct {
	logger logr.Logger
}

func (m *MutateContainerImage) MutateContainerImages(reqBody []byte, imageRules map[string]string) ([]byte, error) {
	review := &v1.AdmissionReview{}
	resp := &v1.AdmissionResponse{}

	if errUnmarshalling := json.Unmarshal(reqBody, review); errUnmarshalling != nil {
		return nil, errUnmarshalling
	}

	podObj := &corev1.Pod{}
	if errUnmarshallingToPod := json.Unmarshal(review.Request.Object.Raw, podObj); errUnmarshallingToPod != nil {
		return nil, errUnmarshallingToPod
	}
	if podObj.GetDeletionTimestamp() != nil {
		return nil, errors.New("ErrPodIsInDeletingState")
	}
	var namespacedName string = fmt.Sprintf("%s/%s", podObj.GetNamespace(), podObj.GetName())

	jsonPatch := v1.PatchTypeJSONPatch
	resp.PatchType = &jsonPatch
	resp.Allowed = true
	resp.AuditAnnotations = map[string]string{"mutatedBy": "pod-mutating-webhook"}

	var patches []map[string]string

	for idx, cont := range podObj.Spec.Containers {
		if val, ok := imageRules[cont.Image]; ok && val != "" {
			p := map[string]string{
				"op":    "replace",
				"path":  fmt.Sprintf("/spec/containers/%d/image", idx),
				"value": val,
			}
			m.logger.Info(fmt.Sprintf("%v changing image from %v to %v", namespacedName, cont.Image, val))
			patches = append(patches, p)
		}

	}
	for idx, initCont := range podObj.Spec.InitContainers {
		if val, ok := imageRules[initCont.Image]; ok && val != "" {
			p := map[string]string{
				"op":    "replace",
				"path":  fmt.Sprintf("/spec/initContainers/%d/image", idx),
				"value": val,
			}
			m.logger.Info(fmt.Sprintf("%v changing image from %v to %v", namespacedName, initCont.Image, val))
			patches = append(patches, p)
		}
	}

	if len(patches) > 0 {
		respPatch, err := json.Marshal(patches)
		if err != nil {
			return nil, err
		}
		resp.Patch = respPatch
		review.Response = resp
		return json.Marshal(review)
	}
	return nil, nil
}
