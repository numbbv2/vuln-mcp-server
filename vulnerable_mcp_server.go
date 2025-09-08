package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// MCPTool represents a tool in the MCP server
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// MCPResponse represents an outgoing MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPContent represents content in a response
type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MCPServer represents the vulnerable MCP server
type MCPServer struct {
	tools []MCPTool
}

// NewMCPServer creates a new vulnerable MCP server
func NewMCPServer() *MCPServer {
	server := &MCPServer{}
	server.setupTools()
	return server
}

// setupTools initializes the MCP tools with intentional vulnerabilities
func (s *MCPServer) setupTools() {
	s.tools = []MCPTool{
		{
			Name:        "search_files",
			Description: "파일명으로 파일을 검색합니다 (Command Injection 취약점 포함)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "검색할 파일명",
					},
				},
				"required": []string{"filename"},
			},
		},
		{
			Name:        "list_directory",
			Description: "디렉토리 내용을 나열합니다 (Command Injection 취약점 포함)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "나열할 디렉토리 경로",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "execute_command",
			Description: "시스템 명령어를 실행합니다 (매우 위험한 Command Injection 취약점)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "실행할 명령어",
					},
				},
				"required": []string{"command"},
			},
		},
	}
}

// HandleRequest processes incoming MCP requests
func (s *MCPServer) HandleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "tools/list":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": s.tools,
			},
		}
	case "tools/call":
		return s.handleToolCall(req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

// handleToolCall processes tool call requests
func (s *MCPServer) handleToolCall(req MCPRequest) MCPResponse {
	params := req.Params
	toolName, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params: name is required",
			},
		}
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	var result string
	var err error

	switch toolName {
	case "search_files":
		filename, _ := arguments["filename"].(string)
		result, err = s.searchFiles(filename)
	case "list_directory":
		path, _ := arguments["path"].(string)
		result, err = s.listDirectory(path)
	case "execute_command":
		command, _ := arguments["command"].(string)
		result, err = s.executeCommand(command)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Unknown tool: %s", toolName),
			},
		}
	}

	if err != nil {
		result = fmt.Sprintf("Error: %v", err)
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []MCPContent{
				{
					Type: "text",
					Text: result,
				},
			},
		},
	}
}

// searchFiles searches for files by name (VULNERABLE: Command Injection)
func (s *MCPServer) searchFiles(filename string) (string, error) {
	// VULNERABILITY: Command Injection
	// 사용자 입력이 직접 shell 명령어에 삽입되어 command injection 공격이 가능합니다.
	// 예: filename = "test.txt; rm -rf /"
	
	// 위험한 코드: 사용자 입력을 직접 shell 명령어에 삽입
	command := fmt.Sprintf("find ./sandbox -name '%s' 2>/dev/null", filename)
	fmt.Printf("[DEBUG] 실행할 명령어: %s\n", command)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Execute command with shell (VULNERABLE!)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Sprintf("검색 실패: %v", err), nil
	}
	
	return fmt.Sprintf("검색 결과:\n%s", string(output)), nil
}

// listDirectory lists directory contents (VULNERABLE: Command Injection)
func (s *MCPServer) listDirectory(path string) (string, error) {
	// VULNERABILITY: Command Injection
	// 사용자 입력이 직접 shell 명령어에 삽입되어 command injection 공격이 가능합니다.
	// 예: path = "/tmp; cat /etc/passwd"
	
	// 위험한 코드: 사용자 입력을 직접 shell 명령어에 삽입
	command := fmt.Sprintf("ls -la '%s' 2>/dev/null", path)
	fmt.Printf("[DEBUG] 실행할 명령어: %s\n", command)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Execute command with shell (VULNERABLE!)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Sprintf("디렉토리 나열 실패: %v", err), nil
	}
	
	return fmt.Sprintf("디렉토리 내용:\n%s", string(output)), nil
}

// executeCommand executes system commands (VULNERABLE: Command Injection)
func (s *MCPServer) executeCommand(command string) (string, error) {
	// VULNERABILITY: 매우 위험한 Command Injection
	// 사용자 입력을 그대로 shell에서 실행하므로 매우 위험합니다.
	
	fmt.Printf("[DEBUG] 실행할 명령어: %s\n", command)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 매우 위험한 코드: 사용자 입력을 그대로 shell에서 실행
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	
	result := fmt.Sprintf("명령어: %s\n", command)
	result += fmt.Sprintf("반환 코드: %d\n", cmd.ProcessState.ExitCode())
	result += fmt.Sprintf("표준 출력:\n%s", string(output))
	
	if err != nil {
		result += fmt.Sprintf("\n오류: %v", err)
	}
	
	return result, nil
}

// main function starts the vulnerable MCP server
func main() {
	server := NewMCPServer()
	
	fmt.Println("Vulnerable MCP Server (Go) 시작됨")
	fmt.Println("주의: 이 서버는 Command Injection 취약점이 있습니다!")
	fmt.Println("교육 목적으로만 사용하세요.\n")
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			errorResp := MCPResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error: &MCPError{
					Code:    -32700,
					Message: "Parse error",
				},
			}
			response, _ := json.Marshal(errorResp)
			fmt.Println(string(response))
			continue
		}
		
		response := server.HandleRequest(req)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}
		
		fmt.Println(string(responseJSON))
	}
	
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}