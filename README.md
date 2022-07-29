# Edgecreaft API Server

<img src="https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;<a href="https://goreportcard.com/report/github.com/create-go-app/fiber-go-template" target="_blank"><img src="https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none" alt="go report" /></a>&nbsp;<img src="https://img.shields.io/badge/license-Apache_2.0-red?style=for-the-badge&logo=none" alt="license" />

## 구성요소
- golang `v1.18`
- echo `v4.72`
- gorp `v2.2.0`


-----
## 🗄 Directory structure
### ./conf
**Folder with configuration files and response message guide files.**

### ./cmd
**Main applications for this project.**

### ./pkg

**Library code that's ok to use by external applications.**. This directory contains all the project-specific code tailored only for your business use case, like _configs_, _middleware_, _routes_ or _utils_.

- `./pkg/configs` folder for configuration functions
- `./pkg/middleware` folder for add middleware (Fiber built-in and yours)
- `./pkg/repository` folder for describe `const` of your project
- `./pkg/utils` folder with utility functions (server starter, error checker, etc)


### ./docs
**Folder with 사용자 문서들. and Swagger 스펙들.**

### ./scripts
**빌드, 설치, 분석, 기타 작업을 위한 스크립트들.**