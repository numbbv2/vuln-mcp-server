# Security Policy

## ⚠️ IMPORTANT SECURITY NOTICE ⚠️

**This repository contains intentionally vulnerable code for educational purposes only.**

## Supported Versions

This software is **NOT SUPPORTED** for production use. All versions are intentionally vulnerable and should only be used in controlled educational environments.

## Known Vulnerabilities

This software intentionally contains the following security vulnerabilities:

### 1. Command Injection in `search_files()`
- **Location**: `vulnerable_mcp_server.py:search_files()`
- **Description**: User input is directly inserted into shell commands
- **Risk**: Arbitrary command execution
- **Example**: `filename = "test.txt; rm -rf /"`

### 2. Command Injection in `list_directory()`
- **Location**: `vulnerable_mcp_server.py:list_directory()`
- **Description**: Path input is directly inserted into shell commands
- **Risk**: Arbitrary command execution
- **Example**: `path = "/tmp; cat /etc/passwd"`

### 3. Command Injection in `execute_command()`
- **Location**: `vulnerable_mcp_server.py:execute_command()`
- **Description**: User input is directly executed as shell commands
- **Risk**: Complete system compromise
- **Example**: `command = "rm -rf /"`

## Reporting a Vulnerability

**IMPORTANT**: These vulnerabilities are intentional and documented for educational purposes.

If you discover additional vulnerabilities or have security concerns:

1. **DO NOT** create public issues for intentional vulnerabilities
2. **DO NOT** use this software in production environments
3. **DO NOT** attempt to exploit these vulnerabilities on systems you don't own

For responsible disclosure of new vulnerabilities:

- **Email**: security@example.com
- **Subject**: [SECURITY] Vulnerability Report - vuln-mcp-server
- **Include**: Detailed description, steps to reproduce, potential impact

## Security Best Practices

### For Educational Use:
1. Use only in isolated virtual machines
2. Use Docker containers with restricted privileges
3. Never use in production environments
4. Ensure proper network isolation
5. Use non-privileged user accounts

### For Developers:
1. Never copy code from this repository to production systems
2. Use proper input validation and sanitization
3. Avoid `shell=True` in subprocess calls
4. Use parameterized commands instead of string concatenation
5. Implement proper access controls

## Malicious Use Prohibition

**STRICTLY PROHIBITED:**
- Using this software to attack systems you don't own
- Using this software in production environments
- Distributing this software without proper warnings
- Modifying this software to hide its educational purpose

## Responsible Disclosure

We follow responsible disclosure practices for any newly discovered vulnerabilities:

1. **Initial Report**: Contact security team privately
2. **Verification**: Confirm vulnerability exists
3. **Timeline**: 90-day disclosure timeline
4. **Coordination**: Work with reporter on disclosure
5. **Public Disclosure**: After patch or 90 days

## Security Resources

- [OWASP Command Injection](https://owasp.org/www-community/attacks/Command_Injection)
- [Python Security Best Practices](https://docs.python.org/3/library/subprocess.html#security-considerations)
- [Secure Coding Guidelines](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)

## Contact Information

- **Security Team**: security@example.com
- **Project Maintainer**: maintainer@example.com
- **Emergency Contact**: emergency@example.com

---

**Remember: This software is for educational purposes only. Use responsibly and ethically.**