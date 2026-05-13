## AjiraScan

AjiraScan is a Go-powered ATS (Applicant Tracking System) analysis platform built to help job seekers optimize their CVs for modern recruitment systems.

The platform analyzes resumes against job descriptions, calculates ATS compatibility scores, identifies missing keywords, and provides practical recommendations to improve interview chances.

AjiraScan is designed with the Kenyan and African job market in mind, helping graduates, professionals, NGO applicants, and corporate job seekers better align their applications with employer expectations.

## Features

- ATS Resume Scoring
- Keyword Matching
- Missing Skills Detection
- Resume Optimization Suggestions
- Job Description Analysis
- Recruiter-Oriented Scoring Engine
- Fast CLI-Based Processing
- Built with Go for performance and scalability

## Project Structure

```bash
ajirascan/
│
├── cmd/
│   └── cli/
│       └── main.go
│
├── internal/
│   ├── ats/
│   │   ├── engine.go
│   │   ├── matcher.go
│   │   ├── scorer.go
│   │   └── *_test.go
│   │
│   └── text/
│       ├── normalize.go
│       ├── tokenize.go
│       ├── frequency.go
│       └── *_test.go
│
├── sample_cv.txt
├── sample_job.txt
├── go.mod
└── README.md
```

## Running the Project

## Clone the repository

```bash
git clone <repository-url>
cd ajirascan
```

## Run Tests

```bash
go test ./...
```

## Run ATS Analysis

```bash
go run ./cmd/cli -cv sample_cv.txt -job sample_job.txt
```

## Example Output

```text
====== ATS RESULT ======
Score: 65
Matched: [go docker backend]
Missing: [kubernetes leadership communication]
```

## MVP Vision

AjiraScan aims to become a career intelligence platform that helps users:

- Understand ATS systems
- Improve CV quality
- Tailor resumes to jobs
- Increase interview opportunities
- Prepare competitive job applications

## Future Features

- AI CV Recommendations
- PDF/DOCX Parsing
- Web Dashboard
- User Accounts
- NGO/UN Job Optimization
- LinkedIn Analyzer
- Cover Letter Generator
- Recruiter Dashboard
- API Integration

## Tech Stack

- Go (Golang)
- CLI-first Architecture
- Modular Internal Packages
- Test-Driven Development