package k8simageadmissioncontroller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	logging "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var COMPRESSED_IMAGE_SIZE_LIMIT = 1_000_000_000

var codecs = serializer.NewCodecFactory(runtime.NewScheme())

func admissionReviewFromRequest(r *http.Request, deserializer runtime.Decoder) (*admissionv1.AdmissionReview, error) {
	// Validate that the incoming content type is correct.
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("expected application/json content-type")
	}

	// Get the body data, which will be the AdmissionReview
	// content for the request.
	var body []byte
	if r.Body != nil {
		requestData, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		body = requestData
	} else {
		logging.Errorf("No body found in AdmissionReview!")
	}
	// Decode the request body into
	admissionReviewRequest := &admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, admissionReviewRequest); err != nil {
		return nil, err
	}

	return admissionReviewRequest, nil
}

func ValidatePod(w http.ResponseWriter, r *http.Request, compressedImageSize int64) {
	logging.Info("Received message on validate")
	deserializer := codecs.UniversalDeserializer()

	// Parse the AdmissionReview from the http request.
	admissionReviewRequest, err := admissionReviewFromRequest(r, deserializer)
	if err != nil {
		msg := fmt.Sprintf("Error getting admission review from request: %v", err)
		logging.Errorf(msg)
		w.WriteHeader(400)
		_, err := w.Write([]byte(msg))
		if err != nil {
			logging.Errorf("Error when writing bytes to response: %v", err)
		}
		return
	}

	// Do server-side validation that we are only dealing with a pod resource. This
	// should also be part of the ValidatingWebhookConfiguration in the cluster, but
	// we should verify here before continuing.
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admissionReviewRequest.Request.Resource != podResource {
		msg := fmt.Sprintf("Did not receive pod, got %s", admissionReviewRequest.Request.Resource.Resource)
		logging.Info(msg)
		w.WriteHeader(400)
		_, err := w.Write([]byte(msg))
		if err != nil {
			logging.Errorf("Error when writing bytes to response: %v", err)
		}
		return
	}

	// Decode the pod from the AdmissionReview.
	rawRequest := admissionReviewRequest.Request.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &pod); err != nil {
		msg := fmt.Sprintf("Error decoding raw pod: %v", err)
		logging.Error(msg)
		w.WriteHeader(500)
		_, err := w.Write([]byte(msg))
		if err != nil {
			logging.Errorf("Error when writing bytes to response: %v", err)
		}
		return
	}

	// Create a response that either allows or rejects the pod creation
	// based off of the value of the hello label. Also, check to see if
	// we should supply a warning message even it is allowed.
	admissionResponse := &admissionv1.AdmissionResponse{}
	admissionResponse.Allowed = true

	containers := pod.Spec.Containers
	for _, container := range containers {
		size, err := GetImageSize(container.Image)
		if err != nil {
			logging.Error("Cannot get image size")
		}
		if size > compressedImageSize {
			msg := "Container image is too big"
			admissionResponse.Allowed = false
			admissionResponse.Result = &metav1.Status{
				Message: msg,
			}
			logging.Errorf(msg)
			break
		}
	}

	// Construct the response, which is just another AdmissionReview.
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admissionResponse
	admissionReviewResponse.SetGroupVersionKind(admissionReviewRequest.GroupVersionKind())
	admissionReviewResponse.Response.UID = admissionReviewRequest.Request.UID

	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		msg := fmt.Sprintf("error marshalling response json: %v", err)
		logging.Info(msg)
		w.WriteHeader(500)
		_, err := w.Write([]byte(msg))
		if err != nil {
			logging.Errorf("Error when writing bytes to response: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	_, err = w.Write(resp)
	if err != nil {
		logging.Errorf("Error when writing bytes to response: %v", err)
	}

}
