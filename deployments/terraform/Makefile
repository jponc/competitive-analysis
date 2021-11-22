init-workspaces:
	terraform workspace new staging-competitive-analysis
	terraform workspace new production-competitive-analysis

plan-staging:
	terraform workspace select staging-competitive-analysis
	terraform plan -var-file="staging.tfvars"

plan-production:
	terraform workspace select production-competitive-analysis
	terraform plan -var-file="production.tfvars"

apply-staging:
	terraform workspace select staging-competitive-analysis
	terraform apply -var-file="staging.tfvars"

apply-production:
	terraform workspace select production-competitive-analysis
	terraform apply -var-file="production.tfvars"
