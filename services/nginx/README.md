

Build image:
```bash
git clone https://github.com/wgerlach/Skycore.git
cd Skycore/services/nginx/docker
docker rm -f mynginx ; docker rmi mynginx
docker build  --no-cache -t mynginx .
```

Start nginx via confd (8003 is just an example)
```bash
docker run -d -p 8003:80 --name mynginx mynginx
or
docker run -d -p 8003:80 --name mynginx mynginx <cmd>
```
or if you want to pull from git first:
```bash
docker run -d -p 8003:80 --name mynginx mynginx bash -c "cd /root/Skycore && git pull && /root/Skycore/services/nginx/confd/run_confd.sh"
```
