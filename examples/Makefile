default:
	@kind delete cluster||true
	@kind create cluster --config kind-config.yml

cleanup:
	@kind delete cluster||true

tail-logs:
	@sudo tail -f logs/audit.log |jq
