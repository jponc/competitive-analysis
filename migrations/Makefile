migrate-staging:
	sls deploy --stage=staging --verbose
	sls invoke --stage=staging -f UpMigration

migrate-production:
	sls deploy --stage=production --verbose
	sls invoke --stage=production -f UpMigration
