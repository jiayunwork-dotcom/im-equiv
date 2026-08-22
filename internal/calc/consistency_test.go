// Package calc: consistency-check tests.
package calc

import (
	"strings"
	"testing"
)

func TestCheckZeroSlipTorqueAcceptsZero(t *testing.T) {
	if err := CheckZeroSlipTorque(0, 0); err != nil {
		t.Errorf("CheckZeroSlipTorque(0, 0): unexpected error: %v", err)
	}
}

func TestCheckZeroSlipTorqueRejectsNonZero(t *testing.T) {
	// A zero-slip solution carrying torque is the signature of the slip being
	// applied in the wrong place and must be rejected.
	err := CheckZeroSlipTorque(0, 42.5)
	if err == nil {
		t.Fatal("CheckZeroSlipTorque(0, 42.5): expected error, got none")
	}
	if !strings.Contains(err.Error(), "zero slip") {
		t.Errorf("error = %q, want mention of zero slip", err)
	}
}

func TestCheckZeroSlipTorqueIgnoresNonZeroSlip(t *testing.T) {
	if err := CheckZeroSlipTorque(0.03, 13.46); err != nil {
		t.Errorf("CheckZeroSlipTorque(0.03, 13.46): unexpected error: %v", err)
	}
}

func TestCheckEnergyBalanceAcceptsSplit(t *testing.T) {
	// Pag = 100 W splits into Pcu2 = 30 W and Pem = 70 W.
	if err := CheckEnergyBalance(100, 30, 70); err != nil {
		t.Errorf("CheckEnergyBalance(100, 30, 70): unexpected error: %v", err)
	}
}

func TestCheckEnergyBalanceRejectsDrift(t *testing.T) {
	// Pag = 100 W but Pcu2 + Pem = 110 W: the split is broken.
	err := CheckEnergyBalance(100, 40, 70)
	if err == nil {
		t.Fatal("CheckEnergyBalance(100, 40, 70): expected error, got none")
	}
	if !strings.Contains(err.Error(), "balance") {
		t.Errorf("error = %q, want mention of the energy balance", err)
	}
}

func TestCheckVoltageScalingAcceptsQuad(t *testing.T) {
	if err := CheckVoltageScaling(10, 40, 2); err != nil {
		t.Errorf("CheckVoltageScaling(10, 40, 2): unexpected error: %v", err)
	}
}

func TestCheckVoltageScalingRejectsWrongQuad(t *testing.T) {
	if err := CheckVoltageScaling(10, 30, 2); err == nil {
		t.Error("CheckVoltageScaling(10, 30, 2): expected error, got none")
	}
}
