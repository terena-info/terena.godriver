push:
	git add .
	git commit -m "$v"
	git push
	git tag v$v
	git push origin v$v