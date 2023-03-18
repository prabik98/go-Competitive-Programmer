package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apikey := viper.GetString("API_KEY")
	if apikey == "" {
		panic("Missing API_KEY in environment variable")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apikey)

	const inputFile = "./input_with_code.txt"
	fileBytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}
	msgPrefix := "Solve this Leetcode Question \n```c++\n"
	msgSuffix := "\n```"
	msg := msgPrefix + string(fileBytes) + msgSuffix

	outputBuilder := strings.Builder{}
	client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			msg,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		outputBuilder.WriteString(resp.Choices[0].Text)
	})
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.TrimSpace(outputBuilder.String())
	const outputFile = "./output.txt"
	os.WriteFile(outputFile, []byte(output), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
}
