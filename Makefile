up: # starts stack
	docker-compose -f ops/docker-compose.yml --project-dir . up -d
build: # build stack
	docker-compose -f ops/docker-compose.yml --project-dir . up -d --build
down: # stop stack
	docker-compose -f ops/docker-compose.yml --project-dir . down
console-front: # run term in frontend container
	docker-compose -f ops/docker-compose.yml --project-dir . exec frontend sh
console-prometheus: # run term in prometheus container
	docker-compose -f ops/docker-compose.yml --project-dir . exec prometheus sh
console-grafana: # run term in grafana container
	docker-compose -f ops/docker-compose.yml --project-dir . exec grafana sh
console-backend: # run term in backend container
	docker-compose -f ops/docker-compose.yml --project-dir . exec backend sh