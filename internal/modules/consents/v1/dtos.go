package v1

type ConsentType string

const (
	ConsentTypeTracking   ConsentType = "tracking"
	ConsentTypeMarketing  ConsentType = "marketing"
	ConsentTypeFunctional ConsentType = "functional"
	ConsentTypeNone       ConsentType = "none"
)

func (c ConsentType) Valid() bool {
	switch c {
	case ConsentTypeTracking, ConsentTypeMarketing, ConsentTypeFunctional, ConsentTypeNone:
		return true
	}
	return false
}

type ConsentStatus string

const (
	ConsentStatusGranted ConsentStatus = "granted"
	ConsentStatusDenied  ConsentStatus = "denied"
	ConsentStatusPending ConsentStatus = "pending"
)

func (s ConsentStatus) Valid() bool {
	switch s {
	case ConsentStatusGranted, ConsentStatusDenied, ConsentStatusPending:
		return true
	}
	return false
}

type UpdateConsentRequest struct {
	ConsentType ConsentType   `json:"consent_type" binding:"required"`
	Status      ConsentStatus `json:"status" binding:"required"`
}

type UpdateConsentResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
