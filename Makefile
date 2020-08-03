.PHONY: build

gocv:
	docker build -f gocv.Dockerfile -t julesyoungberg/gocv . && \
	docker push julesyoungberg/gocv

master:
	docker build -f master.Dockerfile -t julesyoungberg/distributed-evolution-master . && \
	docker push julesyoungberg/distributed-evolution-master

worker:
	docker build -f worker.Dockerfile -t julesyoungberg/distributed-evolution-worker . && \
	docker push julesyoungberg/distributed-evolution-worker

sentinel-master:
	docker build -t julesyoungberg/distributed-evolution-sentinel-master ./sentinel/master && \
	docker push julesyoungberg/distributed-evolution-sentinel-master

sentinel-replica:
	docker build -t julesyoungberg/distributed-evolution-sentinel-replica ./sentinel/replica && \
	docker push julesyoungberg/distributed-evolution-sentinel-replica

ui:
	docker build -f prod.Dockerfile -t julesyoungberg/distributed-evolution-ui ./ui && \
	docker push julesyoungberg/distributed-evolution-ui

build: master worker sentinel-master sentinel-replica ui

apply:
	kubectl delete --all deployments && kubectl apply -f deployment/dev

apply_prod:
	kubectl apply -f deployment/prod

build_apply: build apply

start: 
	docker-compose down -v && \
	docker-compose up --build --scale worker=6 --scale sentinel=3 --scale redis-slave=2

master_logs:
	kubectl logs -l app=master -f

worker_logs:
	kubectl logs -l app=worker -f --max-log-requests 10

clean:
	rm -rf .container-* .dockerfile-*
