model.vw: train.txt
	vw -d train.txt -f model.vw

.PHONY: test
test: test.txt model.vw
	vw -d test.txt -i model.vw -p predictions.txt
