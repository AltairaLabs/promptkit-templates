package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/AltairaLabs/PromptKit/runtime/providers"
    "github.com/AltairaLabs/PromptKit/runtime/providers/openai"
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
    model := envOrDefault("OPENAI_MODEL", "gpt-4o")
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY is required to run the SDK demo")
    }

    provider := openai.NewProvider(
        "openai",
        model,
        "https://api.openai.com/v1",
        providers.ProviderDefaults{Temperature: 0.7, MaxTokens: 900, TopP: 1.0},
        false,
    )

    manager, err := sdk.NewConversationManager(
        sdk.WithProvider(provider),
    )
    if err != nil {
        log.Fatalf("failed to create manager: %v", err)
    }

    pack, err := manager.LoadPack(packPath)
    if err != nil {
        log.Fatalf("failed to load pack %s: %v", packPath, err)
    }

    ctx := context.Background()
    conv, err := manager.NewConversation(ctx, pack, sdk.ConversationConfig{
        PromptName: promptID,
        Variables: map[string]interface{}{
            "customer_id":   customerID,
            "customer_name": customerName,
        },
    })
    if err != nil {
        log.Fatalf("failed to create conversation: %v", err)
    }

    fmt.Printf("┌───────────────────────────────────────────────┐\n")
    fmt.Printf("│  IoT Troubleshooter Demo\n")
    fmt.Printf("│  Customer: %s (%s)\n", customerName, customerID)
    fmt.Printf("│  Prompt:   %s\n", promptID)
    fmt.Printf("└───────────────────────────────────────────────┘\n\n")

    fmt.Printf("👤 User: %s\n\n", userMsg)
    resp, err := conv.Send(ctx, userMsg)
    if err != nil {
        log.Fatalf("send failed: %v", err)
    }

    fmt.Printf("🤖 Assistant: %s\n\n", resp.Content)
    if len(resp.ToolCalls) > 0 {
        fmt.Println("🔧 Tools Called:")
        for _, tc := range resp.ToolCalls {
            fmt.Printf("  - %s\n", tc.Name)
        }
        fmt.Println()
    }

    fmt.Printf("💰 Cost: $%.4f | Tokens: %d\n", resp.Cost, resp.Tokens)
}

func envOrDefault(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
