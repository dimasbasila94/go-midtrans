package midtrans

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// IrisGateway struct
type IrisGateway struct {
	Client Client
}

// Call : base method to call IRIS API
func (gateway *IrisGateway) Call(method, path string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.APIEnvType.IrisURL() + path
	return gateway.Client.Call(method, path, body, v)
}

// Show list of supported banks in IRIS. (https://iris-docs.midtrans.com/#list-banks)
func (gateway *IrisGateway) GetListBeneficiaryBank() (IrisBeneficiaryBanksResponse, error) {
	resp := IrisBeneficiaryBanksResponse{}

	err := gateway.Call("GET", "api/v1/beneficiary_banks", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error getting beneficiary banks: ", err)
		return resp, err
	}

	return resp, nil
}

// Create Beneficiaries (https://iris-docs.midtrans.com/#create-beneficiaries)
func (gateway *IrisGateway) CreateBeneficiaries(req *IrisBeneficiaries) (bool, error) {
	resp := IrisBeneficiariesResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/beneficiaries", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error creating beneficiaries: ", err)
		return false, err
	}

	if resp.Status != "created" {
		gateway.Client.Logger.Println("Error creating beneficiaries: ", resp.Errors)
		return false, errors.New(strings.Join(resp.Errors, ","))
	}

	return true, nil
}

// Update Beneficiaries (https://iris-docs.midtrans.com/#update-beneficiaries)
func (gateway *IrisGateway) UpdateBeneficiaries(aliasName string, req *IrisBeneficiaries) (bool, error) {
	resp := IrisBeneficiariesResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("PATCH", fmt.Sprintf("api/v1/beneficiaries/%s", aliasName), bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error updating beneficiaries: ", err)
		return false, err
	}

	if resp.Status != "updated" {
		gateway.Client.Logger.Println("Error updating beneficiaries: ", resp.Errors)
		return false, errors.New(strings.Join(resp.Errors, ","))
	}

	return true, nil
}

// Get List Beneficiaries (https://iris-docs.midtrans.com/#list-beneficiaries)
func (gateway *IrisGateway) GetListBeneficiaries() ([]IrisBeneficiaries, error) {
	var resp []IrisBeneficiaries

	err := gateway.Call("GET", "api/v1/beneficiaries", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error get list beneficiaries: ", err)
		return resp, err
	}

	return resp, nil
}

// CreatePayouts : This API is for Creator to create a payout. It can be used for single payout and also multiple payouts. (https://iris-docs.midtrans.com/#create-payouts)
func (gateway *IrisGateway) CreatePayouts(req IrisCreatePayoutReq) (IrisCreatePayoutResponse, error) {
	resp := IrisCreatePayoutResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error creating payouts: ", err)
		return resp, err
	}

	if resp.ErrorMessage != "" {
		return resp, errors.New(resp.ErrorMessage)
	}

	return resp, nil
}

// ApprovePayouts : Use this API for Apporver to approve multiple payout request. (https://iris-docs.midtrans.com/#approve-payouts)
func (gateway *IrisGateway) ApprovePayouts(req IrisApprovePayoutReq) (IrisApprovePayoutResponse, error) {
	resp := IrisApprovePayoutResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts/approve", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving payouts: ", err)
		return resp, err
	}

	if len(resp.Errors) > 0 {
		return resp, errors.New(strings.Join(resp.Errors, ", "))
	}

	if resp.Status != "ok" {
		return resp, errors.New("Error approving payouts, status from API not OK")
	}

	return resp, nil
}
