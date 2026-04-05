

git-commit:
	@git add .
	@git status
	@read -p "Commit message: " msg; \
	git commit -m "$$msg"

run-fe:
	cd frontend && npm run dev

run-be:
	cd backend && go run main.go

