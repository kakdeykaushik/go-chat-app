image="go-chat"

# build image if does not exists
if [[ "$(docker images -q $image:latest 2> /dev/null)" == "" ]]; then
	docker build -t $image .
fi

# run container
docker run --network=host -p 8000:8000 $image