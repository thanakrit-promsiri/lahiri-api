package api

import (
	"lahiri-api/internal/calculator"
)

// --- Request Structs (Aliased from domain/calculator for 100% backward compatibility) ---

type CalculateReq = calculator.CalculateReq
type TransitReq = calculator.TransitReq

type CalculateWrapperReq struct {
	CalculateReq CalculateReq `json:"calculateReq"`
}

// TransitWrapperReq รองรับทั้ง JSON แบบ Flat ตามสเปกใหม่ และแบบห่อใน calculateReq / transitReq
type TransitWrapperReq struct {
	TransitReq
	CalculateReq *TransitReq `json:"calculateReq,omitempty"`
	WrappedReq   *TransitReq `json:"transitReq,omitempty"`
}

// --- Response Structs ---

type ErrorResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CalculateWrapperResp struct {
	Status        string        `json:"status"`
	CalculateResp CalculateResp `json:"calculateResp"`
}

type CalculateResp = calculator.CalculateResp
type ApiMetadata = calculator.ApiMetadata
type ApiLocation = calculator.ApiLocation
type ApiPlanet = calculator.ApiPlanet
type ApiHouse = calculator.ApiHouse
type ApiWholeSign = calculator.ApiWholeSign

type DashaWrapperResp struct {
	Status           string                      `json:"status"`
	VimshottariDasha calculator.VimshottariDasha `json:"vimshottari_dasha"`
}

type PromptWrapperResp struct {
	Status string `json:"status"`
	Prompt string `json:"prompt"`
}

// --- Transit Response Structs ---

type TransitWrapperResp struct {
	Status          string              `json:"status"`
	MonthlyTransits []ApiMonthlyTransit `json:"monthly_transits"`
}

type TransitMonthWrapperResp struct {
	Status         string            `json:"status"`
	MonthlyTransit ApiMonthlyTransit `json:"monthly_transit"`
}

type ApiMonthlyTransit = calculator.MonthlyTransit
type ApiTransitDasha = calculator.TransitDasha
type ApiTransitPlanet = calculator.TransitPlanet
