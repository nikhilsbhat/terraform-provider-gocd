scan/code: ## scans code for vulnerabilities
	@docker-compose --project-name trivy -f docker-compose.trivy.yml run --rm trivy fs /terraform-provider-gocd

scan/binary: ## scans binary for vulnerabilities
	@docker-compose --project-name trivy -f docker-compose.trivy.yml run --rm trivy fs /terraform-provider-gocd/dist/terraform-provider-gocd_darwin_amd64_v1/terraform-provider-gocd_v$(VERSION) --scanners vuln
