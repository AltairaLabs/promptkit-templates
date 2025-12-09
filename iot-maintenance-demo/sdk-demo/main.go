package main

import (
"context"
"fmt"
"log"
"os"

"github.com/AltairaLabs/PromptKit/sdk"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "usage: go run ./sdk-demo <pack.json> <prompt-id> <message>\n")
		os.Exit(1)
	}

	packPath := os.Args[1]
	promptID := os.Args[2]
	userMsg := os.Args[3]

	customerID := envOrDefault("CUSTOMER_ID", "acme-corp")
	customerName := envOrDefault("CUSTOMER_NAME", "ACME Corporation")

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY is required to run the SDK demo")
	}

	// SDK v2: Open conversation directly from pack file
	conv, err := sdk.Open(packPath, promptID)
	if err != nil {
		log.Fatalf("failed to open conversation: %v", err)
	}
	defer conv.Close()

	// Set template variables
	conv.SetVars(map[string]any{
"customer_id":   customerID,
"customer_name": customerName,
})

	// Register tool handlers
	registerToolHandlers(conv, customerID)

	fmt.Printf("┌───────────────────────────────────────────────┐\n")
	fmt.Printf("│  IoT Troubleshooter Demo (SDK v2)\n")
	fmt.Printf("│  Customer: %s (%s)\n", customerName, customerID)
	fmt.Printf("│  Prompt:   %s\n", promptID)
	fmt.Printf("└───────────────────────────────────────────────┘\n\n")

	ctx := context.Background()
	fmt.Printf("👤 User: %s\n\n", userMsg)

	fmt.Print("🤖 Assistant: ")

	// Stream the response token by token
	var finalResp *sdk.Response
	for chunk := range conv.Stream(ctx, userMsg) {
		if chunk.Error != nil {
			log.Fatalf("stream error: %v", chunk.Error)
		}

		switch chunk.Type {
		case sdk.ChunkText:
			fmt.Print(chunk.Text)
		case sdk.ChunkToolCall:
			// Show tool calls as they happen
			fmt.Printf("\n  🔧 Calling %s... ✓\n", chunk.ToolCall.Name)
		case sdk.ChunkDone:
			finalResp = chunk.Message
		}
	}
	fmt.Println("\n")

	if finalResp != nil {
		if finalResp.HasToolCalls() {
			fmt.Println("🔧 Tools Called:")
			for _, tc := range finalResp.ToolCalls() {
				fmt.Printf("  - %s\n", tc.Name)
			}
			fmt.Println()
		}

		fmt.Printf("💰 Cost: $%.4f | Duration: %v\n", finalResp.Cost(), finalResp.Duration())
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// registerToolHandlers registers all IoT tool handlers with the conversation
func registerToolHandlers(conv *sdk.Conversation, defaultCustomerID string) {
	// SECURITY: All handlers use the session's customerID, ignoring what the LLM sends.
// This enforces tenant isolation - the LLM cannot access data from other customers.
conv.OnTools(map[string]sdk.ToolHandler{
"list_devices": func(args map[string]any) (any, error) {
return listDevices(defaultCustomerID)
},
"get_sensor_data": func(args map[string]any) (any, error) {
deviceID := getStringArg(args, "device_id", "")
return getSensorData(defaultCustomerID, deviceID)
},
"get_error_logs": func(args map[string]any) (any, error) {
deviceID := getStringArg(args, "device_id", "")
return getErrorLogs(defaultCustomerID, deviceID)
},
"check_maintenance_schedule": func(args map[string]any) (any, error) {
deviceID := getStringArg(args, "device_id", "")
return checkMaintenanceSchedule(defaultCustomerID, deviceID)
},
"get_customer_details": func(args map[string]any) (any, error) {
return getCustomerDetails(defaultCustomerID)
},
"book_engineer": func(args map[string]any) (any, error) {
return bookEngineer(args)
},
})
}

func getStringArg(args map[string]any, key, defaultVal string) string {
if v, ok := args[key].(string); ok && v != "" {
return v
}
return defaultVal
}

// Mock tool implementations
func listDevices(customerID string) (any, error) {
if customerID == "acme-corp" {
return map[string]any{
"devices": []map[string]string{
{"device_id": "MOTOR-001", "type": "Industrial Motor", "location": "Assembly Line 3", "status": "warning"},
{"device_id": "PUMP-002", "type": "Industrial Pump", "location": "Cooling System", "status": "critical"},
{"device_id": "CONVEYOR-003", "type": "Conveyor Belt", "location": "Packaging", "status": "ok"},
},
}, nil
}
return map[string]any{"devices": []any{}, "error": "unknown customer"}, nil
}

func getSensorData(customerID, deviceID string) (any, error) {
if customerID != "acme-corp" {
return map[string]string{"error": "unauthorized customer"}, nil
}
switch deviceID {
case "MOTOR-001":
return map[string]any{
"temperature": 85, "vibration": 2.3, "pressure": 45, "status": "warning",
"message": "Elevated temperature - check cooling fan",
}, nil
case "PUMP-002":
return map[string]any{
"temperature": 120, "vibration": 8.7, "pressure": 12, "status": "critical",
"message": "Bearing failure imminent - excessive vibration",
}, nil
case "CONVEYOR-003":
return map[string]any{
"temperature": 45, "vibration": 0.8, "pressure": 30, "status": "ok",
"message": "All readings nominal",
}, nil
}
return map[string]string{"error": "device not found"}, nil
}

func getErrorLogs(customerID, deviceID string) (any, error) {
if customerID != "acme-corp" {
return map[string]string{"error": "unauthorized"}, nil
}
switch deviceID {
case "PUMP-002":
return map[string]any{
"logs":    []string{"Vibration spike at 09:45", "Vibration spike at 09:30", "Pressure drop at 09:10"},
"summary": "Repeated vibration spikes",
}, nil
default:
return map[string]any{"logs": []string{}, "summary": "No errors"}, nil
}
}

func checkMaintenanceSchedule(customerID, deviceID string) (any, error) {
return map[string]any{
"last_service":     "2025-09-15",
"next_service_due": "2025-12-15",
"notes":            "Due for cooling system cleaning",
}, nil
}

func getCustomerDetails(customerID string) (any, error) {
if customerID == "acme-corp" {
return map[string]any{
"company_name":  "ACME Corporation",
"contact_name":  "Jordan Rivera",
"contact_email": "maintenance@acme.example",
"contact_phone": "+1-555-0101",
}, nil
}
return map[string]string{"error": "unknown customer"}, nil
}

func bookEngineer(params map[string]any) (any, error) {
return map[string]any{
"booking_id":         "ENG-2025-1206-001",
"engineer":           "Sarah Chen",
"scheduled_date":     "2025-12-07",
"estimated_duration": "2 hours",
"confirmation":       "Engineer visit scheduled. Customer will receive email confirmation.",
}, nil
}
