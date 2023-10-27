package ocpp201

type mutability string

const (
	ReadOnly  mutability = "r"
	WriteOnly mutability = "w"
	ReadWrite mutability = "rw"
)

type secEventSeq struct {
	isToTrigger bool
	eventType   securityEventType
	eventInfo   string
}

type securityEventType string

//nolint:varcheck,deadcode // Because this are a list that could be used in the future
const (
	firmwareUpdated                     securityEventType = "FirmwareUpdated"                     // CRITICAL The Charging Station firmware is updated
	failedToAuthenticateAtCsms          securityEventType = "FailedToAuthenticateAtCsms"          // The authentication credentials provided by the Charging Station were rejected by the CSMS
	csmsFailedToAuthenticate            securityEventType = "CsmsFailedToAuthenticate"            // The authentication credentials provided by the CSMS were rejected by the Charging Station
	settingSystemTime                   securityEventType = "SettingSystemTime"                   // CRITICAL The system time on the Charging Station was changed more than ClockCtrlr.TimeAdjustmentReportingThreshold seconds
	startupOfTheDevice                  securityEventType = "StartupOfTheDevice"                  // CRITICAL The Charging Station has booted
	resetOrReboot                       securityEventType = "ResetOrReboot"                       // CRITICAL The Charging Station was rebooted or reset
	securityLogWasCleared               securityEventType = "SecurityLogWasCleared"               // CRITICAL The security log was cleared
	reconfigurationOfSecurityParameters securityEventType = "ReconfigurationOfSecurityParameters" // Security parameters, such as keys or the security profile used, were changed
	memoryExhaustion                    securityEventType = "MemoryExhaustion"                    // CRITICAL The Flash or RAM memory of the Charging Station is getting full
	invalidMessages                     securityEventType = "InvalidMessages"                     // The Charging Station has received messages that are not valid OCPP messages, if signed messages, signage invalid/incorrect
	attemptedReplayAttacks              securityEventType = "AttemptedReplayAttacks"              // The Charging Station has received a replayed message
	tamperDetectionActivated            securityEventType = "TamperDetectionActivated"            // CRITICAL The physical tamper detection sensor was triggered
	invalidFirmwareSignature            securityEventType = "InvalidFirmwareSignature"            // The firmware signature is not valid
	invalidFirmwareSigningCertificate   securityEventType = "InvalidFirmwareSigningCertificate"   // The certificate used to verify the firmware signature is not valid
	invalidCsmsCertificate              securityEventType = "InvalidCsmsCertificate"              // CRITICAL The certificate that the CSMS uses was not valid or could not be verified
	invalidChargingStationCertificate   securityEventType = "InvalidChargingStationCertificate"   // CRITICAL The certificate sent to the Charging Station using the CertificateSignedRequest message is not a valid certificate
	invalidTLSVersion                   securityEventType = "InvalidTLSVersion"                   // CRITICAL The TLS version used by the CSMS is lower than 1.2 and is not allowed by the security specification
	invalidTLSCipherSuite               securityEventType = "InvalidTLSCipherSuite"               // CRITICAL The CSMS did only allow connections using TLS cipher suites that are not allowed by the security specification
	maintenanceLoginAccepted            securityEventType = "MaintenanceLoginAccepted"            // CRITICAL Successful login to the local maintenance interface.
	maintenanceLoginFailed              securityEventType = "MaintenanceLoginFailed"              // CRITICAL Failed login attempt to the local maintenance interface.
)
