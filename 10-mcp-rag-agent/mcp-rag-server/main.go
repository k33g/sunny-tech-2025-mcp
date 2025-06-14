package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"mcp-dd-rag/rag"
)

var client openai.Client
var store rag.MemoryVectorStore
var embeddingsModel string

func main() {
	ctx := context.Background()

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-dd-rag",
		"0.0.0",
	)
	err := godotenv.Load()
	if err != nil {
		// TODO:
	}
	// Ensure MODEL_RUNNER_BASE_URL is set in the environment
	if os.Getenv("MODEL_RUNNER_BASE_URL") == "" {
		os.Setenv("MODEL_RUNNER_BASE_URL", "http://localhost:12434")
	}
	if os.Getenv("EMBEDDING_MODEL") == "" {
		os.Setenv("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest")
	}

	llmURL := os.Getenv("MODEL_RUNNER_BASE_URL") + "/engines/llama.cpp/v1/"
	embeddingsModel = os.Getenv("EMBEDDING_MODEL")

	client = openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)

	// =================================================
	// CHUNKS:
	// =================================================
	contents, err := GetContentFiles(".", ".md")
	if err != nil {
		log.Fatalln("😡 Error getting content files:", err)
	}
	chunks := []string{}
	for _, content := range contents {
		chunks = append(chunks, ChunkText(content, 512, 210)...)
	}

	//fmt.Println(chunks)

	// -------------------------------------------------
	// Create a vector store
	// -------------------------------------------------
	store = rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	// -------------------------------------------------
	// Create and save the embeddings from the chunks
	// -------------------------------------------------
	fmt.Println("⏳ Creating the embeddings...")

	for _, chunk := range chunks {
		embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(chunk),
			},
			Model: embeddingsModel,
		})

		if err != nil {
			fmt.Println(err)
		} else {
			_, errSave := store.Save(rag.VectorRecord{
				Prompt:    chunk,
				Embedding: embeddingsResponse.Data[0].Embedding,
			})
			if errSave != nil {
				fmt.Println("😡:", errSave)
			}
		}
	}

	fmt.Println("✋", "Embeddings created, total of records", len(store.Records))
	fmt.Println()

	// =================================================
	// TOOLS:
	// =================================================
	searchInDoc := mcp.NewTool("question_about_something",
		//mcp.WithDescription(`Find an answer in the internal database.`),
		mcp.WithDescription(`Trouve une réponse dans ta base de données interne.`),
		mcp.WithString("question",
			mcp.Required(),
			//mcp.Description("Search question"),
			mcp.Description("Question de recherche"),
		),
	)
	s.AddTool(searchInDoc, searchInDocHandler)

	// Start the stdio server
	// if err := server.ServeStdio(s); err != nil {
	// 	log.Fatalln("Failed to start server:", err)
	// 	return
	// }

	// Start the HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	).Start(":" + httpPort)
}

func searchInDocHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	args := request.GetArguments()
	userQuestion := args["question"].(string)

	fmt.Println("🔍 Searching for question:", userQuestion)

	// -------------------------------------------------
	// Search for similarities
	// -------------------------------------------------

	fmt.Println("⏳ Searching for similarities...")

	// -------------------------------------------------
	// Create embedding from the user question
	// -------------------------------------------------
	embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(userQuestion),
		},
		Model: embeddingsModel,
	})
	if err != nil {
		log.Fatal("😡:", err)
	}

	// -------------------------------------------------
	// Create a vector record from the user embedding
	// -------------------------------------------------
	embeddingFromUserQuestion := rag.VectorRecord{
		Embedding: embeddingsResponse.Data[0].Embedding,
	}

	strLimit := os.Getenv("LIMIT")
	if strLimit == "" {
		strLimit = "0.6"
	}
	strMax := os.Getenv("MAX_RESULTS")
	if strMax == "" {
		strMax = "2"
	}
	// Convert string to float64 and int
	var limit float64
	fmt.Sscanf(strLimit, "%f", &limit)
	var maxResults int
	fmt.Sscanf(strMax, "%d", &maxResults)

	similarities, _ := store.SearchTopNSimilarities(embeddingFromUserQuestion, limit, maxResults)

	documentsContent := "Documents:\n"

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
		documentsContent += similarity.Prompt
	}
	documentsContent += "\n"
	fmt.Println("✋", "Similarities found, total of records", len(similarities))
	fmt.Println()

	// -------------------------------------------------
	// Generate embeddings from user question
	// -------------------------------------------------
	// EMBEDDINGS...
	return mcp.NewToolResultText(documentsContent), nil
}

// ChunkText takes a text string and divides it into chunks of a specified size with a given overlap.
// It returns a slice of strings, where each string represents a chunk of the original text.
//
// Parameters:
//   - text: The input text to be chunked.
//   - chunkSize: The size of each chunk.
//   - overlap: The amount of overlap between consecutive chunks.
//
// Returns:
//   - []string: A slice of strings representing the chunks of the original text.
func ChunkText(text string, chunkSize, overlap int) []string {
	chunks := []string{}
	for start := 0; start < len(text); start += chunkSize - overlap {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[start:end])
	}
	return chunks
}

// GetContentFiles searches for files with a specific extension in the given directory and its subdirectories.
//
// Parameters:
// - dirPath: The directory path to start the search from.
// - ext: The file extension to search for.
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func GetContentFiles(dirPath string, ext string) ([]string, error) {
	content := []string{}
	_, err := ForEachFile(dirPath, ext, func(path string) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content = append(content, string(data))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ForEachFile iterates over all files with a specific extension in a directory and its subdirectories.
//
// Parameters:
// - dirPath: The root directory to start the search from.
// - ext: The file extension to search for.
// - callback: A function to be called for each file found.
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func ForEachFile(dirPath string, ext string, callback func(string) error) ([]string, error) {
	var textFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			textFiles = append(textFiles, path)
			err = callback(path)
			// generate an error to stop the walk
			if err != nil {
				return err
			}
		}
		return nil
	})
	return textFiles, err
}
