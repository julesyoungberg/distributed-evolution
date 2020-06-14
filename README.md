# Distributed Evolution

TODO:
1. get rpc working: https://www.youtube.com/watch?v=ZONL-6jBevc - DONE
2. get multiple workers running, ideally > 25 - DONE: `docker-compose up --build --scale worker=25`
3. figure out how to import from a util package - DONE
4. create websocket connection between master and interface for displaying the current image - DONE
5. download random image (lorem ipsum?) in master and send to interface - DONE
6. slice image up into jobs
7. send jobs to workers
8. combine result in master
9. iterate
