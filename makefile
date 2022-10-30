push:
	git add .
	git commit -m "$v"
	git push
	git tag $v
	git push origin $v