

Build image:
```bash
git clone https://github.com/wgerlach/Skycore.git
cd Skycore/services/nginx/docker
docker build  --no-cache -t mynginx .
```

Start nginx
```bash
docker run -d -p 8003:80 --name mynginx mynginx nginx
```
or
```bash
docker run -d -p 8003:80 --name mynginx mynginx bash -c "cd /root/Skycore && git pull && nginx"
```
