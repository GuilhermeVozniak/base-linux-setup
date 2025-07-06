---
name: New Preset Request
about: Request support for a new OS/environment preset
title: '[PRESET] Add support for [OS/Environment]'
labels: 'preset', 'enhancement'
assignees: ''

---

**Operating System/Environment**
- **OS Name**: [e.g. Ubuntu Server, CentOS, Fedora]
- **Version**: [e.g. 22.04 LTS, 8.5, 39]
- **Architecture**: [e.g. x86_64, ARM64, ARM]
- **Hardware**: [e.g. Generic PC, Raspberry Pi, Cloud Instance]

**Detection Information**
Please run `neofetch --stdout` on your system and paste the output:
```
[Paste neofetch output here]
```

**Desired Setup Tasks**
What tasks should be included in this preset? Please prioritize them:

**High Priority:**
- [ ] System updates
- [ ] Package manager setup
- [ ] Essential development tools
- [ ] Other: [specify]

**Medium Priority:**
- [ ] Programming languages (Go, Python, Node.js)
- [ ] Development tools (git, vim, etc.)
- [ ] System utilities
- [ ] Other: [specify]

**Low Priority:**
- [ ] Optional software (Docker, VS Code)
- [ ] Customizations (aliases, themes)
- [ ] Other: [specify]

**Special Requirements**
Are there any OS-specific commands, package managers, or configuration requirements?
- Package manager: [e.g. apt, yum, pacman, zypper]
- Service manager: [e.g. systemd, openrc, runit]
- Special considerations: [e.g. SELinux, specific repositories]

**Sample Commands**
If you know the specific commands needed for this OS, please provide them:
```bash
# System update
[commands here]

# Package installation
[commands here]

# Service management
[commands here]
```

**Additional Context**
Any other information that would be helpful for implementing this preset. 