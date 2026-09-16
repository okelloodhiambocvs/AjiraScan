package web

import (
	"html/template"
	"net/http"
)

type legalSection struct{ Heading, Text string }
type legalPage struct {
	Title    string
	Sections []legalSection
}

func PageHandler(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func LegalHandler(page legalPage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		t, err := template.ParseFiles(templatePath("legal.html"))
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		if err := t.Execute(w, page); err != nil {
			http.Error(w, "internal server error", 500)
		}
	}
}

func legalPages() map[string]legalPage {
	return map[string]legalPage{
		"/about": {"About AjiraScan", []legalSection{
			{"Our mission", "AjiraScan helps job seekers understand how their CV communicates relevant experience, while giving organizations a foundation for fair, explainable recruitment workflows."},
			{"What is available today", "The public analyzer compares CV text with a specific job description and highlights matched terms, missing terms, common sections, and a deterministic alignment score. Account holders can access a protected dashboard."},
			{"How we approach recruitment", "A score is not a hiring decision. AjiraScan is designed to support meaningful human review, job-related criteria, documented explanations, and accountable recruiter decisions."},
			{"Fairness boundaries", "The public analyzer excludes common protected selection terms from its comparison input. It does not assess identity, personality, protected characteristics, or a candidate's worth."},
			{"Privacy and security", "We aim to minimize data, protect documents, use server-side sessions, and support user control. Before production processing of candidate CVs, privacy, retention, access-control, and operational controls must be completed."},
			{"Product roadmap", "Secure CV upload, structured job requirements, applications, recruiter review, and auditable shortlisting will be released only as protected workflows with appropriate authorization and human-review controls."},
		}},
		"/privacy": {"Privacy Policy", []legalSection{{"Information we process", "This may include account details, CVs, applications, employer information, documents, interview details, ATS/AI outputs and payment-related records."}, {"Purpose and sharing", "Information supports the service, security, support and lawful obligations. Providers are engaged only when configured and governed by appropriate contracts."}, {"Your choices", "Users may request access, correction, export or deletion where applicable. Retention, cross-border transfers, children's data, privacy contact and Kenya legal requirements require owner and DPO/legal review."}}},
		"/terms":   {"Terms and Conditions", []legalSection{{"Use of the service", "Users are responsible for accurate, authorized content and secure accounts. Prohibited use includes unlawful content, abuse and attempting to bypass security."}, {"Content and decisions", "Users retain their content subject to service rights needed to operate it. Scores and AI content are not guarantees and do not replace human recruitment decisions."}, {"Business terms", "Subscriptions, refunds, availability, limitations, termination, governing law and contact details require business and legal confirmation before production."}}},
		"/cookies": {"Cookie Policy", []legalSection{{"Currently used", "When authentication is configured, AjiraScan uses an essential HttpOnly session cookie. The public analyzer does not require analytics or advertising cookies."}, {"May be used later", "Preference, analytics or payment cookies may be added only after implementation, notice and consent where required."}, {"Your control", "You can manage cookies through browser controls. Cookie duration, consent tooling and third-party details require final owner review."}}},
	}
}

func AuthUnavailableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "authentication persistence is not configured",
	})
}
