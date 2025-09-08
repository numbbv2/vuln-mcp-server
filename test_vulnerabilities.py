#!/usr/bin/env python3
"""
Command Injection 취약점 테스트 스크립트
이 스크립트는 교육 목적으로만 사용되어야 합니다.
"""

import json
import subprocess
import time
import sys

def send_mcp_request(method, params=None, request_id=1):
    """MCP 서버에 JSON-RPC 요청을 보냅니다."""
    request = {
        "jsonrpc": "2.0",
        "id": request_id,
        "method": method
    }
    
    if params:
        request["params"] = params
    
    return json.dumps(request)

def test_search_files_vulnerability():
    """search_files 함수의 Command Injection 취약점을 테스트합니다."""
    print("=== search_files Command Injection 테스트 ===")
    
    # 정상적인 요청
    print("\n1. 정상적인 요청:")
    normal_request = send_mcp_request(
        "tools/call",
        {
            "name": "search_files",
            "arguments": {"filename": "*.txt"}
        }
    )
    print(f"요청: {normal_request}")
    
    # Command Injection 공격 시도
    print("\n2. Command Injection 공격 시도:")
    malicious_request = send_mcp_request(
        "tools/call",
        {
            "name": "search_files",
            "arguments": {"filename": "test.txt; echo 'COMMAND INJECTION SUCCESSFUL'"}
        }
    )
    print(f"요청: {malicious_request}")
    
    # 더 위험한 공격 시도
    print("\n3. 더 위험한 공격 시도:")
    dangerous_request = send_mcp_request(
        "tools/call",
        {
            "name": "search_files",
            "arguments": {"filename": "test.txt && whoami"}
        }
    )
    print(f"요청: {dangerous_request}")

def test_list_directory_vulnerability():
    """list_directory 함수의 Command Injection 취약점을 테스트합니다."""
    print("\n=== list_directory Command Injection 테스트 ===")
    
    # 정상적인 요청
    print("\n1. 정상적인 요청:")
    normal_request = send_mcp_request(
        "tools/call",
        {
            "name": "list_directory",
            "arguments": {"path": "/tmp"}
        }
    )
    print(f"요청: {normal_request}")
    
    # Command Injection 공격 시도
    print("\n2. Command Injection 공격 시도:")
    malicious_request = send_mcp_request(
        "tools/call",
        {
            "name": "list_directory",
            "arguments": {"path": "/tmp; echo 'DIRECTORY INJECTION SUCCESSFUL'"}
        }
    )
    print(f"요청: {malicious_request}")

def test_execute_command_vulnerability():
    """execute_command 함수의 Command Injection 취약점을 테스트합니다."""
    print("\n=== execute_command Command Injection 테스트 ===")
    
    # 정상적��� 요청
    print("\n1. 정상적인 요청:")
    normal_request = send_mcp_request(
        "tools/call",
        {
            "name": "execute_command",
            "arguments": {"command": "ls -la"}
        }
    )
    print(f"요청: {normal_request}")
    
    # Command Injection 공격 시도
    print("\n2. Command Injection 공격 시도:")
    malicious_request = send_mcp_request(
        "tools/call",
        {
            "name": "execute_command",
            "arguments": {"command": "echo 'EXECUTE COMMAND INJECTION SUCCESSFUL'"}
        }
    )
    print(f"요청: {malicious_request}")
    
    # 매우 위험한 공격 시도 (실제로는 실행하지 않음)
    print("\n3. 매우 위험한 공격 시도 (예시만):")
    dangerous_request = send_mcp_request(
        "tools/call",
        {
            "name": "execute_command",
            "arguments": {"command": "rm -rf /tmp/test_file"}
        }
    )
    print(f"요청: {dangerous_request}")
    print("⚠️  실제로는 실행하지 않습니다!")

def demonstrate_safe_alternatives():
    """안전한 대안 방법들을 보여줍니다."""
    print("\n=== 안전한 대안 방법들 ===")
    
    print("\n1. subprocess.run에서 shell=False 사용:")
    safe_code = '''
import subprocess

# 안전한 방법
filename = "test.txt"
result = subprocess.run(
    ["find", "./sandbox", "-name", filename],
    capture_output=True,
    text=True
)
'''
    print(safe_code)
    
    print("\n2. 입력 검증 및 이스케이프:")
    validation_code = '''
import shlex
import re

def validate_filename(filename):
    # 파일명 검증
    if not re.match(r'^[a-zA-Z0-9._-]+$', filename):
        raise ValueError("Invalid filename")
    return shlex.quote(filename)

# 사용 예시
safe_filename = validate_filename("test.txt")
command = f"find ./sandbox -name {safe_filename}"
'''
    print(validation_code)

def main():
    """메인 함수"""
    print("Command Injection 취약점 테스트 스크립트")
    print("=" * 50)
    print("⚠️  경고: 이 스크립트는 교육 목적으로만 사용되어야 합니다!")
    print("⚠️  실제 시스템에서 테스트할 때는 격리된 환경에서만 수행하세요.")
    print("=" * 50)
    
    # 취약점 테스트
    test_search_files_vulnerability()
    test_list_directory_vulnerability()
    test_execute_command_vulnerability()
    
    # 안전한 대안 제시
    demonstrate_safe_alternatives()
    
    print("\n" + "=" * 50)
    print("테스트 완료!")
    print("실제 MCP 서버를 실행하려면: python vulnerable_mcp_server.py")
    print("=" * 50)

if __name__ == "__main__":
    main()