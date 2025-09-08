package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// MCPRequest represents an MCP request for testing
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// sendMCPRequest creates an MCP request for testing
func sendMCPRequest(method string, params map[string]interface{}, requestID int) string {
	request := MCPRequest{
		JSONRPC: "2.0",
		ID:      requestID,
		Method:  method,
		Params:  params,
	}
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return ""
	}
	
	return string(jsonData)
}

// testSearchFilesVulnerability tests Command Injection in search_files function
func testSearchFilesVulnerability() {
	fmt.Println("=== search_files Command Injection 테스트 ===")
	
	// 정상적인 요청
	fmt.Println("\n1. 정상적인 요청:")
	normalRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "search_files",
		"arguments": map[string]interface{}{
			"filename": "*.txt",
		},
	}, 1)
	fmt.Printf("요청: %s\n", normalRequest)
	
	// Command Injection 공격 시도
	fmt.Println("\n2. Command Injection 공격 시도:")
	maliciousRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "search_files",
		"arguments": map[string]interface{}{
			"filename": "test.txt; echo 'COMMAND INJECTION SUCCESSFUL'",
		},
	}, 2)
	fmt.Printf("요청: %s\n", maliciousRequest)
	
	// 더 위험한 공격 시도
	fmt.Println("\n3. 더 위험한 공격 시도:")
	dangerousRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "search_files",
		"arguments": map[string]interface{}{
			"filename": "test.txt && whoami",
		},
	}, 3)
	fmt.Printf("요청: %s\n", dangerousRequest)
}

// testListDirectoryVulnerability tests Command Injection in list_directory function
func testListDirectoryVulnerability() {
	fmt.Println("\n=== list_directory Command Injection 테스트 ===")
	
	// 정상적인 요청
	fmt.Println("\n1. 정상적인 요청:")
	normalRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "list_directory",
		"arguments": map[string]interface{}{
			"path": "/tmp",
		},
	}, 4)
	fmt.Printf("요청: %s\n", normalRequest)
	
	// Command Injection 공격 시도
	fmt.Println("\n2. Command Injection 공격 시도:")
	maliciousRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "list_directory",
		"arguments": map[string]interface{}{
			"path": "/tmp; echo 'DIRECTORY INJECTION SUCCESSFUL'",
		},
	}, 5)
	fmt.Printf("요청: %s\n", maliciousRequest)
}

// testExecuteCommandVulnerability tests Command Injection in execute_command function
func testExecuteCommandVulnerability() {
	fmt.Println("\n=== execute_command Command Injection 테스트 ===")
	
	// 정상적인 요청
	fmt.Println("\n1. 정상적인 요청:")
	normalRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "execute_command",
		"arguments": map[string]interface{}{
			"command": "ls -la",
		},
	}, 6)
	fmt.Printf("요청: %s\n", normalRequest)
	
	// Command Injection 공격 시도
	fmt.Println("\n2. Command Injection 공격 시도:")
	maliciousRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "execute_command",
		"arguments": map[string]interface{}{
			"command": "echo 'EXECUTE COMMAND INJECTION SUCCESSFUL'",
		},
	}, 7)
	fmt.Printf("요청: %s\n", maliciousRequest)
	
	// 매우 위험한 공격 시도 (실제로는 실행하지 않음)
	fmt.Println("\n3. 매우 위험한 공격 시도 (예시만):")
	dangerousRequest := sendMCPRequest("tools/call", map[string]interface{}{
		"name": "execute_command",
		"arguments": map[string]interface{}{
			"command": "rm -rf /tmp/test_file",
		},
	}, 8)
	fmt.Printf("요청: %s\n", dangerousRequest)
	fmt.Println("⚠️  실제로는 실행하지 않습니다!")
}

// demonstrateSafeAlternatives shows safe coding practices
func demonstrateSafeAlternatives() {
	fmt.Println("\n=== 안전한 대안 방법들 ===")
	
	fmt.Println("\n1. exec.Command에서 직접 인수 전달:")
	safeCode := `
package main

import (
    "os/exec"
    "context"
)

// 안전한 방법
func safeSearchFiles(filename string) (string, error) {
    // 직접 인수를 전달하여 shell injection 방지
    cmd := exec.CommandContext(context.Background(), "find", "./sandbox", "-name", filename)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
`
	fmt.Println(safeCode)
	
	fmt.Println("\n2. 입력 검증 및 이스케이프:")
	validationCode := `
package main

import (
    "regexp"
    "strings"
    "os/exec"
)

func validateFilename(filename string) error {
    // 파일명 검증
    matched, err := regexp.MatchString("^[a-zA-Z0-9._-]+$", filename)
    if err != nil || !matched {
        return fmt.Errorf("invalid filename")
    }
    return nil
}

// 사용 예시
func safeSearchFiles(filename string) (string, error) {
    if err := validateFilename(filename); err != nil {
        return "", err
    }
    
    // 안전한 명령어 실행
    cmd := exec.Command("find", "./sandbox", "-name", filename)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
`
	fmt.Println(validationCode)
	
	fmt.Println("\n3. 허용된 명령어만 실행:")
	whitelistCode := `
package main

import (
    "fmt"
    "os/exec"
)

var allowedCommands = map[string]bool{
    "ls":   true,
    "find": true,
    "grep": true,
}

func executeSafeCommand(commandName string, args ...string) (string, error) {
    if !allowedCommands[commandName] {
        return "", fmt.Errorf("command not allowed: %s", commandName)
    }
    
    cmd := exec.Command(commandName, args...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
`
	fmt.Println(whitelistCode)
}

// main function runs all vulnerability tests
func main() {
	fmt.Println("Command Injection 취약점 테스트 스크립트 (Go)")
	fmt.Println("=" * 50)
	fmt.Println("⚠️  경고: 이 스크립트는 교육 목적으로만 사용되어야 합니다!")
	fmt.Println("⚠️  실제 시스템에서 테스트할 때는 격리된 환경에서만 수행하세요.")
	fmt.Println("=" * 50)
	
	// 취약점 테스트
	testSearchFilesVulnerability()
	testListDirectoryVulnerability()
	testExecuteCommandVulnerability()
	
	// 안전한 대안 제시
	demonstrateSafeAlternatives()
	
	fmt.Println("\n" + "=" * 50)
	fmt.Println("테스트 완료!")
	fmt.Println("실제 MCP 서버를 실행하려면: go run vulnerable_mcp_server.go")
	fmt.Println("=" * 50)
}