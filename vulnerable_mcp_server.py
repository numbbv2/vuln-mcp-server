#!/usr/bin/env python3
"""
Vulnerable MCP Server - Command Injection Example
이 코드는 교육 목적으로만 사용되어야 하며, 실제 프로덕션 환경에서는 사용하면 안 됩니다.
"""

import asyncio
import json
import subprocess
import sys
from typing import Any, Dict, List

# MCP 서버 기본 클래스
class MCPServer:
    def __init__(self):
        self.tools = []
        self.setup_tools()
    
    def setup_tools(self):
        """MCP 도구들을 등록합니다."""
        self.tools = [
            {
                "name": "search_files",
                "description": "파일명으로 파일을 검색합니다 (Command Injection 취약점 포함)",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "filename": {
                            "type": "string",
                            "description": "검색할 파일명"
                        }
                    },
                    "required": ["filename"]
                }
            },
            {
                "name": "list_directory",
                "description": "디렉토리 내용을 나열합니다 (Command Injection 취약점 포함)",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "path": {
                            "type": "string",
                            "description": "나열할 디렉토리 경로"
                        }
                    },
                    "required": ["path"]
                }
            },
            {
                "name": "execute_command",
                "description": "시스템 명령어를 실행합니다 (매우 위험한 Command Injection 취약점)",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "command": {
                            "type": "string",
                            "description": "실행할 명령어"
                        }
                    },
                    "required": ["command"]
                }
            }
        ]
    
    async def handle_request(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """MCP 요청을 처리합니다."""
        method = request.get("method")
        params = request.get("params", {})
        
        if method == "tools/list":
            return {
                "jsonrpc": "2.0",
                "id": request.get("id"),
                "result": {
                    "tools": self.tools
                }
            }
        
        elif method == "tools/call":
            tool_name = params.get("name")
            arguments = params.get("arguments", {})
            
            if tool_name == "search_files":
                result = await self.search_files(arguments.get("filename", ""))
            elif tool_name == "list_directory":
                result = await self.list_directory(arguments.get("path", ""))
            elif tool_name == "execute_command":
                result = await self.execute_command(arguments.get("command", ""))
            else:
                result = {"error": f"Unknown tool: {tool_name}"}
            
            return {
                "jsonrpc": "2.0",
                "id": request.get("id"),
                "result": {
                    "content": [
                        {
                            "type": "text",
                            "text": str(result)
                        }
                    ]
                }
            }
        
        else:
            return {
                "jsonrpc": "2.0",
                "id": request.get("id"),
                "error": {
                    "code": -32601,
                    "message": f"Method not found: {method}"
                }
            }
    
    async def search_files(self, filename: str) -> str:
        """
        파일명으로 파일을 검색합니다.
        
        VULNERABILITY: Command Injection
        사용자 입력이 직접 shell 명령어에 삽입되어 command injection 공격이 가능합니다.
        예: filename = "test.txt; rm -rf /"
        """
        try:
            # 위험한 코드: 사용자 입력을 직접 shell 명령어에 삽입
            command = f"find ./sandbox -name '{filename}' 2>/dev/null"
            print(f"[DEBUG] 실행할 명령어: {command}")
            
            result = subprocess.run(
                command,
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            
            if result.returncode == 0:
                return f"검색 결과:\n{result.stdout}"
            else:
                return f"검색 실패: {result.stderr}"
                
        except subprocess.TimeoutExpired:
            return "명령어 실행 시간 초과"
        except Exception as e:
            return f"오류 발생: {str(e)}"
    
    async def list_directory(self, path: str) -> str:
        """
        디렉토리 내용을 나열합니다.
        
        VULNERABILITY: Command Injection
        사용자 입력이 직접 shell 명령어에 삽입되어 command injection 공격이 가능합니다.
        예: path = "/tmp; cat /etc/passwd"
        """
        try:
            # 위험한 코드: 사용자 입력을 직접 shell 명령어에 삽입
            command = f"ls -la '{path}' 2>/dev/null"
            print(f"[DEBUG] 실행할 명령어: {command}")
            
            result = subprocess.run(
                command,
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            
            if result.returncode == 0:
                return f"디렉토리 내용:\n{result.stdout}"
            else:
                return f"디렉토리 나열 실패: {result.stderr}"
                
        except subprocess.TimeoutExpired:
            return "명령어 실행 시간 초과"
        except Exception as e:
            return f"오류 발생: {str(e)}"
    
    async def execute_command(self, command: str) -> str:
        """
        시스템 명령어를 실행합니다.
        
        VULNERABILITY: 매우 위험한 Command Injection
        사용자 입력을 그대로 shell에서 실행하므로 매우 위험합니다.
        """
        try:
            print(f"[DEBUG] 실행할 명령어: {command}")
            
            # 매우 위험한 코드: 사용자 입력을 그대로 shell에서 실행
            result = subprocess.run(
                command,
                shell=True,
                capture_output=True,
                text=True,
                timeout=30
            )
            
            output = f"명령어: {command}\n"
            output += f"반환 코드: {result.returncode}\n"
            output += f"표준 출력:\n{result.stdout}\n"
            if result.stderr:
                output += f"표준 에러:\n{result.stderr}"
            
            return output
                
        except subprocess.TimeoutExpired:
            return "명령어 실행 시간 초과"
        except Exception as e:
            return f"오류 발생: {str(e)}"

async def main():
    """MCP 서버를 실행합니다."""
    server = MCPServer()
    
    print("Vulnerable MCP Server 시작됨")
    print("주의: 이 서버는 Command Injection 취약점이 있습니다!")
    print("교육 목적으로만 사용하세요.\n")
    
    try:
        while True:
            # JSON-RPC 요청 읽기
            line = await asyncio.get_event_loop().run_in_executor(None, sys.stdin.readline)
            if not line:
                break
            
            try:
                request = json.loads(line.strip())
                response = await server.handle_request(request)
                print(json.dumps(response, ensure_ascii=False))
                sys.stdout.flush()
            except json.JSONDecodeError:
                error_response = {
                    "jsonrpc": "2.0",
                    "id": None,
                    "error": {
                        "code": -32700,
                        "message": "Parse error"
                    }
                }
                print(json.dumps(error_response))
                sys.stdout.flush()
                
    except KeyboardInterrupt:
        print("\n서버 종료됨")
    except Exception as e:
        print(f"서버 오류: {e}")

if __name__ == "__main__":
    asyncio.run(main())